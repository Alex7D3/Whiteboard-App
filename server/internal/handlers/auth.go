package handlers

import (
	"net/http"
	"time"
	"fmt"
	"context"
	"strings"
	"drawing-api/internal/storage"
	"drawing-api/internal/service"
	"drawing-api/internal/model"
	"drawing-api/internal/util"
	"drawing-api/internal/api"
)

type AuthHandler struct {
	userStorage    storage.UserStorage
	sessionStorage storage.SessionStorage
	tokenService   *service.TokenService
	cookieService  *service.CookieService
	timeout        time.Duration
	accessExpiry   time.Duration
	refreshExpiry  time.Duration
}

func NewAuthHandler(
	userStorage storage.UserStorage,
	sessionStorage storage.SessionStorage,
	tokenService *service.TokenService,
	cookieService *service.CookieService,
	timeout, accessExpiry, refreshExpiry time.Duration,
) *AuthHandler {
	return  &AuthHandler {
		userStorage,
		sessionStorage,
		tokenService,
		cookieService,
		timeout,
		accessExpiry,
		refreshExpiry,
	}
}

func (h *AuthHandler) extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", api.NewAPIError("missing auth header", http.StatusUnauthorized)
	}

	prefix := "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", api.NewAPIError("missing Bearer prefix", http.StatusUnauthorized)
	}

	token := strings.TrimPrefix(authHeader, prefix)
	token = strings.Trim(token, " ")
	if token == "" {
		return "", api.NewAPIError("missing bearer token", http.StatusUnauthorized)
	}

	return token, nil
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()
	var userReq model.UserRequest
	if err := api.ParseJSON(r, &userReq); err != nil {
		return err
	}
	if userReq.Email == "" || userReq.Username == "" {
        return api.NewAPIError("email and username are required", http.StatusBadRequest)
    }
    if len(userReq.Password) < 6 {
        return api.NewAPIError("password must be at least 6 characters", http.StatusBadRequest)
    }
	_, err := h.userStorage.GetByEmail(ctx, userReq.Email)
	if err == nil {
		return api.NewAPIError(
			fmt.Sprintf("user with email '%s' already exists", userReq.Email),
			http.StatusConflict,
		)
	}
	passwordHash, err := util.HashPassword(userReq.Password)
	if err != nil {
		return err
	}
	user := &model.User{Username: userReq.Username, Email: userReq.Email, PasswordHash: passwordHash}
	id, err := h.userStorage.Create(ctx, user)
	if err != nil {
		return err
	}
	userRes := &model.UserResponse{ID: id, Username: user.Username, Email: user.Email}
	return api.WriteJSON(w, http.StatusOK, userRes)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	var userReq model.UserRequest

	if err := api.ParseJSON(r, &userReq); err != nil {
		return api.NewAPIError("Invalid message body", http.StatusBadRequest)
	}

	user, err := h.userStorage.GetByEmail(ctx, userReq.Email)
	if err != nil {
		msg := fmt.Sprintf("No account with email '%s' found", userReq.Email)
		return api.NewAPIError(msg, http.StatusNotFound)
	}

	if !util.CheckPasswordHash(userReq.Password, user.PasswordHash) {
		return api.NewAPIError("Incorrect password", http.StatusUnauthorized)
	}

	accessTok, err := h.tokenService.MakeJwtToken(h.accessExpiry, user)
	if err != nil {
		return err
	}

	refreshTok, err := h.tokenService.MakeOpaqueToken()
	if err != nil {
		return err
	}

	refreshTokHash := util.HashToken(refreshTok)

	if _, err = h.sessionStorage.Create(ctx, &model.Session{
		TokenHash: refreshTokHash,
		UserID: user.ID,
		ExpiresAt: time.Now().Add(h.refreshExpiry),
	}); err != nil {
		return err
	}

	h.cookieService.SetAuthCookie(w, h.refreshExpiry, refreshTok)

	return api.WriteJSON(w, http.StatusOK, &model.LoginResponse{
		User: user,
		AccessToken: accessTok,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	refreshTok, err := h.cookieService.ExtractSignedString(r)
	if err != nil {
		h.sessionStorage.Revoke(ctx, refreshTok)
		h.cookieService.RemoveAuthCookie(w)
	}
	return api.WriteJSONMessage(w, http.StatusOK, "Logged out")
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	oldRefreshTok, err := h.cookieService.ExtractSignedString(r)
	if err != nil {
		return err
	}

	newRefreshTok, err := h.tokenService.MakeOpaqueToken()
	if err != nil {
		return err
	}

	session, err := h.sessionStorage.RotateToken(ctx, h.refreshExpiry, oldRefreshTok, newRefreshTok)
	if err != nil {
		return err
	}

	user, err := h.userStorage.GetByID(ctx, session.UserID)
	if err != nil {
		return err
	}

	newAccessTok, err := h.tokenService.MakeJwtToken(h.accessExpiry, user)
	if err != nil {
		return err
	}

	h.cookieService.SetAuthCookie(w, h.refreshExpiry, newRefreshTok)
	return api.WriteJSON(w, http.StatusOK, &model.RefreshResponse{AccessToken: newAccessTok})
}
