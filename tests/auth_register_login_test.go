package tests

import (
	"auth_grpc/tests/suite"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	ssoy1 "github.com/stanislavCasciuc/protos_grpc_auth_go/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	emptyAppId = 0
	appId      = 1
	appSecret  = "secret1"

	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()

	password := randomFakePassword()

	respReg, err := st.AuthClient.Register(
		ctx, &ssoy1.RegisterRequest{
			Email:    email,
			Password: password,
		},
	)
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(
		ctx, &ssoy1.LoginRequest{
			Email:    email,
			Password: password,
			AppId:    appId,
		},
	)

	loginTime := time.Now()

	require.NoError(t, err)
	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(
		token, func(token *jwt.Token) (interface{}, error) {
			return []byte(appSecret), nil
		},
	)
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appId, int(claims["app_id"].(float64)))
	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))

	const deltaSeconds = 10

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
