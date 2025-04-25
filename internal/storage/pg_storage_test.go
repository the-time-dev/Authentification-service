package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewStorage(t *testing.T) {
	connString, ok := os.LookupEnv("PG_CONN")
	if !ok {
		t.Fatalf("PG_CONN environment variable is not set")
	}
	stor, err := NewStorage(connString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer func() {
		if err := stor.Close(); err != nil {
			t.Fatalf("Error closing storage: %v", err)
		}
	}()
}

func TestAddToken(t *testing.T) {
	connString, ok := os.LookupEnv("PG_CONN")
	if !ok {
		t.Fatalf("PG_CONN environment variable is not set")
	}
	stor, err := NewStorage(connString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer func() {
		if err := stor.Close(); err != nil {
			t.Fatalf("Error closing storage: %v", err)
		}
	}()

	tokenUUID := func() string {
		uuid, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("Expected no error on uuid generation, got %v", err)
		}
		return uuid.String()
	}()
	tokenHash := time.Now().String()
	time.Sleep(time.Second / 10)

	err = stor.AddToken(tokenUUID, tokenHash)
	if err != nil {
		t.Fatalf("Expected no error on add token, got %v", err)
	}

	var hash string
	err = stor.conn.QueryRow(context.Background(), "SELECT token_hash FROM refresh_tokens WHERE jti = $1", tokenUUID).Scan(&hash)
	if err != nil {
		t.Fatalf("Expected no error on token retrieval, got %v", err)
	}

	if hash != tokenHash {
		t.Errorf("Expected hash to be %s, got %s", tokenHash, hash)
	}
}

func TestGetHash(t *testing.T) {
	connString, ok := os.LookupEnv("PG_CONN")
	if !ok {
		t.Fatalf("PG_CONN environment variable is not set")
	}
	stor, err := NewStorage(connString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer func() {
		if err := stor.Close(); err != nil {
			t.Fatalf("Error closing storage: %v", err)
		}
	}()

	tokenUUID := func() string {
		uuid, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("Expected no error on uuid generation, got %v", err)
		}
		return uuid.String()
	}()
	tokenHash := time.Now().String()
	time.Sleep(time.Second / 10)

	err = stor.AddToken(tokenUUID, tokenHash)
	if err != nil {
		t.Fatalf("Expected no error on add token, got %v", err)
	}

	hash, err := stor.GetHash(tokenUUID)
	if err != nil {
		t.Fatalf("Expected no error on get hash, got %v", err)
	}

	if hash != tokenHash {
		t.Errorf("Expected hash to be %s, got %s", tokenHash, hash)
	}
}

func TestUpdateToken(t *testing.T) {
	connString, ok := os.LookupEnv("PG_CONN")
	if !ok {
		t.Fatalf("PG_CONN environment variable is not set")
	}
	stor, err := NewStorage(connString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer func() {
		if err := stor.Close(); err != nil {
			t.Fatalf("Error closing storage: %v", err)
		}
	}()

	tokenUUID := func() string {
		uuid, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("Expected no error on uuid generation, got %v", err)
		}
		return uuid.String()
	}()
	tokenHash := time.Now().String()
	time.Sleep(time.Second / 10)
	newTokenHash := time.Now().String()
	time.Sleep(time.Second / 10)

	err = stor.AddToken(tokenUUID, tokenHash)
	if err != nil {
		t.Fatalf("Expected no error on add token, got %v", err)
	}

	err = stor.UpdateToken(tokenUUID, newTokenHash)
	if err != nil {
		t.Fatalf("Expected no error on update token, got %v", err)
	}

	hash, err := stor.GetHash(tokenUUID)
	if err != nil {
		t.Fatalf("Expected no error on get hash, got %v", err)
	}

	if hash != newTokenHash {
		t.Errorf("Expected hash to be %s, got %s", newTokenHash, hash)
	}
}
