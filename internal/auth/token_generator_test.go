package auth

import (
	"auth-service/internal/storage"
	"log"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateTokenPair(t *testing.T) {
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
	authorizer := NewAuthorizer(stor)
	userID := "test-user"
	clientIP := "127.0.0.1"

	accessToken, refreshToken, err := authorizer.GenerateTokenPair(userID, clientIP)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if accessToken == "" || refreshToken == "" {
		t.Error("Expected non-empty tokens")
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		t.Fatalf("Expected valid access token, got error %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Expected valid claims")
	}

	if claims["sub"] != userID || claims["ip"] != clientIP {
		t.Error("Expected claims to match input values")
	}
}

func TestValidateAccessToken(t *testing.T) {
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
	authorizer := NewAuthorizer(stor)
	userID := "test-user"
	clientIP := "127.0.0.1"

	accessToken, _, err := authorizer.GenerateTokenPair(userID, clientIP)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	valid, err := authorizer.ValidateAccessToken(accessToken)
	if err != nil || !valid {
		t.Fatalf("Expected token to be valid, got valid=%v, err=%v", valid, err)
	}

	invalidToken := "invalid.token.here"
	valid, err = authorizer.ValidateAccessToken(invalidToken)
	if valid || err == nil {
		t.Fatalf("Expected token to be invalid, got valid=%v, err=%v", valid, err)
	}
}

func TestRefreshAccessToken(t *testing.T) {
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
	authorizer := NewAuthorizer(stor)
	userID := "test-user"
	clientIP := "127.0.0.1"

	accessToken, refreshToken, err := authorizer.GenerateTokenPair(userID, clientIP)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	newAccessToken, newRefreshToken, err := authorizer.RefreshAccessToken(accessToken, refreshToken)
	if err != nil {
		t.Fatalf("Expected no error on refresh, got %v", err)
	}

	if newAccessToken == "" || newRefreshToken == "" {
		t.Error("Expected non-empty new tokens")
	}

	if newAccessToken == accessToken || newRefreshToken == refreshToken {
		t.Errorf("Expected new tokens to be different from old ones \n %s \n %s \n %s \n %s", newAccessToken, accessToken, newRefreshToken, refreshToken)
	}
}
