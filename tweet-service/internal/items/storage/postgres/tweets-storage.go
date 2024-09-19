package psql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	pb "github.com/abdulazizax/mini-twitter/tweet-service/genproto/tweet"
	"github.com/abdulazizax/mini-twitter/tweet-service/internal/items/repository"
	"github.com/abdulazizax/mini-twitter/tweet-service/internal/pkg/config"
	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type TweetStorage struct {
	postgres     *sql.DB
	queryBuilder sq.StatementBuilderType
	cfg          *config.Config
	logger       *slog.Logger
}

func NewTweetStorage(postgres *sql.DB, queryBuilder sq.StatementBuilderType, cfg *config.Config, logger *slog.Logger) repository.TweetI {
	return &TweetStorage{
		postgres:     postgres,
		queryBuilder: queryBuilder,
		cfg:          cfg,
		logger:       logger,
	}
}

// 1
func (s *TweetStorage) CreateTweet(ctx context.Context, in *pb.CreateTweetRequest) (*pb.CreateTweetResponse, error) {
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

	// Get the latest tweet_serial for the user
	query, args, err := s.queryBuilder.Select("COALESCE(MAX(tweet_serial), 0)").
		From("tweets").
		Where(sq.Eq{"username": in.Username}).
		ToSql()

	if err != nil {
		s.logger.Error("Error while building tweet_serial query", slog.String("err", err.Error()))
		return nil, err
	}

	var tweet_serial int32
	err = tx.QueryRow(query, args...).Scan(&tweet_serial)
	if err != nil && err != sql.ErrNoRows {
		s.logger.Error("Error while getting tweet serial", slog.String("err", err.Error()))
		return nil, err
	}

	// Prepare data for new tweet
	id := uuid.New().String()
	created_at := time.Now()

	// Insert the new tweet
	query, args, err = s.queryBuilder.Insert("tweets").
		Columns(
			"id",
			"username",
			"tweet_serial",
			"content",
			"media",
			"created_at",
			"updated_at").
		Values(
			id,
			in.Username,
			tweet_serial+1, // Increment the tweet serial
			in.Content,
			pq.Array(in.Media), // Use pq.Array for array types
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

	// Commit transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error("Error while committing transaction", slog.String("err", err.Error()))
		return nil, err
	}

	// Set tx to nil after commit to avoid Rollback in defer
	tx = nil

	return &pb.CreateTweetResponse{
		Tweet: &pb.Tweet{
			Id:            id,
			Username:      in.Username,
			TweetSerial:   tweet_serial + 1, // Increment the tweet serial
			Content:       in.Content,
			Media:         in.Media,
			CommentsCount: 0,
			ViewsCount:    0,
			RepostCount:   0,
			SharesCount:   0,
			CreatedAt:     timestamppb.New(created_at),
			UpdatedAt:     timestamppb.New(created_at),
		},
	}, nil
}

// 2
func (s *TweetStorage) GetTweet(ctx context.Context, in *pb.GetTweetRequest) (*pb.GetTweetResponse, error) {
	// Query to select the tweet based on username and tweet_serial
	query, args, err := s.queryBuilder.Select(
		"id",
		"content",
		"media",
		"comments_count",
		"views_count",
		"repost_count",
		"shares_count",
		"created_at",
		"updated_at").
		From("tweets").
		Where(sq.Eq{"tweet_serial": in.TweetSerial}).
		Where(sq.Eq{"username": in.Username}).
		ToSql()

	if err != nil {
		s.logger.Error("Error while building query", slog.String("err", err.Error()))
		return nil, err
	}

	var tweet pb.Tweet
	var media pq.StringArray // To handle VARCHAR[] from Postgres
	var created_at, updated_at time.Time

	// Scan the result into the tweet struct
	err = s.postgres.QueryRow(query, args...).Scan(
		&tweet.Id,
		&tweet.Content,
		&media, // Scanning the media as pq.StringArray
		&tweet.CommentsCount,
		&tweet.ViewsCount,
		&tweet.RepostCount,
		&tweet.SharesCount,
		&created_at,
		&updated_at,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("Tweet not found", slog.String("username", in.Username), slog.Any("tweet_serial", in.TweetSerial))
			return nil, fmt.Errorf("tweet not found")
		}
		s.logger.Error("Error while scanning tweet", slog.String("err", err.Error()))
		return nil, err
	}

	// Assign the scanned media array to the tweet's Media field
	tweet.Media = media

	// Assign username and tweet_serial from input
	tweet.Username = in.Username
	tweet.TweetSerial = in.TweetSerial

	tweet.CreatedAt = timestamppb.New(created_at)
	tweet.UpdatedAt = timestamppb.New(updated_at)

	return &pb.GetTweetResponse{
		Tweet: &tweet,
	}, nil
}

