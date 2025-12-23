package handlers

import (
	"net/http"
	"time"
	"fmt"
	"context"
	"drawing-api/internal/storage"
	"drawing-api/internal/model"
	"drawing-api/internal/util"
	"drawing-api/internal/api"
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthHandler struct {
	userStorage storage.UserStorage
	jwtSecret []byte
	claims JwtClaims
	timeout time.Duration
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
	expTime := time.Now().Add(5 * time.Minute)
	claims := &JwtClaims{
		Username: userReq.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedStr, err := token.SignedString(h.jwtSecret)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name: "auth_token",
		Value: signedStr,
		Expires: expTime,
	})
	return api.WriteJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r http.Request) error {
	return nil
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r http.Request) error {
	ctx := r.Context()
	id, err := h.userStorage.GetById(ctx, ctx.Value("id").(int))
	if id != nil {
		return err
	}
	return nil
}
