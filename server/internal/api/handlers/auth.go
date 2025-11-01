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
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	userStorage storage.UserStorage
	jwtSecret []byte
	jwt.RegisteredClaims
	tokenExpiration time.Duration
	timeout time.Duration	
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(r.Context, h.timeout)
	defer cancel()

	var userReq model.UserRequest
	if err := ParseJSON(r, &userReq); err != nil {
		return err
	}

	_, err := h.userStorage.GetByEmail(ctx, userReq.Email)
	if err == nil {
		return NewApiError(
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
	return WriteJSON(w, http.StatusOK, userRes)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(r.Context, h.timeout)
	defer cancel()

	var userReq model.UserRequest
	if err := json.NewDecoder(r.Body).Decode(*userReq); err != nil {
		return NewAPIError("Invalid message body", http.StatusBadRequest)
	}

	var user model.User
	user, err := h.userStorage.GetByEmail(userReq.Email)
	if err != nil {
		msg := fmt.Sprintf("No account with email '%s' found", userReq.Email)
		return NewAPIError(msg, http.StatusNotFound)
	}

	if !util.CheckPasswordHash(userReq.Password, user.PasswordHash) {
		return NewAPIError("Incorrect password", http.StatusUnauthorized)
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, )

	return WriteJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r http.Request) error {

}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r http.Request) error {
	ctx := r.Context()
	id, err := h.userStorage.GetById(ctx.Value("id").(int))
	if id != nil
}