// 3
func (s *TweetStorage) UpdateTweet(ctx context.Context, in *pb.UpdateTweetRequest) (*pb.UpdateTweetResponse, error) {
	updatedAt := time.Now()

	// List of dynamic updates
	queryBuilder := s.queryBuilder.Update("tweets")

	// Add non-empty values for updating
	fieldsUpdated := false

	if in.Content != "" {
		queryBuilder = queryBuilder.Set("content", in.Content)
		fieldsUpdated = true
	}
	if len(in.Media) > 0 { // Check if media is not empty
		queryBuilder = queryBuilder.Set("media", pq.Array(in.Media))
		fieldsUpdated = true
	}

	// If no updates, return the existing tweet
	if !fieldsUpdated {
		return nil, fmt.Errorf("no fields to update")
	}

	// Add 'updated_at' field
	queryBuilder = queryBuilder.Set("updated_at", updatedAt).
		Where(sq.Eq{"tweet_serial": in.TweetSerial}).
		Where(sq.Eq{"username": in.Username})

	// Create SQL query
	query, args, err := queryBuilder.Suffix("RETURNING id, content, media, comments_count, views_count, repost_count, shares_count, created_at, updated_at").ToSql()
	if err != nil {
		s.logger.Error("Error while building update query", slog.String("err", err.Error()))
		return nil, err
	}

	// Execute query and scan results to return the tweet
	var tweet pb.Tweet
	var media pq.StringArray
	var created_at, updated_at time.Time

	err = s.postgres.QueryRow(query, args...).Scan(
		&tweet.Id,
		&tweet.Content,
		&media,
		&tweet.CommentsCount,
		&tweet.ViewsCount,
		&tweet.RepostCount,
		&tweet.SharesCount,
		&created_at,
		&updated_at,
	)
	if err != nil {
		s.logger.Error("Error while scanning updated tweet", slog.String("err", err.Error()))
		return nil, err
	}

	tweet.Username = in.Username
	tweet.TweetSerial = in.TweetSerial
	tweet.Media = media

	tweet.CreatedAt = timestamppb.New(created_at)
	tweet.UpdatedAt = timestamppb.New(updated_at)

	// Return the updated tweet
	return &pb.UpdateTweetResponse{
		Tweet: &tweet,
	}, nil
}

// 4
func (s *TweetStorage) DeleteTweet(ctx context.Context, in *pb.DeleteTweetRequest) (*pb.Status, error) {
	// Create a query to update the 'deleted_at' field
	query, args, err := s.queryBuilder.Update("tweets").
		Set("deleted_at", time.Now()).
		Where(sq.Eq{"tweet_serial": in.TweetSerial}).
		Where(sq.Eq{"username": in.Username}).
		ToSql()
	if err != nil {
		s.logger.Error("Error while building update query", slog.String("err", err.Error()))
		return nil, err
	}

	// Execute the SQL query
	result, err := s.postgres.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Error while executing delete query", slog.String("err", err.Error()))
		return nil, err
	}

	// Check the number of affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("Error while checking rows affected", slog.String("err", err.Error()))
		return nil, err
	}

	// If no rows were modified, return an error
	if rowsAffected == 0 {
		s.logger.Error("Tweet not found or already deleted", slog.Any("tweet_serial", in.TweetSerial))
		return nil, fmt.Errorf("tweet not found or already deleted")
	}

	// Return a successful deletion message
	return &pb.Status{
		Success: true,
	}, nil
}

