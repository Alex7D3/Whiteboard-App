package service

import (
	"time"
	"net/http"
	"drawing-api/internal/api"
)

type CookieService struct {
	cookieName    string
	expireMinutes time.Duration
}

func NewCookieService(cookieName string, expireMinutes time.Duration) *CookieService {
	return &CookieService {
		cookieName,
		expireMinutes,
	}
}

func (s *CookieService) SetAuthCookie(w http.ResponseWriter, signedString string) {
	http.SetCookie(w, &http.Cookie{
		Name: s.cookieName,
		Value: signedString,
		Path: "/",
		Expires: time.Now().Add(s.expireMinutes * time.Minute),
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
