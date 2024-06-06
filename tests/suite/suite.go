package suite

import (
	"auth_grpc/internal/config"
	"context"
	ssoy1 "github.com/stanislavCasciuc/protos_grpc_auth_go/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
)

const grpcHost = "localhost"

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssoy1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/local.yaml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(
		func() {
			t.Helper()
			cancelCtx()
		},
	)

	cc, err := grpc.DialContext(
		context.Background(), grpcAddress(cfg), grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}
	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssoy1.NewAuthClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {

	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