// 5
func (s *TweetStorage) GetAllTweets(ctx context.Context, in *pb.GetAllTweetsRequest) (*pb.GetAllTweetsResponse, error) {
	// Create SQL query
	query, args, err := s.queryBuilder.Select(
		"id",
		"tweet_serial",
		"content",
		"media",
		"comments_count",
		"views_count",
		"repost_count",
		"shares_count",
		"created_at",
		"updated_at").
		From("tweets").
		Where(sq.Eq{"username": in.Username}).
		Where(sq.Eq{"deleted_at": nil}).
		OrderBy("tweet_serial ASC").
		ToSql()
	if err != nil {
		s.logger.Error("Error while building query", slog.String("err", err.Error()))
		return nil, err
	}

	// Execute SQL query
	rows, err := s.postgres.Query(query, args...)
	if err != nil {
		s.logger.Error("Error while executing query", slog.String("err", err.Error()))
		return nil, err
	}
	defer rows.Close()

	// Collect results
	var tweets pb.GetAllTweetsResponse
	for rows.Next() {
		var tweet pb.Tweet
		var createdAt, updatedAt sql.NullTime
		var media pq.StringArray

		err := rows.Scan(
			&tweet.Id,
			&tweet.TweetSerial,
			&tweet.Content,
			&media,
			&tweet.CommentsCount,
			&tweet.ViewsCount,
			&tweet.RepostCount,
			&tweet.SharesCount,
			&createdAt,
			&updatedAt,
			// 'deleted_at' field is not read here
		)
		if err != nil {
			s.logger.Error("Error while scanning tweet row", slog.String("err", err.Error()))
			return nil, err
		}

		// Convert times to 'google.protobuf.Timestamp' format
		if createdAt.Valid {
			tweet.CreatedAt = timestamppb.New(createdAt.Time)
		}
		if updatedAt.Valid {
			tweet.UpdatedAt = timestamppb.New(updatedAt.Time)
		}

		tweet.Username = in.Username
		tweet.Media = media

		// Add tweet to the results list
		tweets.Tweets = append(tweets.Tweets, &tweet)
	}

	// Check for errors after scanning rows
	if err := rows.Err(); err != nil {
		s.logger.Error("Error after scanning rows", slog.String("err", err.Error()))
		return nil, err
	}

	return &tweets, nil
}

// 6
func (s *TweetStorage) IncreaseViewsCount(ctx context.Context, in *pb.Id) (*pb.Status, error) {
	// Start a transaction
	tx, err := s.postgres.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("Error while starting a transaction", slog.String("err", err.Error()))
		return nil, err
	}
	defer tx.Rollback()

	// Get the current views count
	var currentViewsCount int32
	query, args, err := s.queryBuilder.Select("views_count").
		From("tweets").
		Where(sq.Eq{"username": in.Username}).
		Where(sq.Eq{"tweet_serial": in.TweetSerial}).
		ToSql()
	if err != nil {
		s.logger.Error("Error while building query", slog.String("err", err.Error()))
		return nil, err
	}

	err = tx.QueryRowContext(ctx, query, args...).Scan(&currentViewsCount)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("Tweet not found", slog.String("username", in.Username), slog.Any("tweet_serial", in.TweetSerial))
			return nil, fmt.Errorf("tweet not found")
		}
		s.logger.Error("Error while scanning views count", slog.String("err", err.Error()))
		return nil, err
	}

	// Increment the views count by 1
	query, args, err = s.queryBuilder.Update("tweets").
		Set("views_count", currentViewsCount+1).
		Where(sq.Eq{"username": in.Username}).
		Where(sq.Eq{"tweet_serial": in.TweetSerial}).
		ToSql()
	if err != nil {
		s.logger.Error("Error while building update query", slog.String("err", err.Error()))
		return nil, err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Error while executing update query", slog.String("err", err.Error()))
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error("Error while committing transaction", slog.String("err", err.Error()))
		return nil, err
	}

	return &pb.Status{
		Success: true,
	}, nil
}

