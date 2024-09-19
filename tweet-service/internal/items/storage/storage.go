package storage

import (
	"database/sql"
	"log/slog"

	"github.com/abdulazizax/mini-twitter/tweet-service/internal/items/repository"
	psql "github.com/abdulazizax/mini-twitter/tweet-service/internal/items/storage/postgres"
	"github.com/abdulazizax/mini-twitter/tweet-service/internal/pkg/config"

	sq "github.com/Masterminds/squirrel"
)

type StrorageI interface {
	Comments() repository.CommentI
	Likes() repository.LikeI
	Tweets() repository.TweetI
}

type Storage struct {
	tweetRepo   repository.TweetI
	likeRepo    repository.LikeI
	commentRepo repository.CommentI
}

func New(postgres *sql.DB, queryBuilder sq.StatementBuilderType, cfg *config.Config, logger *slog.Logger) StrorageI {
	return &Storage{
		tweetRepo:   psql.NewTweetStorage(postgres, queryBuilder, cfg, logger),
		likeRepo:    psql.NewLikeStorage(postgres, queryBuilder, cfg, logger),
		commentRepo: psql.NewCommentStorage(postgres, queryBuilder, cfg, logger),
	}
}

func (s *Storage) Comments() repository.CommentI {
	return s.commentRepo
}

func (s *Storage) Likes() repository.LikeI {
	return s.likeRepo
}

func (s *Storage) Tweets() repository.TweetI {
	return s.tweetRepo
}
