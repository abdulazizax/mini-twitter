package repository

import (
	"context"

	pb "github.com/abdulazizax/mini-twitter/tweet-service/genproto/like"
)

type LikeI interface {
	Like(ctx context.Context, in *pb.LikeRequest) (*pb.LikeResponse, error)
	Unlike(ctx context.Context, in *pb.UnlikeRequest) (*pb.UnlikeResponse, error)
	GetLikes(ctx context.Context, in *pb.GetLikesRequest) (*pb.GetLikesResponse, error)
}
