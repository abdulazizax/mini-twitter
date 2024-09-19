package service

import (
	"context"
	"log/slog"

	"github.com/abdulazizax/mini-twitter/tweet-service/internal/items/repository"

	pb "github.com/abdulazizax/mini-twitter/tweet-service/genproto/tweet"
)

type TweetService struct {
	pb.UnimplementedTweetServiceServer
	tweetstorage repository.TweetI
	logger       *slog.Logger
}

func NewTweetsService(tweetstorage repository.TweetI, logger *slog.Logger) *TweetService {
	return &TweetService{
		tweetstorage: tweetstorage,
		logger:       logger,
	}
}

// 1
func (s *TweetService) CreateTweet(ctx context.Context, in *pb.CreateTweetRequest) (*pb.CreateTweetResponse, error) {
	s.logger.Info("Creating tweet", "username", in.Username)
	return s.tweetstorage.CreateTweet(ctx, in)
}

// 2
func (s *TweetService) GetTweet(ctx context.Context, in *pb.GetTweetRequest) (*pb.GetTweetResponse, error) {
	s.logger.Info("Getting tweet", "username", in.Username)
	return s.tweetstorage.GetTweet(ctx, in)
}

// 3
func (s *TweetService) UpdateTweet(ctx context.Context, in *pb.UpdateTweetRequest) (*pb.UpdateTweetResponse, error) {
	s.logger.Info("Updating tweet", "username", in.Username)
	return s.tweetstorage.UpdateTweet(ctx, in)
}

// 4
func (s *TweetService) DeleteTweet(ctx context.Context, in *pb.DeleteTweetRequest) (*pb.Status, error) {
	s.logger.Info("Deleting tweet", "username", in.Username)
	return s.tweetstorage.DeleteTweet(ctx, in)
}

// 5
func (s *TweetService) GetAllTweets(ctx context.Context, in *pb.GetAllTweetsRequest) (*pb.GetAllTweetsResponse, error) {
	s.logger.Info("Getting all tweets", "username", in.Username)
	return s.tweetstorage.GetAllTweets(ctx, in)
}

// 6
func (s *TweetService) IncreaseViewsCount(ctx context.Context, in *pb.Id) (*pb.Status, error) {
	s.logger.Info("Increasing views count", "username", in.Username)
	return s.tweetstorage.IncreaseViewsCount(ctx, in)
}

// 7
func (s *TweetService) IncreaseRepostCount(ctx context.Context, in *pb.Id) (*pb.Status, error) {
	s.logger.Info("Increasing repost count", "username", in.Username)
	return s.tweetstorage.IncreaseRepostCount(ctx, in)
}

// 8
func (s *TweetService) IncreaseSharesCount(ctx context.Context, in *pb.Id) (*pb.Status, error) {
	s.logger.Info("Increasing shares count", "username", in.Username)
	return s.tweetstorage.IncreaseSharesCount(ctx, in)
}
