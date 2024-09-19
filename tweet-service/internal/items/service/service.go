package service

import (
	"log/slog"

	"github.com/abdulazizax/mini-twitter/tweet-service/internal/items/storage"
)

type Service struct {
	CommentService *CommentService
	LikeService    *LikeService
	TweetService   *TweetService
}

func New(storage storage.StrorageI, logger *slog.Logger) *Service {
	return &Service{
		CommentService: NewCommentsService(storage.Comments(), logger),
		LikeService:    NewLikesService(storage.Likes(), logger),
		TweetService:   NewTweetsService(storage.Tweets(), logger),
	}

}
