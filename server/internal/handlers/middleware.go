package handlers

import (
	"net/http"
	"drawing-api/internal/api"
	"context"
)

const ClaimsKey string = "user_claims"


func (h *AuthHandler) Authorize(next api.AppHandler) api.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
        ss, err := h.cookieService.ExtractSignedString(r)
        if err != nil {
            return err
        }

        claims, err := h.tokenService.VerifyToken(ss)
        if err != nil {
            return err
        }

        ctx := context.WithValue(r.Context(), ClaimsKey, claims)
        return next(w, r.WithContext(ctx))
	}
}
