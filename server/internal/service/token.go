package service

import (
	"time"
	"drawing-api/internal/model"
	"drawing-api/internal/api"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type JwtClaims struct {
	ID 		 int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type TokenService struct {
	jwtSecret     []byte
	expireMinutes time.Duration
}

func NewTokenService(jwtSecret []byte, expireMinutes time.Duration) *TokenService {
	return &TokenService{
		jwtSecret,
		expireMinutes,
	}
}

func (s *TokenService) GetSignedString(user *model.User) (string, error) {
	expTime := time.Now().Add(s.expireMinutes * time.Minute)
	claims := &JwtClaims{
		ID: user.ID,
		Email: user.Email,
		Username: user.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *TokenService) VerifyToken(signedString string) (*JwtClaims, error) {
	extractedClaims := &JwtClaims{}
	token, err := jwt.ParseWithClaims(signedString, extractedClaims, func(token *jwt.Token) (any, error) {
		return s.jwtSecret, nil
	})
	if err == jwt.ErrSignatureInvalid || !token.Valid {
		return nil, api.NewAPIError("Invalid signature", http.StatusUnauthorized)	
	} else if err != nil {
		return nil, err
	}

	return extractedClaims, nil
}

func (s *TokenService) RefreshString(claims *JwtClaims) (string, error) {
	if time.Now().Before(claims.ExpiresAt.Time) {
		return "", api.NewAPIError("Token has not expired", http.StatusBadRequest)
	}
	expTime := time.Now().Add(s.expireMinutes * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expTime)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return newToken.SignedString(s.jwtSecret)
}
