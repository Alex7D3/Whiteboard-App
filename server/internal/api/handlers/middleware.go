package handlers

import (
	"net/http"
	"time"
	"fmt"
	"context"
	"encoding/json"
	"drawing-api/internal/storage"
	"drawing-api/internal/model"
	"drawing-api/internal/util"
	"drawing-api/internal/api"
	"github.com/golang-jwt/jwt/v5"
)

func (h *AuthHandler) Authorize(next api.AppHandler) api.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		cookie, err := r.Cookie("auth_token")
		if err == http.ErrNoCookie {
			return api.NewAPIError("Missing cookie", http.StatusUnauthorized)
		}
		if err != nil {
			return api.NewAPIError("Invalid cookie", http.StatusBadRequest)
		}
		tokenStr := cookie.Value
		claims := JwtClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
			return h.jwtSecret, nil
		})
		return next(w, r)
	}
}
