package auth

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ssoy1 "github.com/stanislavCasciuc/protos_grpc_auth_go/gen/go/sso"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(ctx context.Context,
		email string,
		password string,
	) (userId int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

const (
	emptyValue = 0
)

type serverAPI struct {
	ssoy1.UnimplementedAuthServer
	auth Auth
}

func Login(gRPC *grpc.Server, auth Auth) {
	ssoy1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssoy1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func isAdmin(gRPC *grpc.Server) {
	ssoy1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssoy1.IsAdminRequest) (*ssoy1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}
	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		// TODO ....
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssoy1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func (s *serverAPI) Login(ctx context.Context, req *ssoy1.LoginRequest) (*ssoy1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		// TODO ...
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssoy1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssoy1.RegisterRequest) (*ssoy1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		// TODO ...
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssoy1.RegisterResponse{
		UserId: userID,
	}, nil
}

func validateLogin(req *ssoy1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is null")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is null")
	}
	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is null")

	}

	return nil
}

func validateRegister(req *ssoy1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is null")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is null")
	}
	return nil
}

func validateIsAdmin(req *ssoy1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}
