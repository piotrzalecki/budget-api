package config

import (
	"log"

	"github.com/piotrzalecki/budget-api/internal/data"
)

type AppConfig struct {
	Port        int
	Env         string
	Version     string
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	Models      data.Models
}
