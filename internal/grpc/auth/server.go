package auth

import (
	"context"
	"google.golang.org/grpc"

	ssoy1 "github.com/stanislavCasciuc/protos_grpc_auth_go/gen/go/sso"
)

type serverAPI struct {
	ssoy1.UnimplementedAuthServer
}

func Login(gRPC *grpc.Server) {
	ssoy1.RegisterAuthServer(gRPC, &serverAPI{})
}

func Register(gRPC *grpc.Server) {
	ssoy1.RegisterAuthServer(gRPC, &serverAPI{})
}

func isAdmin(gRPC *grpc.Server) {
	ssoy1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssoy1.IsAdminRequest) (*ssoy1.IsAdminResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Login(ctx context.Context, req *ssoy1.LoginRequest) (*ssoy1.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Register(ctx context.Context, req *ssoy1.RegisterRequest) (*ssoy1.RegisterResponse, error) {
	panic("implement me")
}
