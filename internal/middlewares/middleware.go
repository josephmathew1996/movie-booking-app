package middlewares

import (
	"movie-booking-app/users-service/pkg/utils"
	"strings"
)

//go:generate mockery --name IAuthMiddleware --case underscore
type IAuth interface {
	Authenticate(string) (*utils.Claims, error)
}

type Auth struct{}

func NewAuth() IAuth {
	return Auth{}
}

func (a Auth) Authenticate(token string) (*utils.Claims, error) {
	tokenStr := strings.TrimPrefix(token, "Bearer ")
	claims, err := utils.ParseToken(tokenStr)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
