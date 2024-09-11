package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/piotrzalecki/budget-api/internal/handlers"
	mid "github.com/piotrzalecki/budget-api/internal/middleware"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Post("/users/login", handlers.Repo.Login)
	mux.Post("/users/logout", handlers.Repo.Logout)

	mux.Post("/validate-token", handlers.Repo.ValidateToken)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(mid.Mid.AuthTokenMiddlewere)

		mux.Post("/users", handlers.Repo.AllUsers)
		mux.Post("/users/save", handlers.Repo.EditUser)
		mux.Post("/users/get/{id}", handlers.Repo.GetUser)
		mux.Post("/users/delete", handlers.Repo.DeleteUser)
		mux.Post("/log-user-out/{id}", handlers.Repo.LogUserOutAdnSetInactive)
	})

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
