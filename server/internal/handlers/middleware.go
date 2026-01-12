package handlers

import (
	"net/http"
	"drawing-api/internal/api"
	"context"
	"time"
	"fmt"
)

const ClaimsKey string = "user_claims"

func (h *AuthHandler) Authorize(next api.AppHandler) api.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		accessTok, err := h.extractBearerToken(r)
		if err != nil {
			return err
		}

		claims, err := h.tokenService.VerifyJwtToken(accessTok)
		if err != nil {
			return api.NewAPIError("Invalid token", http.StatusUnauthorized)
		}

		if time.Now().After(claims.ExpiresAt.Time) {
			return api.NewAPIError("Token has expired", http.StatusUnauthorized)
		}

        ctx := context.WithValue(r.Context(), ClaimsKey, claims)
        return next(w, r.WithContext(ctx))
	}
}

// Websocket JS API does not allow custom headers, so send as a query param instead
func (h *AuthHandler) AuthorizeWS(next api.AppHandler) api.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		accessTok := r.URL.Query().Get("access_token")
		if accessTok == "" {
			return api.NewAPIError("Missing access_token parameter", http.StatusBadRequest)
		}

		claims, err := h.tokenService.VerifyJwtToken(accessTok)
		if err != nil {
			return api.NewAPIError("Invalid token", http.StatusUnauthorized)
		}

		if time.Now().After(claims.ExpiresAt.Time) {
			return api.NewAPIError("Token has expired", http.StatusUnauthorized)
		}

        ctx := context.WithValue(r.Context(), ClaimsKey, claims)
        return next(w, r.WithContext(ctx))
	}
}
