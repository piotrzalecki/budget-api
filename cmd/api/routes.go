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

		// users
		mux.Post("/users", handlers.Repo.AllUsers)
		mux.Post("/users/save", handlers.Repo.EditUser)
		mux.Post("/users/get/{id}", handlers.Repo.GetUser)
		mux.Post("/users/delete", handlers.Repo.DeleteUser)
		mux.Post("/log-user-out/{id}", handlers.Repo.LogUserOutAdnSetInactive)

		// tags
		mux.Post("/dashboard/tags", handlers.Repo.Tags)
		mux.Put("/dashboard/tags", handlers.Repo.TagsCreateUpdate)
		mux.Patch("/dashboard/tags", handlers.Repo.TagsCreateUpdate)
		mux.Delete("/dashboard/tags", handlers.Repo.TagsDelete)
		mux.Post("/dashboard/tags/{id}", handlers.Repo.TagById)

		// budgets
		mux.Post("/dashboard/budgets", handlers.Repo.Budgets)
		mux.Put("/dashboard/budgets", handlers.Repo.BudgetsCreateUpdate)
		mux.Patch("/dashboard/budgets", handlers.Repo.BudgetsCreateUpdate)
		mux.Delete("/dashboard/budgets", handlers.Repo.BudgetsDelete)
		mux.Post("/dashboard/budgets/{id}", handlers.Repo.BudgetsById)

		// transactions recurrences
		mux.Post("/dashboard/transactions-recurrences", handlers.Repo.TransactionsRecurrences)
		mux.Put("/dashboard/transactions-recurrences", handlers.Repo.TransactionsRecurrencesCreate)

		//transactions
		mux.Post("/dashboard/transactions", handlers.Repo.TransactionsAll)
		mux.Post("/dashboard/transactions/{id}", handlers.Repo.TransactionsById)
		mux.Delete("/dashboard/transactions", handlers.Repo.TransactionsDelete)
		mux.Put("/dashboard/transactions", handlers.Repo.TransactionsCreateUpdate)
		mux.Patch("/dashboard/transactions", handlers.Repo.TransactionsCreateUpdate)
		mux.Post("/dashboard/transactions/set-all-active", handlers.Repo.TransactionsSetStatusAllActive)
		mux.Post("/dashboard/transactions/set-status", handlers.Repo.TransactionsSetStatus)

		//Logs
		mux.Post("/dashboard/logs", handlers.Repo.Logs)

	})

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
