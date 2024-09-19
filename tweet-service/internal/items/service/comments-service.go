package service

import (
	"context"
	"log/slog"

	"github.com/abdulazizax/mini-twitter/tweet-service/internal/items/repository"

	pb "github.com/abdulazizax/mini-twitter/tweet-service/genproto/comment"
)

type CommentService struct {
	pb.UnimplementedCommentServiceServer
	commentstorage repository.CommentI
	logger         *slog.Logger
}

func NewCommentsService(commentstorage repository.CommentI, logger *slog.Logger) *CommentService {
	return &CommentService{
		commentstorage: commentstorage,
		logger:         logger,
	}
}

func (s *CommentService) CreateComment(ctx context.Context, in *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	s.logger.Info("Creating comment", "username", in.Username)
	return s.commentstorage.CreateComment(ctx, in)
}

func (s *CommentService) DeleteComment(ctx context.Context, in *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	s.logger.Info("Deleting comment", "username", in.Username)
	return s.commentstorage.DeleteComment(ctx, in)
}

func (s *CommentService) GetCommentsForTweet(ctx context.Context, in *pb.GetCommentsForTweetRequest) (*pb.GetCommentsForTweetResponse, error) {
	s.logger.Info("Getting comments for tweet")
	return s.commentstorage.GetCommentsForTweet(ctx, in)
}
