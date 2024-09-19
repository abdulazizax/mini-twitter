package repository

import (
	"context"

	pb "github.com/abdulazizax/mini-twitter/tweet-service/genproto/comment"
)

type CommentI interface {
	// 1
	CreateComment(ctx context.Context, in *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error)
	// 2
	DeleteComment(ctx context.Context, in *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error)
	// 3
	GetCommentsForTweet(ctx context.Context, in *pb.GetCommentsForTweetRequest) (*pb.GetCommentsForTweetResponse, error)
}
