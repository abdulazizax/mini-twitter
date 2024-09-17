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

func New(postgres *sql.DB, queryBuilder sq.StatementBuilderType, cfg *config.Config, logger *slog.Logger) repository.IAuthRepo {
	return &Storage{
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
		s.logger.Error("Error while hashing password:", slog.String("err:", err.Error()))
		return nil, err
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
		password_hash,
		created_at,
		created_at,
	).ToSql()
	if err != nil {
		s.logger.Error("Error while building a query")
		return nil, err
	}

	_, err = s.postgres.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Error while executing a query")
		return nil, err
	}

	return &pb.RegisterUserResponse{
		UserId: id,
	}, nil
}

func (s *Storage) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	tx, err := s.postgres.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("Error while starting a transaction")
		return nil, err
	}
	defer tx.Rollback()

	query, args, err := s.queryBuilder.Select(
		"id",
		"role",
		"password_hash").
		From("users").
		Where(sq.Eq{"email": in.Email}).
		ToSql()
	if err != nil {
		s.logger.Error("Error while building a query")
		return nil, err
	}

	var id, role, hashedPassword string

	err = s.postgres.QueryRowContext(ctx, query, args...).Scan(&id, &role, &hashedPassword)
	if err != nil {
		s.logger.Error("Error while executing a query")
		return nil, err
	}

	_, err = checkPassword(hashedPassword, in.Password)
	if err != nil {
		s.logger.Error("Error while checking password")
		return nil, err
	}

	accessToken, err := jwttokens.GenerateAccessToken(id, in.Email, role, s.cfg.JWT.SecretKey)
	if err != nil {
		s.logger.Error("Error while generating access token")
		return nil, err
	}

	refreshToken, err := jwttokens.GenerateRefreshToken(id, in.Email, role, s.cfg.JWT.SecretKey)
	if err != nil {
		s.logger.Error("Error while generating refresh token")
		return nil, err
	}

	query, args, err = s.queryBuilder.Update("users").
		Set("refresh_token", refreshToken).
		Set("is_active", true).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		s.logger.Error("Error while building a query")
		return nil, err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Error while executing a query")
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		s.logger.ErrorContext(ctx, "error while committing transaction", slog.String("error", err.Error()))
		return nil, err
	}

	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Storage) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	tx, err := s.postgres.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("Error while starting a transaction")
		return nil, err
	}
	defer tx.Rollback()

	query, args, err := s.queryBuilder.Update("users").
		Set("refresh_token", "").
		Set("is_active", false).
		Where(sq.Eq{"id": in.UserId}).
		ToSql()
	if err != nil {
		s.logger.Error("Error while building a query")
		return nil, err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Error while executing a query")
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		s.logger.ErrorContext(ctx, "error while committing transaction", slog.String("error", err.Error()))
		return nil, err
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
		s.logger.Error("Error while building a query")
		return nil, err
	}

	var email, username, first_name, last_name, phone_number, bio, profile_picture sql.NullString
	var created_at, updated_at time.Time

	err = s.postgres.QueryRowContext(ctx, query, args...).Scan(&email, &username, &first_name, &last_name, &phone_number, &bio, &profile_picture, &created_at, &updated_at)
	if err != nil {
		s.logger.Error("Error while executing a query")
		return nil, err
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
		s.logger.Error("Error while building update query", "error", err)
		return nil, err
	}

	// SQL so'rovini bajarish
	_, err = s.postgres.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Error while executing update query", "error", err)
		return nil, err
	}

	// Yangilangan foydalanuvchi ma'lumotlarini qaytarish
	return &pb.UpdateUserResponse{
		User: &pb.User{
			Id:                in.UserId,
			Email:             in.Email,
			Username:          in.Username,
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
		s.logger.Error("Error while building a query")
		return nil, err
	}

	_, err = s.postgres.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Error while executing a query")
		return nil, err
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
		s.logger.Error("Error while building a query")
		return nil, err
	}

	var id string

	err = s.postgres.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		s.logger.Error("Error while executing a query")
		return nil, err
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
		s.logger.Error("Error while building a query")
		return nil, err
	}

	var id string

	err = s.postgres.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		s.logger.Error("Error while executing a query")
		return nil, err
	}

	return &pb.GetUserByUsernameResponse{
		UserId: id,
	}, nil
}

func (s *Storage) UpdatePassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (*pb.RawResponse, error) {
	err := s.verifyEmail(context.Background(), req.Email, int(req.VerificationCode))
	if err != nil {
		s.logger.Error("Error while verifying email")
		return nil, err
	}

	tx, err := s.postgres.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("Error while starting a transaction")
		return nil, err
	}
	defer tx.Rollback()

	query, args, err := s.queryBuilder.Update("users").
		Set("password_hash", req.NewPassword).
		Where(sq.Eq{"id": req.UserId}).
		ToSql()
	if err != nil {
		s.logger.Error("Error while building a query")
		return nil, err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Error while executing a query")
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		s.logger.ErrorContext(ctx, "error while committing transaction", slog.String("error", err.Error()))
		return nil, err
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
		s.logger.Error("Error while building a query")
		return nil, err
	}

	var email string
	err = s.postgres.QueryRowContext(ctx, query, args...).Scan(&email)
	if err != nil {
		s.logger.Error("Error while executing a query")
		return nil, err
	}

	if email != req.Email {
		s.logger.Info("Email not found")
		return nil, fmt.Errorf("email not found")
	}

	if err := s.sendVerificationCode(context.Background(), req.Email); err != nil {
		s.logger.Error("Error while sending verification code")
		return nil, err
	}

	return &pb.RawResponse{
		Message: "Verification code sent to your email!",
	}, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPassword(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil, err
}
