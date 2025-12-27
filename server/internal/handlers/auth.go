package handlers

import (
	"net/http"
	"time"
	"fmt"
	"context"
	"drawing-api/internal/storage"
	"drawing-api/internal/service"
	"drawing-api/internal/model"
	"drawing-api/internal/util"
	"drawing-api/internal/api"
)

type AuthHandler struct {
	userStorage   storage.UserStorage
	tokenService  *service.TokenService
	cookieService *service.CookieService
	timeout       time.Duration
}

func NewAuthHandler(
	userStorage storage.UserStorage,
	tokenService *service.TokenService,
	cookieService *service.CookieService,
	timeout time.Duration,
) *AuthHandler {
	return  &AuthHandler {
		userStorage,
		tokenService,
		cookieService,
		timeout,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()
	var userReq model.UserRequest
	if err := api.ParseJSON(r, &userReq); err != nil {
		return err
	}
	_, err := h.userStorage.GetByEmail(ctx, userReq.Email)
	if err == nil {
		return api.NewAPIError(
			fmt.Sprintf("user with email '%s' already exists", userReq.Email),
			http.StatusConflict,
		)
	}
	hash, err := util.HashPassword(userReq.Password)
	if err != nil {
		return err
	}
	user := model.NewUser(userReq.Username, userReq.Email, hash)
	id, err := h.userStorage.Create(ctx, user)
	if err != nil {
		return err
	}
	userRes := model.NewUserResponse(id, user.UserName, user.Email)
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
	ss, err := h.tokenService.GetSignedString(user)
	if err != nil {
		return err
	}

	h.cookieService.SetAuthCookie(w, ss)
	return api.WriteJSON(w, http.StatusOK, model.NewUserResponse(
		user.ID,
		user.UserName,
		user.Email,
	))
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	h.cookieService.RemoveAuthCookie(w)
	api.WriteJSONMessage(w, http.StatusOK, "Logged out")
	return nil
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) error {
	ss, err := h.cookieService.ExtractSignedString(r)
	if err != nil {
		return err
	}
	extractedClaims, err := h.tokenService.VerifyToken(ss)
	if err != nil {
		return err
	}
	refreshedSS, err := h.tokenService.RefreshString(extractedClaims)
	h.cookieService.SetAuthCookie(w, refreshedSS)
	api.WriteJSONMessage(w, http.StatusOK, "Refreshed")
	return nil
}
