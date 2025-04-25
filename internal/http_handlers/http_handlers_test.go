package http_handlers

import (
	"auth-service/internal/auth"
	"auth-service/internal/storage"
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAcessHandler(t *testing.T) {
	connString, ok := os.LookupEnv("PG_CONN")
	if !ok {
		t.Fatalf("PG_CONN environment variable is not set")
	}
	stor, err := storage.NewStorage(connString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer func() {
		if err := stor.Close(); err != nil {
			log.Printf("Error closing storage: %v", err)
		}
	}()
	authorizer := auth.NewAuthorizer(stor)
	server, err := NewServer(authorizer)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("POST", "/access/test-user", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.acessHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	if !bytes.Contains(rr.Body.Bytes(), []byte("accessToken")) {
		t.Errorf("Expected response to contain accessToken, got %v", rr.Body.String())
	}
}

func TestRefreshHandler(t *testing.T) {
	connString, ok := os.LookupEnv("PG_CONN")
	if !ok {
		t.Fatalf("PG_CONN environment variable is not set")
	}
	stor, err := storage.NewStorage(connString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer func() {
		if err := stor.Close(); err != nil {
			log.Printf("Error closing storage: %v", err)
		}
	}()
	authorizer := auth.NewAuthorizer(stor)
	server, err := NewServer(authorizer)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	accessToken, refreshToken, err := authorizer.GenerateTokenPair("test-user", "127.0.0.1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("POST", "/refresh", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Refresh", refreshToken)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.refreshHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	if !bytes.Contains(rr.Body.Bytes(), []byte("accessToken")) {
		t.Errorf("Expected response to contain accessToken, got %v", rr.Body.String())
	}
}
