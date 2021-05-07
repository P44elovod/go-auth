package service

import (
	"github.com/p44elovod/auth-with-gopg/models"
	"github.com/p44elovod/auth-with-gopg/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	generateAuthToken(username, password string) (string, error)
	generateRefreshToken() string
	ParseAuthToken(accessToken string) (int, error)
	GenerateTokenPair(username, password string) (map[string]string, error)
	RefreshTokens(refreshToken string) (map[string]string, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.PGAuthorization, repos.RedisAuthorization),
	}
}
