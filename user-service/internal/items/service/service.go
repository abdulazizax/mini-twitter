package service

import (
	"context"
	"log/slog"

	pb "github.com/abdulazizax/mini-twitter/user-service/genproto/user"

	"github.com/abdulazizax/mini-twitter/user-service/internal/items/repository"
)

type Service struct {
	pb.UnimplementedUserServiceServer
	storage repository.IAuthRepo
	logger  *slog.Logger
}

func New(storage repository.IAuthRepo, logger *slog.Logger) *Service {
	return &Service{
		storage: storage,
		logger:  logger,
	}
}

// 1 ------
func (s *Service) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	s.logger.Info("Register function was invoked", slog.String("request", in.String()))
	return s.storage.RegisterUser(ctx, in)
}

// 2
func (s *Service) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.UserResponse, error) {
	s.logger.Info("GetUser function was invoked", slog.String("request", in.String()))
	return s.storage.GetUser(ctx, in)
}

// 3 ------
func (s *Service) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	s.logger.Info("UpdateUser function was invoked", slog.String("request", in.String()))
	return s.storage.UpdateUser(ctx, in)
}

// 4 ------
func (s *Service) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	s.logger.Info("DeleteUser function was invoked", slog.String("request", in.String()))
	return s.storage.DeleteUser(ctx, in)
}

// 5 ------
func (s *Service) FollowUser(ctx context.Context, in *pb.FollowUserRequest) (*pb.FollowUserResponse, error) {
	s.logger.Info("FollowUser function was invoked", slog.String("request", in.String()))
	return s.storage.FollowUser(ctx, in)
}

// 6 ------
func (s *Service) UnfollowUser(ctx context.Context, in *pb.UnfollowUserRequest) (*pb.UnfollowUserResponse, error) {
	s.logger.Info("UnfollowUser function was invoked", slog.String("request", in.String()))
	return s.storage.UnfollowUser(ctx, in)
}

// 7
func (s *Service) GetFollowers(ctx context.Context, in *pb.GetFollowersRequest) (*pb.GetFollowersResponse, error) {
	s.logger.Info("GetFollowers function was invoked", slog.String("request", in.String()))
	return s.storage.GetFollowers(ctx, in)
}

// 8
func (s *Service) GetFollowing(ctx context.Context, in *pb.GetFollowingRequest) (*pb.GetFollowingResponse, error) {
	s.logger.Info("GetFollowing function was invoked", slog.String("request", in.String()))
	return s.storage.GetFollowing(ctx, in)
}

// 9
func (s *Service) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	s.logger.Info("Login function was invoked", slog.String("request", in.String()))
	return s.storage.Login(ctx, in)
}

// 10
func (s *Service) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	s.logger.Info("Logout function was invoked", slog.String("request", in.String()))
	return s.storage.Logout(ctx, in)
}

// 11
func (s *Service) GetUserByEmail(ctx context.Context, in *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	s.logger.Info("GetUserByEmail function was invoked", slog.String("request", in.String()))
	return s.storage.GetUserByEmail(ctx, in)
}

// 12
func (s *Service) GetUserByUsername(ctx context.Context, in *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error) {
	s.logger.Info("GetUserByUsername function was invoked", slog.String("request", in.String()))
	return s.storage.GetUserByUsername(ctx, in)
}

// 13 ------
func (s *Service) UpdateUserPassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (*pb.RawResponse, error) {
	s.logger.Info("UpdatePasswordService function was invoked", slog.String("request", req.String()))
	return s.storage.UpdateUserPassword(ctx, req)
}

// 14 ------
func (s *Service) SendVerificationCode(ctx context.Context, req *pb.SendVerificationCodeRequest) (*pb.RawResponse, error) {
	s.logger.Info("SendVerificationCode function was invoked", slog.String("request", req.String()))
	return s.storage.SendVerificationCode(ctx, req)
}
