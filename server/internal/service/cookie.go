package service

import (
	"time"
	"net/http"
	"drawing-api/internal/api"
)

type CookieService struct {
	cookieName string
}

func NewCookieService(cookieName string) *CookieService {
	return &CookieService {
		cookieName,
	}
}

func (s *CookieService) SetAuthCookie(w http.ResponseWriter, expiry time.Duration, signedString string) {
	http.SetCookie(w, &http.Cookie{
		Name: s.cookieName,
		Value: signedString,
		Path: "/api/refresh",
		Expires: time.Now().Add(expiry),
		HttpOnly: true,
	})
}

func (s *CookieService) RemoveAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name: s.cookieName,
		Value: "",
		Path: "/",
		MaxAge: -1,
	});
}

func (s *CookieService) ExtractSignedString(r *http.Request) (string, error) {
		cookie, err := r.Cookie(s.cookieName)
		if err == http.ErrNoCookie {
			return "", api.NewAPIError("Missing cookie", http.StatusUnauthorized)
		}
		if err != nil {
			return "", api.NewAPIError("Invalid cookie", http.StatusBadRequest)
		}
		return cookie.Value, nil
}
