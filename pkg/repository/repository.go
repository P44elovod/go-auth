package repository

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-redis/redis/v8"
	"github.com/p44elovod/auth-with-gopg/models"
)

const (
	usersTable = "users"
)

type PGAuthorization interface {
	CreateUser(user models.User) (int, error)
	GetUserID(username, password string) (int, error)
	GetUserByID(id int) (models.User, error)
}

type RedisAuthorization interface {
	SetValue(refreshToken, authToken string) error
	GetValue(refreshToken string) (string, error)
	DelValue(rerefreshToken string) error
}

type Repository struct {
	PGAuthorization
	RedisAuthorization
}

func NewRepository(db *pg.DB, rc *redis.Client) *Repository {
	return &Repository{
		PGAuthorization:    NewAuthPostgres(db),
		RedisAuthorization: NewAuthRedis(rc),
	}
}
