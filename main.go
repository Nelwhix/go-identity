package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	gHandlers "github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"go-identity/handlers"
	"go-identity/pkg"
	"go-identity/pkg/middlewares"
	"go-identity/pkg/models"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	ServerPort = ":8080"
)

var validate *validator.Validate

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	fileName := filepath.Join("logs", "app_logs.txt")
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := pkg.CreateNewLogger(f)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	validate = validator.New(validator.WithRequiredStructEnabled())

	conn, err := pkg.CreateDbConn()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	model := models.Model{
		Conn: conn,
	}

	handler := handlers.Handler{
		Model:     model,
		Logger:    logger,
		Validator: validate,
	}

	r := http.NewServeMux()
	r.HandleFunc("POST /auth/signup", handler.SignUp)

	fmt.Printf("Go-identity started at http://localhost:%s\n", ServerPort)

	err = http.ListenAndServe(ServerPort, gHandlers.CombinedLoggingHandler(os.Stdout, middlewares.ContentTypeMiddleware(r)))

	if err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}
