package psql

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	pb "github.com/abdulazizax/mini-twitter/tweet-service/genproto/like"
	"github.com/abdulazizax/mini-twitter/tweet-service/internal/items/repository"
	"github.com/abdulazizax/mini-twitter/tweet-service/internal/pkg/config"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	sq "github.com/Masterminds/squirrel"
)

type LikeStorage struct {
	postgres     *sql.DB
	queryBuilder sq.StatementBuilderType
	cfg          *config.Config
	logger       *slog.Logger
}

func NewLikeStorage(postgres *sql.DB, queryBuilder sq.StatementBuilderType, cfg *config.Config, logger *slog.Logger) repository.LikeI {
	return &LikeStorage{
		postgres:     postgres,
		queryBuilder: queryBuilder,
		cfg:          cfg,
		logger:       logger,
	}
}

func (s *LikeStorage) Like(ctx context.Context, in *pb.LikeRequest) (*pb.LikeResponse, error) {
	id := uuid.New().String()
	created_at := time.Now()

	query, args, err := s.queryBuilder.Insert("likes").
		Columns(
			"id",
			"username",
			"target_id",
			"target_type",
			"created_at",
		).Values(
		id,
		in.Username,
		in.TargetId,
		in.TargetType,
		created_at,
	).ToSql()

	if err != nil {
		s.logger.Error("Error while building insert query", slog.String("err", err.Error()))
		return &pb.LikeResponse{}, err
	}

	_, err = s.postgres.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Error while executing insert query", slog.String("err", err.Error()))
		return &pb.LikeResponse{}, err
	}

	return &pb.LikeResponse{
		Message: "Like created successfully",
	}, nil
}

func (s *LikeStorage) Unlike(ctx context.Context, in *pb.UnlikeRequest) (*pb.UnlikeResponse, error) {
	query, args, err := s.queryBuilder.Delete("likes").
		Where(sq.Eq{
			"username":    in.Username,
			"target_id":   in.TargetId,
			"target_type": in.TargetType,
		}).ToSql()

	if err != nil {
		s.logger.Error("Error while building delete query", slog.String("err", err.Error()))
		return &pb.UnlikeResponse{}, err
	}

	_, err = s.postgres.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Error while executing delete query", slog.String("err", err.Error()))
		return &pb.UnlikeResponse{}, err
	}

	return &pb.UnlikeResponse{
		Message: "Like removed successfully",
	}, nil
}

func (s *LikeStorage) GetLikes(ctx context.Context, in *pb.GetLikesRequest) (*pb.GetLikesResponse, error) {
	query, args, err := s.queryBuilder.Select("username", "created_at").
		From("likes").
		Where(sq.Eq{
			"target_id":   in.TargetId,
			"target_type": in.TargetType,
		}).
		OrderBy("created_at DESC").
		ToSql()

	if err != nil {
		s.logger.Error("Error while building select query", slog.String("err", err.Error()))
		return &pb.GetLikesResponse{}, err
	}

	rows, err := s.postgres.QueryContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Error while executing select query", slog.String("err", err.Error()))
		return &pb.GetLikesResponse{}, err
	}
	defer rows.Close()

	var likes []*pb.LikeInfo
	for rows.Next() {
		like := &pb.LikeInfo{}
		var createdAt time.Time
		err := rows.Scan(&like.Username, &createdAt)
		if err != nil {
			s.logger.Error("Error while scanning row", slog.String("err", err.Error()))
			return &pb.GetLikesResponse{}, err
		}
		like.LikedAt = timestamppb.New(createdAt)
		likes = append(likes, like)
	}

	if err = rows.Err(); err != nil {
		s.logger.Error("Error after iterating rows", slog.String("err", err.Error()))
		return &pb.GetLikesResponse{}, err
	}

	return &pb.GetLikesResponse{
		Likes: likes,
	}, nil
}
