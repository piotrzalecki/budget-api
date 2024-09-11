package middleware

import (
	"net/http"

	"github.com/piotrzalecki/budget-api/internal/config"
	"github.com/piotrzalecki/budget-api/internal/handlers"
)

type Middleware struct {
	App *config.AppConfig
	// DB  repository.DatabaseRepo
}

var Mid *Middleware

func NewMid(a *config.AppConfig) *Middleware {
	return &Middleware{
		App: a,
		// DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewMiddleware(m *Middleware) {
	Mid = m
}

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (mid *Middleware) AuthTokenMiddlewere(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := mid.App.Models.Token.AuthenticateToken(r)
		if err != nil {
			mid.App.ErrorLogger.Println(err)
			payload := jsonResponse{
				Error:   true,
				Message: "invalid authentication token",
			}

			_ = handlers.Repo.WriteJSON(w, http.StatusUnauthorized, payload)
			return
		}
		next.ServeHTTP(w, r)
	})
}
