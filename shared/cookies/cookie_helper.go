package cookies

import (
	"net/http"
	"time"
)

type CookieConfig struct {
	Name     string
	Value    string
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
}

func SetSecureCookie(w http.ResponseWriter, config CookieConfig) {
	cookie := &http.Cookie{
		Name:     config.Name,
		Value:    config.Value,
		Path:     config.Path,
		Domain:   config.Domain,
		MaxAge:   config.MaxAge,
		Secure:   config.Secure,
		HttpOnly: config.HttpOnly,
		SameSite: config.SameSite,
	}

	if config.MaxAge > 0 {
		cookie.Expires = time.Now().Add(time.Duration(config.MaxAge) * time.Second)
	}

	http.SetCookie(w, cookie)
}

// ClearCookie удаляет cookie
func ClearCookie(w http.ResponseWriter, name, path string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     path,
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// GetSessionCookie извлекает сессию из cookie
func GetSessionCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// GetCSRFCookie извлекает CSRF токен из cookie
func GetCSRFCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("csrf_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
