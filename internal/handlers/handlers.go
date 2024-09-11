package handlers

import "github.com/piotrzalecki/budget-api/internal/config"

type Repository struct {
	App *config.AppConfig
	// DB  repository.DatabaseRepo
}

var Repo *Repository

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		// DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}
