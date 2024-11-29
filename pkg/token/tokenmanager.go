package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims interface {
	jwt.Claims
}

type JwtTokenManager struct {
	SecretKey string
	Lifetime  int
}

func NewJwtTokenGenerator(secretKey string, lifetime int) *JwtTokenManager {
	return &JwtTokenManager{SecretKey: secretKey, Lifetime: lifetime}
}

func (m *JwtTokenManager) Generate(claims Claims) (string, error) {
	if expTime, _ := claims.GetExpirationTime(); expTime == nil {
		claims.(*jwt.RegisteredClaims).ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(m.Lifetime)))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.SecretKey))
}

func (m *JwtTokenManager) Valid(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.SecretKey), nil
	})

	if err != nil {
		return false, err
	}

	return token.Valid, nil
}

func (m *JwtTokenManager) Decode(tokenString string, toClaims Claims) (*jwt.Token, *Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, toClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.SecretKey), nil
	})

	if err != nil {
		return nil, nil, err
	}

	if claims, ok := token.Claims.(Claims); ok && token.Valid {
		return token, &claims, nil
	}

	return nil, nil, fmt.Errorf("decode error")
}
