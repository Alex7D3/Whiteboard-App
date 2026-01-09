package service

import (
	"time"
	"drawing-api/internal/model"
	"drawing-api/internal/api"
	"github.com/golang-jwt/jwt/v5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
)

type JwtClaims struct {
	ID 		 int64  `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type TokenService struct {
	jwtSecret     []byte

}

func NewTokenService(jwtSecret []byte) *TokenService {
	return &TokenService{
		jwtSecret,
	}
}

func (s *TokenService) MakeJwtToken(ttl time.Duration, user *model.User) (string, error) {
	claims := &JwtClaims{
		ID: user.ID,
		Username: user.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Email,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *TokenService) VerifyJwtToken(signedString string) (*JwtClaims, error) {
	extractedClaims := &JwtClaims{}
	token, err := jwt.ParseWithClaims(signedString, extractedClaims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}
		return s.jwtSecret, nil
	})
	if err == jwt.ErrSignatureInvalid || !token.Valid {
		return nil, api.NewAPIError("Invalid signature", http.StatusUnauthorized)	
	} else if err != nil {
		return nil, fmt.Errorf("error parsing token %w", err)
	}

	return extractedClaims, nil
}

func (s *TokenService) RefreshJwtToken(ttl time.Duration, claims *JwtClaims) (string, error) {
	if time.Now().Before(claims.ExpiresAt.Time) {
		return "", api.NewAPIError("Token has not expired", http.StatusBadRequest)
	}
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(ttl))
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return newToken.SignedString(s.jwtSecret)
}

func (s *TokenService) MakeOpaqueToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("Failed to generate bytes %v", err)
	}
	return hex.EncodeToString(b), nil
}
