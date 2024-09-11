package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/piotrzalecki/budget-api/internal/config"
	"github.com/piotrzalecki/budget-api/internal/data"
	"github.com/piotrzalecki/budget-api/internal/driver"
	"github.com/piotrzalecki/budget-api/internal/handlers"
	mid "github.com/piotrzalecki/budget-api/internal/middleware"
)

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dsn := os.Getenv("DSN")
	environment := os.Getenv("ENV")

	// "host=localhost port=5432 user=postgres password=password dbname=vueapi sslmode=disable timezone=UTC connect_timeout=5"
	db, err := driver.ConnectPostgres(dsn)
	if err != nil {
		log.Fatal("Can not connect to database")
	}
	defer db.SQL.Close()

	appConfig := &config.AppConfig{
		Port:        8081,
		Env:         environment,
		Version:     "0.0.1",
		InfoLogger:  infoLog,
		ErrorLogger: errorLog,
		Models:      data.New(db.SQL),
	}

	repo := handlers.NewRepo(appConfig)
	handlers.NewHandlers(repo)

	midRepo := mid.NewMid(appConfig)
	mid.NewMiddleware(midRepo)

	err = serve(appConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func serve(cfg *config.AppConfig) error {
	cfg.InfoLogger.Println("API listening on port", cfg.Port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: routes(),
	}

	return srv.ListenAndServe()
}
