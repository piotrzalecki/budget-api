package handlers

import (
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/piotrzalecki/budget-api/internal/config"
	"github.com/piotrzalecki/budget-api/internal/data"
)

var testApp config.AppConfig
var mockDB sqlmock.Sqlmock

func TestMain(m *testing.M) {

	testDB, myMock, _ := sqlmock.New()
	mockDB = myMock

	defer testDB.Close()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	testApp = config.AppConfig{
		Port: 8081,
		Env: "dev",
		Version: "0.0.1",
		InfoLogger: infoLog,
		ErrorLogger: errorLog,
		Models: data.New(testDB),
	}
	

	os.Exit(m.Run())
}
