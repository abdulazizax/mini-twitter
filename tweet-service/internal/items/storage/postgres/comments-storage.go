package psql

import (
	"database/sql"
	"log/slog"
	"time"

	pb "github.com/abdulazizax/mini-twitter/tweet-service/genproto/comment"
	"github.com/abdulazizax/mini-twitter/tweet-service/internal/items/repository"
	"github.com/abdulazizax/mini-twitter/tweet-service/internal/pkg/config"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/timestamppb"

	sq "github.com/Masterminds/squirrel"
)

type CommentStorage struct {
	postgres     *sql.DB
	queryBuilder sq.StatementBuilderType
	cfg          *config.Config
	logger       *slog.Logger
}

func NewCommentStorage(postgres *sql.DB, queryBuilder sq.StatementBuilderType, cfg *config.Config, logger *slog.Logger) repository.CommentI {
	return &CommentStorage{
		postgres:     postgres,
		queryBuilder: queryBuilder,
		cfg:          cfg,
		logger:       logger,
	}
}

func (s *CommentStorage) CreateComment(ctx context.Context, in *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	tx, err := s.postgres.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("Error while starting a transaction", slog.String("err", err.Error()))
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Rollback() // Rollback if not committed
		}
	}()

	// Get the latest comment_serial for the tweet
	query, args, err := s.queryBuilder.Select("COALESCE(MAX(comment_serial), 0)").
		From("comments").
		Where(sq.Eq{"tweet_id": in.TweetId}).
		ToSql()
	if err != nil {
		s.logger.Error("Error while building comment_serial query", slog.String("err", err.Error()))
		return nil, err
	}

	var comment_serial int32
	err = tx.QueryRow(query, args...).Scan(&comment_serial)
	if err != nil && err != sql.ErrNoRows {
		s.logger.Error("Error while getting comment serial", slog.String("err", err.Error()))
		return nil, err
	}

	// Prepare data for new comment
	id := uuid.New().String()
	created_at := time.Now()

	// Insert the new comment
	query, args, err = s.queryBuilder.Insert("comments").
		Columns(
			"id",
			"username",
			"tweet_id",
			"comment_serial",
			"content",
			"created_at",
			"updated_at",
		).Values(
		id,
		in.Username,
		in.TweetId,
		comment_serial+1, // Increment comment serial
		in.Content,
		created_at,
		created_at,
	).ToSql()

	if err != nil {
		s.logger.Error("Error while building insert query", slog.String("err", err.Error()))
		return nil, err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		s.logger.Error("Error while executing insert query", slog.String("err", err.Error()))
		return nil, err
	}

	// Update the comments_count in the tweet
	query, args, err = s.queryBuilder.Update("tweets").
		Set("comments_count", sq.Expr("comments_count + 1")).
		Where(sq.Eq{"id": in.TweetId}).
		ToSql()
	if err != nil {
		s.logger.Error("Error while building update query for comments_count", slog.String("err", err.Error()))
		return nil, err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		s.logger.Error("Error while updating comments_count", slog.String("err", err.Error()))
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error("Error while committing transaction", slog.String("err", err.Error()))
		return nil, err
	}

	// Set tx to nil after commit to avoid Rollback in defer
	tx = nil

	return &pb.CreateCommentResponse{
		Comment: &pb.Comment{
			Id:            id,
			Username:      in.Username,
			TweetId:       in.TweetId,
			CommentSerial: comment_serial + 1,
			Content:       in.Content,
			CreatedAt:     timestamppb.New(created_at),
			UpdatedAt:     timestamppb.New(created_at),
		},
	}, nil
}

func (s *CommentStorage) DeleteComment(ctx context.Context, in *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	tx, err := s.postgres.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("Error while starting a transaction", slog.String("err", err.Error()))
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Rollback() // Rollback if not committed
		}
	}()

	// Delete the comment by setting deleted_at timestamp
	query, args, err := s.queryBuilder.Update("comments").
		Set("deleted_at", time.Now()).
		Where(sq.Eq{"id": in.CommentId}).
		ToSql()

	if err != nil {
		s.logger.Error("Error while building delete query", slog.String("err", err.Error()))
		return nil, err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		s.logger.Error("Error while executing delete query", slog.String("err", err.Error()))
		return nil, err
	}

	// Decrease the comments_count in the tweet
	query, args, err = s.queryBuilder.Update("tweets").
		Set("comments_count", sq.Expr("comments_count - 1")).
		Where(sq.Eq{"id": in.TweetId}).
		ToSql()
	if err != nil {
		s.logger.Error("Error while building update query for comments_count", slog.String("err", err.Error()))
		return nil, err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		s.logger.Error("Error while updating comments_count", slog.String("err", err.Error()))
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error("Error while committing transaction", slog.String("err", err.Error()))
		return nil, err
	}

	// Set tx to nil after commit to avoid Rollback in defer
	tx = nil

	return &pb.DeleteCommentResponse{
		Message: "Comment deleted successfully",
	}, nil
}

func (s *CommentStorage) GetCommentsForTweet(ctx context.Context, in *pb.GetCommentsForTweetRequest) (*pb.GetCommentsForTweetResponse, error) {
	// Define the query to select comments for the given tweet_id
	query, args, err := s.queryBuilder.Select(
		"id",
		"username",
		"tweet_id",
		"comment_serial",
		"content",
		"created_at",
		"updated_at").
		From("comments").
		Where(sq.Eq{"tweet_id": in.TweetId}).
		Where(sq.Eq{"deleted_at": nil}). // Make sure we only fetch non-deleted comments
		OrderBy("comment_serial ASC").   // Sort by comment_serial in ascending order
		ToSql()

	if err != nil {
		s.logger.Error("Error while building query for comments", slog.String("err", err.Error()))
		return nil, err
	}

	// Execute the query
	rows, err := s.postgres.QueryContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Error while executing query for comments", slog.String("err", err.Error()))
		return nil, err
	}
	defer rows.Close()

	// Prepare the response structure
	var comments []*pb.Comment
	for rows.Next() {
		var comment pb.Comment
		var createdAt, updatedAt time.Time

		// Scan each row into a Comment object
		err := rows.Scan(
			&comment.Id,
			&comment.Username,
			&comment.TweetId,
			&comment.CommentSerial,
			&comment.Content,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			s.logger.Error("Error while scanning comment row", slog.String("err", err.Error()))
			return nil, err
		}

		// Convert the timestamps to protobuf format
		comment.CreatedAt = timestamppb.New(createdAt)
		comment.UpdatedAt = timestamppb.New(updatedAt)

		// Append the comment to the list
		comments = append(comments, &comment)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		s.logger.Error("Error while iterating over comment rows", slog.String("err", err.Error()))
		return nil, err
	}

	// Return the list of comments
	return &pb.GetCommentsForTweetResponse{
		Comments: comments,
	}, nil
}
