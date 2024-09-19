package repository

import (
	"context"

	pb "github.com/abdulazizax/mini-twitter/tweet-service/genproto/tweet"
)

type TweetI interface {
	// 1
	CreateTweet(ctx context.Context, in *pb.CreateTweetRequest) (*pb.CreateTweetResponse, error)
	// 2
	GetTweet(ctx context.Context, in *pb.GetTweetRequest) (*pb.GetTweetResponse, error)
	// 3
	UpdateTweet(ctx context.Context, in *pb.UpdateTweetRequest) (*pb.UpdateTweetResponse, error)
	// 4
	DeleteTweet(ctx context.Context, in *pb.DeleteTweetRequest) (*pb.Status, error)
	// 5
	GetAllTweets(ctx context.Context, in *pb.GetAllTweetsRequest) (*pb.GetAllTweetsResponse, error)
	// 6
	IncreaseViewsCount(ctx context.Context, in *pb.Id) (*pb.Status, error)
	// 7
	IncreaseRepostCount(ctx context.Context, in *pb.Id) (*pb.Status, error)
	// 8	
	IncreaseSharesCount(ctx context.Context, in *pb.Id) (*pb.Status, error)
}
