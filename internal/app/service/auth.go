package service

import (
	"casino-back/internal/app/logger"
	"crypto/sha1"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AuthService struct {
	authRepository IUserRepository
}

func NewAuthService(authRepository IUserRepository) *AuthService {
	return &AuthService{authRepository: authRepository}
}

func (s *AuthService) CreateUser(name, login, password string) (int, error) {
	password = s.generatePasswordHash(password)
	return s.authRepository.CreateUser(name, login, password)
}

func (s *AuthService) GetUserId(login, password string) (int, error) {
	password = s.generatePasswordHash(password)
	return s.authRepository.GetUserId(login, password)
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

type tokenClaims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

const (
	signingKey = "secret_key"
)

func (s *AuthService) GenerateToken(login, password string) (string, error) {
	id, err := s.authRepository.GetUserId(login, s.generatePasswordHash(password))

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.InfoKV("service err", "err", fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		logger.InfoKV("service error", "err", "token claims are not of type *tokenClaims")
		return 0, fmt.Errorf("invalid token or token claims")
	}

	return claims.UserId, nil
}
