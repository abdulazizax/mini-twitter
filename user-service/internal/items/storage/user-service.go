package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	pb "github.com/abdulazizax/mini-twitter/user-service/genproto/user"
	jwttokens "github.com/abdulazizax/mini-twitter/user-service/internal/items/jwt"
	"github.com/abdulazizax/mini-twitter/user-service/internal/items/redisservice"
	"github.com/abdulazizax/mini-twitter/user-service/internal/items/repository"
	"github.com/abdulazizax/mini-twitter/user-service/internal/pkg/config"
	"google.golang.org/protobuf/types/known/timestamppb"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Storage struct {
	redisService *redisservice.RedisService
	postgres     *sql.DB
	queryBuilder sq.StatementBuilderType
	cfg          *config.Config
	logger       *slog.Logger
}

func New(redisService *redisservice.RedisService, postgres *sql.DB, queryBuilder sq.StatementBuilderType, cfg *config.Config, logger *slog.Logger) repository.IAuthRepo {
	return &Storage{
		redisService: redisService,
		postgres:     postgres,
		queryBuilder: queryBuilder,
		cfg:          cfg,
		logger:       logger,
	}
}

func (s *Storage) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	id := uuid.New().String()
	created_at := time.Now()

	password_hash, err := hashPassword(in.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	query, args, err := s.queryBuilder.Insert("users").
		Columns(
			"id",
			"email",
			"username",
			"password_hash",
			"created_at",
			"updated_at",
		).Values(
		id,
		in.Email,
		in.Username,
		password_hash,
		created_at,
		created_at,
	).ToSql()
	if err != nil {
		s.logger.Error("Failed to build query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	_, err = s.postgres.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Failed to execute query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &pb.RegisterUserResponse{
		UserId: id,
	}, nil
}

func (s *Storage) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	tx, err := s.postgres.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("Failed to start transaction", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	query, args, err := s.queryBuilder.Select(
		"id",
		"role",
		"username",
		"password_hash").
		From("users").
		Where(sq.Eq{"email": in.Email}).
		ToSql()
	if err != nil {
		s.logger.Error("Failed to build query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var id, role, username, hashedPassword string

	err = s.postgres.QueryRowContext(ctx, query, args...).Scan(&id, &role, &username, &hashedPassword)
	if err != nil {
		s.logger.Error("Failed to execute query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	_, err = checkPassword(hashedPassword, in.Password)
	if err != nil {
		s.logger.Error("Failed to check password", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to check password: %w", err)
	}

	accessToken, err := jwttokens.GenerateAccessToken(id, username, in.Email, role, s.cfg.JWT.SecretKey)
	if err != nil {
		s.logger.Error("Failed to generate access token", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := jwttokens.GenerateRefreshToken(id, username, in.Email, role, s.cfg.JWT.SecretKey)
	if err != nil {
		s.logger.Error("Failed to generate refresh token", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	query, args, err = s.queryBuilder.Update("users").
		Set("refresh_token", refreshToken).
		Set("is_active", true).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		s.logger.Error("Failed to build query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Failed to execute query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Storage) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	tx, err := s.postgres.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("Failed to start transaction", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	query, args, err := s.queryBuilder.Update("users").
		Set("refresh_token", "").
		Set("is_active", false).
		Where(sq.Eq{"id": in.UserId}).
		ToSql()
	if err != nil {
		s.logger.Error("Failed to build query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Failed to execute query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &pb.LogoutResponse{
		Success: true,
	}, nil
}

func (s *Storage) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.UserResponse, error) {
	query, args, err := s.queryBuilder.Select(
		"email",
		"username",
		"first_name",
		"last_name",
		"phone_number",
		"bio",
		"profile_picture",
		"created_at",
		"updated_at",
	).From("users").
		Where(sq.Eq{"id": in.UserId}).
		ToSql()
	if err != nil {
		s.logger.Error("Failed to build query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var email, username, first_name, last_name, phone_number, bio, profile_picture sql.NullString
	var created_at, updated_at time.Time

	err = s.postgres.QueryRowContext(ctx, query, args...).Scan(&email, &username, &first_name, &last_name, &phone_number, &bio, &profile_picture, &created_at, &updated_at)
	if err != nil {
		s.logger.Error("Failed to execute query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &pb.UserResponse{
		User: &pb.User{
			Id:                in.UserId,
			Email:             email.String,
			Username:          username.String,
			FirstName:         first_name.String,
			LastName:          last_name.String,
			PhoneNumber:       phone_number.String,
			Bio:               bio.String,
			ProfilePictureUrl: profile_picture.String,
			CreatedAt:         timestamppb.New(created_at),
			UpdatedAt:         timestamppb.New(updated_at),
		},
	}, nil
}

func (s *Storage) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	updatedAt := time.Now()

	// Dinamik yangilanishlar ro'yxati
	queryBuilder := s.queryBuilder.Update("users")

	// Bo'sh bo'lmagan qiymatlarni yangilash uchun qo'shamiz
	fieldsUpdated := false

	if in.Email != "" {
		queryBuilder = queryBuilder.Set("email", in.Email)
		fieldsUpdated = true
	}
	if in.Username != "" {
		queryBuilder = queryBuilder.Set("username", in.Username)
		fieldsUpdated = true
	}
	if in.PhoneNumber != "" {
		queryBuilder = queryBuilder.Set("phone_number", in.Username)
		fieldsUpdated = true
	}
	if in.FirstName != "" {
		queryBuilder = queryBuilder.Set("first_name", in.FirstName)
		fieldsUpdated = true
	}
	if in.LastName != "" {
		queryBuilder = queryBuilder.Set("last_name", in.LastName)
		fieldsUpdated = true
	}
	if in.Bio != "" {
		queryBuilder = queryBuilder.Set("bio", in.Bio)
		fieldsUpdated = true
	}
	if in.ProfilePictureUrl != "" {
		queryBuilder = queryBuilder.Set("profile_picture", in.ProfilePictureUrl)
		fieldsUpdated = true
	}

	// Agar hech qanday yangilanish bo'lmasa, xatolik qaytaramiz
	if !fieldsUpdated {
		// Hech qanday yangilanish bo'lmasa, foydalanuvchi ID bilan qaytamiz
		return &pb.UpdateUserResponse{
			User: &pb.User{
				Id: in.UserId,
			},
		}, nil
	}

	// 'updated_at' maydonini qo'shamiz
	queryBuilder = queryBuilder.Set("updated_at", updatedAt).
		Where(sq.Eq{"id": in.UserId})

	// SQL so'rovini yaratamiz
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		s.logger.Error("Failed to build update query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to build update query: %w", err)
	}

	// SQL so'rovini bajarish
	_, err = s.postgres.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Failed to execute update query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to execute update query: %w", err)
	}

	// Yangilangan foydalanuvchi ma'lumotlarini qaytarish
	return &pb.UpdateUserResponse{
		User: &pb.User{
			Id:                in.UserId,
			Email:             in.Email,
			Username:          in.Username,
			PhoneNumber:       in.PhoneNumber,
			FirstName:         in.FirstName,
			LastName:          in.LastName,
			Bio:               in.Bio,
			ProfilePictureUrl: in.ProfilePictureUrl,
			UpdatedAt:         timestamppb.New(updatedAt),
		},
	}, nil
}

func (s *Storage) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	deleted_at := time.Now()

	query, args, err := s.queryBuilder.Update("users").
		Set("deleted_at", deleted_at).
		Where(sq.Eq{"id": in.UserId}).
		ToSql()
	if err != nil {
		s.logger.Error("Failed to build query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	_, err = s.postgres.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Failed to execute query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, in *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	query, args, err := s.queryBuilder.Select("id").
		From("users").
		Where(sq.Eq{"email": in.Email}).
		ToSql()
	if err != nil {
		s.logger.Error("Failed to build query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var id string

	err = s.postgres.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		s.logger.Error("Failed to execute query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &pb.GetUserByEmailResponse{
		UserId: id,
	}, nil
}

func (s *Storage) GetUserByUsername(ctx context.Context, in *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error) {
	query, args, err := s.queryBuilder.Select("id").
		From("users").
		Where(sq.Eq{"username": in.Username}).
		ToSql()
	if err != nil {
		s.logger.Error("Failed to build query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var id string

	err = s.postgres.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		s.logger.Error("Failed to execute query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &pb.GetUserByUsernameResponse{
		UserId: id,
	}, nil
}

func (s *Storage) UpdateUserPassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (*pb.RawResponse, error) {
	err := s.verifyEmail(context.Background(), req.Email, int(req.VerificationCode))
	if err != nil {
		s.logger.Error("Failed to verify email", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to verify email: %w", err)
	}

	tx, err := s.postgres.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("Failed to start transaction", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	password_hash, err := hashPassword(req.NewPassword)
	if err != nil {
		s.logger.Error("Failed to hash password", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	query, args, err := s.queryBuilder.Update("users").
		Set("password_hash", password_hash).
		Where(sq.Eq{"email": req.Email}).
		ToSql()
	if err != nil {
		s.logger.Error("Failed to build query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Failed to execute query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &pb.RawResponse{
		Message: "Password updated successfully",
	}, nil
}

func (s *Storage) SendVerificationCode(ctx context.Context, req *pb.SendVerificationCodeRequest) (*pb.RawResponse, error) {
	query, args, err := s.queryBuilder.Select("email").
		From("users").
		Where(sq.Eq{"email": req.Email}).
		ToSql()
	if err != nil {
		s.logger.Error("Failed to build query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var email string
	err = s.postgres.QueryRowContext(ctx, query, args...).Scan(&email)
	if err != nil {
		s.logger.Error("Failed to execute query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	if email != req.Email {
		s.logger.Info("Email not found", slog.String("email", req.Email))
		return nil, fmt.Errorf("email not found: %s", req.Email)
	}

	if err := s.sendVerificationCode(context.Background(), req.Email); err != nil {
		s.logger.Error("Failed to send verification code", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to send verification code: %w", err)
	}

	return &pb.RawResponse{
		Message: "Verification code sent to your email!",
	}, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func checkPassword(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, fmt.Errorf("failed to compare password: %w", err)
	}
	return true, nil
}
