package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"example.com/pz10-auth/internal/core"
	"example.com/pz10-auth/internal/http/middleware"
	"example.com/pz10-auth/internal/platform/config"
	"example.com/pz10-auth/internal/platform/jwt"
	"example.com/pz10-auth/internal/repo"
)

func Build(cfg config.Config) http.Handler {
	r := chi.NewRouter()

	// Инициализация зависимостей
	userRepo := repo.NewUserMem()
	refreshStore := repo.NewRefreshStore()
	jwtService := jwt.NewHS256(cfg.JWTSecret, cfg.AccessTTL, cfg.RefreshTTL)
	service := core.NewService(userRepo, jwtService, refreshStore)

	// Публичные маршруты
	r.Post("/api/v1/login", service.LoginHandler)
	r.Post("/api/v1/refresh", service.RefreshHandler)
	r.Post("/api/v1/logout", service.LogoutHandler)

	// Защищённые маршруты
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.AuthN(jwtService))
		r.Use(middleware.AuthZRoles("user", "admin"))

		r.Get("/me", service.MeHandler)
		r.Get("/users/{id}", service.GetUserHandler) // ABAC защищенный эндпоинт

		// Админские маршруты
		r.Route("/admin", func(r chi.Router) {
			r.Use(middleware.AuthZRoles("admin"))
			r.Get("/stats", service.AdminStats)
		})
	})

	return r
}
