// Package main Notes API server.
//
// @title           Notes API
// @version         1.0
// @description     Учебный REST API для заметок (CRUD).
// @contact.name    Backend Course
// @contact.email   example@university.ru
// @BasePath        /api/v1
package main

import httpSwagger "github.com/swaggo/http-swagger"

import (
	"example.com/notes-api/internal/http"
	"example.com/notes-api/internal/http/handlers"
	"example.com/notes-api/internal/repo"
	"log"
	"net/http"
)

func main() {
	repo := repo.NewNoteRepoMem()
	h := &handlers.Handler{Repo: repo}
	r := httpx.NewRouter(h)

	r.Get("/docs/*", httpSwagger.WrapHandler) // для chi: r.Get

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
