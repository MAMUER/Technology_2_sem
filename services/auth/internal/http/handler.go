package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"tech-ip-sem2/services/auth/internal/service"
	"tech-ip-sem2/shared/cookies"
	"tech-ip-sem2/shared/logger"
	"tech-ip-sem2/shared/middleware"
)

type Handlers struct {
	authService    *service.AuthService
	sessionService *service.SessionService
	log            *logger.Logger
}

func NewHandlers(authService *service.AuthService, sessionService *service.SessionService, log *logger.Logger) *Handlers {
	return &Handlers{
		authService:    authService,
		sessionService: sessionService,
		log:            log,
	}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	Message     string `json:"message,omitempty"`
}

type verifyResponse struct {
	Valid   bool   `json:"valid"`
	Subject string `json:"subject,omitempty"`
	Error   string `json:"error,omitempty"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	log := h.log.WithRequestID(requestID)

	if r.Method != http.MethodPost {
		log.Warn("method not allowed", zap.String("method", r.Method))
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid request format"})
		return
	}

	if req.Username == "" || req.Password == "" {
		log.Warn("missing credentials", zap.String("username", req.Username))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "username and password are required"})
		return
	}

	if !h.authService.ValidateCredentials(req.Username, req.Password) {
		log.Info("invalid login attempt", zap.String("username", req.Username))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid credentials"})
		return
	}

	// subject
	subject := req.Username
	if req.Username == "admin" {
		subject = "admin"
	}

	// Создание сессии и получение CSRF токена
	sessionID, csrfToken, err := h.sessionService.CreateSession(req.Username, subject)
	if err != nil {
		log.Error("failed to create session", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: "internal server error"})
		return
	}

	// Session cookie (HttpOnly, Secure, SameSite)
	cookies.SetSecureCookie(w, cookies.CookieConfig{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   86400, // 24 часа
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	// CSRF cookie
	cookies.SetSecureCookie(w, cookies.CookieConfig{
		Name:     "csrf_token",
		Value:    csrfToken,
		Path:     "/",
		MaxAge:   86400, // 24 часа
		Secure:   true,
		HttpOnly: false, // JS должен иметь доступ
		SameSite: http.SameSiteLaxMode,
	})

	log.Info("user logged in",
		zap.String("username", req.Username),
		zap.String("session_id", sessionID[:8]+"..."),
	)

	// CSRF токен в теле ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse{
		Message: "Login successful",
	})
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	log := h.log.WithRequestID(requestID)

	// session cookie
	sessionID, err := cookies.GetSessionCookie(r)
	if err == nil && sessionID != "" {
		// Удаление сессии
		h.sessionService.DeleteSession(sessionID)
	}

	// Очистка cookies
	cookies.ClearCookie(w, "session_id", "/")
	cookies.ClearCookie(w, "csrf_token", "/")

	log.Info("user logged out")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}

func (h *Handlers) Verify(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	log := h.log.WithRequestID(requestID)

	// Несколько способов аутентификации:
	// 1. Bearer токен
	// 2. Session cookie

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Bearer token authentication
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			token := parts[1]
			valid, subject := h.authService.ValidateToken(token)
			if valid {
				log.Info("token verified", zap.String("subject", subject))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(verifyResponse{
					Valid:   true,
					Subject: subject,
				})
				return
			}
		}
	}

	// Проверка session cookie
	sessionID, err := cookies.GetSessionCookie(r)
	if err == nil && sessionID != "" {
		session, err := h.sessionService.GetSession(sessionID)
		if err == nil {
			log.Info("session verified", zap.String("subject", session.Subject))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(verifyResponse{
				Valid:   true,
				Subject: session.Subject,
			})
			return
		}
	}

	log.Warn("authentication failed")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(verifyResponse{
		Valid: false,
		Error: "unauthorized",
	})
}

// Эндпоинт для получения CSRF токена
func (h *Handlers) GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	log := h.log.WithRequestID(requestID)

	// Проверка аутентификацию через сессию
	sessionID, err := cookies.GetSessionCookie(r)
	if err != nil {
		log.Warn("no session cookie")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse{Error: "unauthorized"})
		return
	}

	csrfToken, err := h.sessionService.GetCSRFToken(sessionID)
	if err != nil {
		log.Warn("failed to get CSRF token", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse{Error: "unauthorized"})
		return
	}

	log.Info("CSRF token retrieved", zap.String("session_id", sessionID[:8]+"..."))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"csrf_token": csrfToken,
	})
}

// Health check
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "auth",
	})
}
