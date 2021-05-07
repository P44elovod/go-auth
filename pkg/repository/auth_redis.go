package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type AuthRedis struct {
	rc *redis.Client
}

func NewAuthRedis(rc *redis.Client) *AuthRedis {
	return &AuthRedis{rc: rc}
}

func (r *AuthRedis) SetValue(refreshToken, authToken string) error {
	err := r.rc.Set(r.rc.Context(), refreshToken, authToken, time.Duration(10*time.Minute))
	if err != nil {
		redisErr := fmt.Sprintf("error accure at inputing value to redis: %s", err)
		return errors.New(redisErr)
	}
	return nil
}

func (r *AuthRedis) GetValue(refreshToken string) (string, error) {
	val, err := r.rc.Get(r.rc.Context(), refreshToken).Result()
	if err != nil {
		redisErr := fmt.Sprintf("error accure at getting value to redis: %s", err)
		return "", errors.New(redisErr)
	}
	return val, nil
}

func (r *AuthRedis) DelValue(rerefreshToken string) error {
	err := r.rc.Del(r.rc.Context(), rerefreshToken)
	if err != nil {
		redisErr := fmt.Sprintf("error accure at deleting value to redis: %s", err)
		return errors.New(redisErr)
	}

	return nil
}