// 7
func (s *TweetStorage) IncreaseRepostCount(ctx context.Context, in *pb.Id) (*pb.Status, error) {
	// Start a transaction
	tx, err := s.postgres.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("Error while starting a transaction", slog.String("err", err.Error()))
		return nil, err
	}
	defer tx.Rollback()

	// Get the current repost count
	var currentRepostCount int32
	query, args, err := s.queryBuilder.Select("repost_count").
		From("tweets").
		Where(sq.Eq{"username": in.Username}).
		Where(sq.Eq{"tweet_serial": in.TweetSerial}).
		ToSql()
	if err != nil {
		s.logger.Error("Error while building query", slog.String("err", err.Error()))
		return nil, err
	}

	err = tx.QueryRowContext(ctx, query, args...).Scan(&currentRepostCount)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("Tweet not found", slog.String("username", in.Username), slog.Any("tweet_serial", in.TweetSerial))
			return nil, fmt.Errorf("tweet not found")
		}
		s.logger.Error("Error while scanning repost count", slog.String("err", err.Error()))
		return nil, err
	}

	// Increment the repost count by 1
	query, args, err = s.queryBuilder.Update("tweets").
		Set("repost_count", currentRepostCount+1).
		Where(sq.Eq{"username": in.Username}).
		Where(sq.Eq{"tweet_serial": in.TweetSerial}).
		ToSql()
	if err != nil {
		s.logger.Error("Error while building update query", slog.String("err", err.Error()))
		return nil, err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Error while executing update query", slog.String("err", err.Error()))
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error("Error while committing transaction", slog.String("err", err.Error()))
		return nil, err
	}

	return &pb.Status{
		Success: true,
	}, nil
}

// 8
func (s *TweetStorage) IncreaseSharesCount(ctx context.Context, in *pb.Id) (*pb.Status, error) {
	// Start a transaction
	tx, err := s.postgres.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("Error while starting a transaction", slog.String("err", err.Error()))
		return nil, err
	}
	defer tx.Rollback()

	// Get the current shares count
	var currentSharesCount int32
	query, args, err := s.queryBuilder.Select("shares_count").
		From("tweets").
		Where(sq.Eq{"username": in.Username}).
		Where(sq.Eq{"tweet_serial": in.TweetSerial}).
		ToSql()
	if err != nil {
		s.logger.Error("Error while building query", slog.String("err", err.Error()))
		return nil, err
	}

	err = tx.QueryRowContext(ctx, query, args...).Scan(&currentSharesCount)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("Tweet not found", slog.String("username", in.Username), slog.Any("tweet_serial", in.TweetSerial))
			return nil, fmt.Errorf("tweet not found")
		}
		s.logger.Error("Error while scanning shares count", slog.String("err", err.Error()))
		return nil, err
	}

	// Increment the shares count by 1
	query, args, err = s.queryBuilder.Update("tweets").
		Set("shares_count", currentSharesCount+1).
		Where(sq.Eq{"username": in.Username}).
		Where(sq.Eq{"tweet_serial": in.TweetSerial}).
		ToSql()
	if err != nil {
		s.logger.Error("Error while building update query", slog.String("err", err.Error()))
		return nil, err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Error while executing update query", slog.String("err", err.Error()))
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error("Error while committing transaction", slog.String("err", err.Error()))
		return nil, err
	}

	return &pb.Status{
		Success: true,
	}, nil
}
