package service

import (
	"context"
	"log/slog"

	"github.com/abdulazizax/mini-twitter/tweet-service/internal/items/repository"

	pb "github.com/abdulazizax/mini-twitter/tweet-service/genproto/like"
)

type LikeService struct {
	pb.UnimplementedLikeServiceServer
	likestorage repository.LikeI
	logger      *slog.Logger
}

func NewLikesService(likestorage repository.LikeI, logger *slog.Logger) *LikeService {
	return &LikeService{
		likestorage: likestorage,
		logger:      logger,
	}
}

func (s *LikeService) Like(ctx context.Context, in *pb.LikeRequest) (*pb.LikeResponse, error) {
	s.logger.Info("Liked", "username", in.Username)
	return s.likestorage.Like(ctx, in)
}

func (s *LikeService) Unlike(ctx context.Context, in *pb.UnlikeRequest) (*pb.UnlikeResponse, error) {
	s.logger.Info("Unliked", "username", in.Username)
	return s.likestorage.Unlike(ctx, in)
}

func (s *LikeService) GetLikes(ctx context.Context, in *pb.GetLikesRequest) (*pb.GetLikesResponse, error) {
	s.logger.Info("Getting likes")
	return s.likestorage.GetLikes(ctx, in)
}
