package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	pb "github.com/abdulazizax/mini-twitter/user-service/genproto/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Storage) FollowUser(ctx context.Context, in *pb.FollowUserRequest) (*pb.FollowUserResponse, error) {
	tx, err := s.postgres.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("Error while starting a transaction", "error", err)
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now()

	// Check if the user is already following the target user
	query, args, err := s.queryBuilder.Select("follower_id", "followed_id", "status").
		From("followers").
		Where(sq.Eq{"follower_id": in.FollowerId, "followed_id": in.FollowedId}).
		ToSql()
	if err != nil {
		s.logger.Error("Error while building the query", "error", err)
		return nil, err
	}

	var followerID, followedID string
	var status bool
	err = tx.QueryRowContext(ctx, query, args...).Scan(&followerID, &followedID, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			// User is not following the target user, proceed with the follow action
			query, args, err := s.queryBuilder.Insert("followers").
				Columns("follower_id", "followed_id", "status", "created_at", "updated_at").
				Values(in.FollowerId, in.FollowedId, true, now, now).
				ToSql()
			if err != nil {
				s.logger.Error("Error while building the insert query", "error", err)
				return nil, err
			}
			_, err = tx.ExecContext(ctx, query, args...)
			if err != nil {
				s.logger.Error("Error while executing the insert query", "error", err)
				return nil, err
			}
		} else {
			s.logger.Error("Error while checking if user is already following", "error", err)
			return nil, err
		}
	} else if !status {
		// User has previously unfollowed, update the status to true
		query, args, err := s.queryBuilder.Update("followers").
			Set("status", true).
			Set("updated_at", now).
			Where(sq.Eq{"follower_id": in.FollowerId, "followed_id": in.FollowedId}).
			ToSql()
		if err != nil {
			s.logger.Error("Error while building the update query", "error", err)
			return nil, err
		}

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			s.logger.Error("Error while executing the update query", "error", err)
			return nil, err
		}
	} else {
		// User is already following the target user
		return &pb.FollowUserResponse{Success: true}, nil
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("Error while committing the transaction", "error", err)
		return nil, err
	}

	return &pb.FollowUserResponse{Success: true}, nil
}

func (s *Storage) UnfollowUser(ctx context.Context, in *pb.UnfollowUserRequest) (*pb.UnfollowUserResponse, error) {
	// Start a transaction
	tx, err := s.postgres.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("Error while starting a transaction", "error", err)
		return nil, err
	}
	defer tx.Rollback()

	updatedAt := time.Now()

	// Check if the user is following
	query, args, err := s.queryBuilder.Select("follower_id", "followed_id", "status").
		From("followers").
		Where(sq.Eq{"follower_id": in.FollowerId, "followed_id": in.FollowedId}).
		ToSql()
	if err != nil {
		s.logger.Error("Error while building the query", "error", err)
		return nil, err
	}

	var followerID, followedID string
	var status bool
	err = tx.QueryRowContext(ctx, query, args...).Scan(&followerID, &followedID, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			// If the user has never followed before
			return &pb.UnfollowUserResponse{
				Success: false,
			}, nil
		}
		s.logger.Error("Error while checking if user is following", "error", err)
		return nil, err
	}

	if status {
		// If the user is following and wants to unfollow now
		query, args, err = s.queryBuilder.Update("followers").
			Set("status", false).
			Set("updated_at", updatedAt).
			Where(sq.Eq{"follower_id": in.FollowerId, "followed_id": in.FollowedId}).
			ToSql()
		if err != nil {
			s.logger.Error("Error while building the update query", "error", err)
			return nil, err
		}

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			s.logger.Error("Error while executing the update query", "error", err)
			return nil, err
		}
	} else {
		// If the user has already unfollowed
		return &pb.UnfollowUserResponse{
			Success: false,
		}, nil
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("Error while committing the transaction", "error", err)
		return nil, err
	}

	return &pb.UnfollowUserResponse{Success: true}, nil
}

func (s *Storage) GetFollowers(ctx context.Context, in *pb.GetFollowersRequest) (*pb.GetFollowersResponse, error) {
	query, args, err := s.queryBuilder.Select(
		"u.id",
		"u.email",
		"u.username",
		"u.first_name",
		"u.last_name",
		"u.phone_number",
		"u.bio",
		"u.profile_picture",
		"u.created_at",
		"u.updated_at").
		From("followers f").
		Join("users u ON u.id = f.follower_id").
		Where(sq.Eq{"f.followed_id": in.UserId, "f.status": true}).
		ToSql()

	if err != nil {
		s.logger.Error("Error while building the query", "error", err)
		return nil, err
	}

	log.Println(query)

	rows, err := s.postgres.QueryContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Error while executing the query", "error", err)
		return nil, err
	}
	defer rows.Close()

	var followers []*pb.User
	for rows.Next() {
		var user pb.User
		var createdAt, updatedAt time.Time
		var phoneNumber sql.NullString // Using sql.NullString to handle NULL values

		err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&phoneNumber, // Scanning phone_number into sql.NullString
			&user.Bio,
			&user.ProfilePictureUrl,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			s.logger.Error("Error while scanning row", "error", err)
			return nil, err
		}

		// Set PhoneNumber to empty string if it's NULL
		if phoneNumber.Valid {
			user.PhoneNumber = phoneNumber.String
		} else {
			user.PhoneNumber = ""
		}

		user.CreatedAt = timestamppb.New(createdAt)
		user.UpdatedAt = timestamppb.New(updatedAt)

		followers = append(followers, &user)
	}

	if err := rows.Err(); err != nil {
		s.logger.Error("Error after iterating through rows", "error", err)
		return nil, err
	}

	return &pb.GetFollowersResponse{Followers: followers}, nil
}

func (s *Storage) GetFollowing(ctx context.Context, in *pb.GetFollowingRequest) (*pb.GetFollowingResponse, error) {
	query, args, err := s.queryBuilder.Select(
		"u.id",
		"u.email",
		"u.username",
		"u.first_name",
		"u.last_name",
		"u.phone_number",
		"u.bio",
		"u.profile_picture",
		"u.created_at",
		"u.updated_at").
		From("followers f").
		Join("users u ON u.id = f.followed_id").
		Where(sq.Eq{"f.follower_id": in.UserId, "f.status": true}).
		ToSql()

	if err != nil {
		s.logger.Error("Error while building the query", "error", err)
		return nil, err
	}

	rows, err := s.postgres.QueryContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Error while executing the query", "error", err)
		return nil, err
	}
	defer rows.Close()

	var following []*pb.User
	for rows.Next() {
		var user pb.User
		var createdAt, updatedAt time.Time
		var phoneNumber sql.NullString // Using sql.NullString to handle NULL values

		err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&phoneNumber, // Scanning phone_number into sql.NullString
			&user.Bio,
			&user.ProfilePictureUrl,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			s.logger.Error("Error while scanning row", "error", err)
			return nil, err
		}

		// Set PhoneNumber to empty string if it's NULL
		if phoneNumber.Valid {
			user.PhoneNumber = phoneNumber.String
		} else {
			user.PhoneNumber = ""
		}

		user.CreatedAt = timestamppb.New(createdAt)
		user.UpdatedAt = timestamppb.New(updatedAt)

		following = append(following, &user)
	}

	if err := rows.Err(); err != nil {
		s.logger.Error("Error after iterating through rows", "error", err)
		return nil, err
	}

	return &pb.GetFollowingResponse{Following: following}, nil
}
