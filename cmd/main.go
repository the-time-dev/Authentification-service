package main

import (
	"auth-service/internal/auth"
	"auth-service/internal/http_handlers"
	"auth-service/internal/storage"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using system environment variables")
	}

	conn, ok := os.LookupEnv("PG_CONN")
	if !ok {
		log.Fatal("PG_CONN environment variable is not set")
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	stor, err := storage.NewStorage(conn)
	if err != nil {
		log.Fatal("Error creating storage:", err)
	}
	log.Println("Storage created")
	defer stor.Close()
	server, err := http_handlers.NewServer(auth.NewAuthorizer(stor))
	if err != nil {
		log.Fatal("Error creating server:", err)
	}
	log.Println("Server created")

	log.Println("Starting server on port", port)
	http.ListenAndServe(":"+port, server)
}
