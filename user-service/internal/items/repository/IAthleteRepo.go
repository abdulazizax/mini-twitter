package repository

import (
	"context"

	pb "github.com/abdulazizax/mini-twitter/user-service/genproto/user"
)

type IAuthRepo interface {
	// 1
	RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error)
	// 2
	GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.UserResponse, error)
	// 3
	UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)
	// 4
	DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error)

	// 5
	FollowUser(ctx context.Context, in *pb.FollowUserRequest) (*pb.FollowUserResponse, error)
	// 6
	UnfollowUser(ctx context.Context, in *pb.UnfollowUserRequest) (*pb.UnfollowUserResponse, error)
	// 7
	GetFollowers(ctx context.Context, in *pb.GetFollowersRequest) (*pb.GetFollowersResponse, error)
	// 8
	GetFollowing(ctx context.Context, in *pb.GetFollowingRequest) (*pb.GetFollowingResponse, error)

	// 9
	Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error)
	// 10
	Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutResponse, error)
	// 11
	GetUserByEmail(ctx context.Context, in *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error)
	// 12
	GetUserByUsername(ctx context.Context, in *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error)
	// 13
	UpdateUserPassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (*pb.RawResponse, error)
	// 14
	SendVerificationCode(ctx context.Context, req *pb.SendVerificationCodeRequest) (*pb.RawResponse, error)
}
