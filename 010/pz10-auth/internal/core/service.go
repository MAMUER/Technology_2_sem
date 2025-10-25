package core

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

type UserRepo interface {
	CheckPassword(email, pass string) (*User, error)
	GetUserByID(id int64) (*User, error)
}

type JWTService interface {
	SignAccess(userID int64, email, role string) (string, error)
	SignRefresh(userID int64) (string, error)
	Parse(tokenStr string) (jwt.MapClaims, error)
}

type RefreshStore interface {
	Store(token string, exp time.Time) error
	IsRevoked(token string) bool
	Revoke(token string) error
	Cleanup()
}

type Service struct {
	repo         UserRepo
	jwt          JWTService
	refreshStore RefreshStore
}

func NewService(r UserRepo, j JWTService, rs RefreshStore) *Service {
	service := &Service{repo: r, jwt: j, refreshStore: rs}
	// Запускаем очистку устаревших токенов
	go service.cleanupRoutine()
	return service
}

func (s *Service) cleanupRoutine() {
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		s.refreshStore.Cleanup()
	}
}

func (s *Service) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpError(w, 400, "invalid_request", "Invalid JSON")
		return
	}

	user, err := s.repo.CheckPassword(in.Email, in.Password)
	if err != nil {
		httpError(w, 401, "invalid_credentials", "Wrong email or password")
		return
	}

	accessToken, err := s.jwt.SignAccess(user.ID, user.Email, user.Role)
	if err != nil {
		httpError(w, 500, "token_creation_error", "Cannot create access token")
		return
	}

	refreshToken, err := s.jwt.SignRefresh(user.ID)
	if err != nil {
		httpError(w, 500, "token_creation_error", "Cannot create refresh token")
		return
	}

	// Сохраняем refresh токен
	refreshExp := time.Now().Add(7 * 24 * time.Hour)
	if err := s.refreshStore.Store(refreshToken, refreshExp); err != nil {
		httpError(w, 500, "storage_error", "Cannot store refresh token")
		return
	}

	jsonOK(w, map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    900, // 15 минут в секундах
		"user": map[string]any{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (s *Service) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var in struct {
		RefreshToken string `json:"refresh_token"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpError(w, 400, "invalid_request", "Invalid JSON")
		return
	}

	// Проверяем, не отозван ли токен
	if s.refreshStore.IsRevoked(in.RefreshToken) {
		httpError(w, 401, "token_revoked", "Refresh token was revoked")
		return
	}

	// Парсим refresh токен
	claims, err := s.jwt.Parse(in.RefreshToken)
	if err != nil {
		httpError(w, 401, "invalid_token", "Invalid refresh token")
		return
	}

	// Проверяем, что это refresh токен
	if claims["type"] != "refresh" {
		httpError(w, 401, "invalid_token_type", "Not a refresh token")
		return
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		httpError(w, 401, "invalid_token", "Invalid user ID in token")
		return
	}

	// Получаем пользователя
	user, err := s.repo.GetUserByID(int64(userID))
	if err != nil {
		httpError(w, 401, "user_not_found", "User not found")
		return
	}

	// Отзываем старый refresh токен
	s.refreshStore.Revoke(in.RefreshToken)

	// Генерируем новую пару токенов
	accessToken, err := s.jwt.SignAccess(user.ID, user.Email, user.Role)
	if err != nil {
		httpError(w, 500, "token_creation_error", "Cannot create access token")
		return
	}

	newRefreshToken, err := s.jwt.SignRefresh(user.ID)
	if err != nil {
		httpError(w, 500, "token_creation_error", "Cannot create refresh token")
		return
	}

	// Сохраняем новый refresh токен
	refreshExp := time.Now().Add(7 * 24 * time.Hour)
	if err := s.refreshStore.Store(newRefreshToken, refreshExp); err != nil {
		httpError(w, 500, "storage_error", "Cannot store refresh token")
		return
	}

	jsonOK(w, map[string]any{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
		"token_type":    "Bearer",
		"expires_in":    900,
	})
}

func (s *Service) MeHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(map[string]interface{})
	if !ok {
		httpError(w, 401, "unauthorized", "No claims in context")
		return
	}

	userID, _ := claims["sub"].(float64)
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)

	jsonOK(w, map[string]any{
		"id":    int64(userID),
		"email": email,
		"role":  role,
	})
}

// ABAC: пользователь может получать только свой профиль
func (s *Service) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(map[string]interface{})
	if !ok {
		httpError(w, 401, "unauthorized", "No claims in context")
		return
	}

	// Извлекаем ID пользователя из URL
	userIDStr := chi.URLParam(r, "id")
	requestedID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		httpError(w, 400, "invalid_user_id", "Invalid user ID format")
		return
	}

	// Извлекаем ID из токена
	tokenUserID, _ := claims["sub"].(float64)
	userRole, _ := claims["role"].(string)

	// ABAC проверка: обычный пользователь может получать только свой профиль
	if userRole == "user" && int64(tokenUserID) != requestedID {
		httpError(w, 403, "forbidden", "You can only access your own profile")
		return
	}

	// Получаем пользователя
	user, err := s.repo.GetUserByID(requestedID)
	if err != nil {
		httpError(w, 404, "user_not_found", "User not found")
		return
	}

	jsonOK(w, map[string]any{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	})
}

func (s *Service) AdminStats(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(map[string]interface{})
	if !ok {
		httpError(w, 401, "unauthorized", "No claims in context")
		return
	}

	// Дополнительная проверка роли
	role, _ := claims["role"].(string)
	if role != "admin" {
		httpError(w, 403, "forbidden", "Admin access required")
		return
	}

	jsonOK(w, map[string]any{
		"stats": "admin only data",
		"users": 42,
		"time":  time.Now().Format(time.RFC3339),
	})
}

func (s *Service) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	var in struct {
		RefreshToken string `json:"refresh_token"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpError(w, 400, "invalid_request", "Invalid JSON")
		return
	}

	// Отзываем refresh токен
	s.refreshStore.Revoke(in.RefreshToken)

	jsonOK(w, map[string]any{
		"message": "Successfully logged out",
	})
}

func jsonOK(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func httpError(w http.ResponseWriter, code int, errorType, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   errorType,
		"details": details,
	})
}