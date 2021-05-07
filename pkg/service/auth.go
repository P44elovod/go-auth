package service

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/p44elovod/auth-with-gopg/models"
	"github.com/p44elovod/auth-with-gopg/pkg/repository"
	uuid "github.com/satori/go.uuid"
)

var (
	signingKey = os.Getenv("SIGNING_KEY")
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo  repository.PGAuthorization
	redis repository.RedisAuthorization
}

func NewAuthService(repo repository.PGAuthorization, redis repository.RedisAuthorization) *AuthService {
	return &AuthService{
		repo:  repo,
		redis: redis,
	}
}

func (s *AuthService) CreateUser(user models.User) (int, error) {
	return s.repo.CreateUser(user)
}

func (s *AuthService) RefreshTokens(refreshToken string) (map[string]string, error) {
	tokens := make(map[string]string)

	authToken, err := s.redis.GetValue(refreshToken)
	if err != nil {
		return nil, err
	}
	userID, err := s.ParseAuthToken(authToken)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	tokens, err = s.GenerateTokenPair(user.Username, user.Password)
	if err != nil {

		return nil, err
	}

	s.redis.DelValue(refreshToken)
	return tokens, nil

}

func (s *AuthService) GenerateTokenPair(username, password string) (map[string]string, error) {

	tokens := make(map[string]string)
	authToken, err := s.generateAuthToken(username, password)
	if err != nil {

		return nil, err
	}

	refreshToken := s.generateRefreshToken()
	s.redis.SetValue(refreshToken, authToken)

	tokens["authToken"] = authToken
	tokens["refreshToken"] = refreshToken

	return tokens, nil
}

func (s *AuthService) generateAuthToken(username, password string) (string, error) {
	userID, err := s.repo.GetUserID(username, password)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userID,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseAuthToken(authToken string) (int, error) {
	token, err := jwt.ParseWithClaims(authToken, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return 0, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil

}

func (s *AuthService) generateRefreshToken() string {

	token := uuid.NewV4().String()

	return token
}
