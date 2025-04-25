package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	conn *pgx.Conn
}

func NewStorage(connString string) (*Storage, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	var exists bool
	err = conn.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'refresh_tokens')",
	).Scan(&exists)
	if err != nil {
		conn.Close(context.Background())
		return nil, err
	}

	if !exists {
		_, err = conn.Exec(context.Background(), `
			CREATE TABLE refresh_tokens (
				id SERIAL PRIMARY KEY,
				token_hash TEXT NOT NULL UNIQUE,
				jti UUID NOT NULL UNIQUE,
				created_at TIMESTAMP NOT NULL DEFAULT NOW()
			);
		`,
		)
		if err != nil {
			conn.Close(context.Background())
			return nil, err
		}

		_, err = conn.Exec(context.Background(),
			"CREATE INDEX idx_refresh_tokens_jti ON refresh_tokens(jti);",
		)
		if err != nil {
			conn.Close(context.Background())
			return nil, err
		}
	}

	return &Storage{conn: conn}, nil
}

func (s *Storage) AddToken(tokenUUID, tokenHash string) error {
	query := "INSERT INTO refresh_tokens (jti, token_hash) VALUES ($1, $2)"
	_, err := s.conn.Exec(context.Background(), query, tokenUUID, tokenHash)
	return err
}

func (s *Storage) GetHash(tokenUUID string) (string, error) {
	var hash string
	query := "SELECT token_hash FROM refresh_tokens WHERE jti = $1"
	err := s.conn.QueryRow(context.Background(), query, tokenUUID).Scan(&hash)
	return hash, err
}

func (s *Storage) UpdateToken(tokenUUID, tokenHash string) error {
	query := "UPDATE refresh_tokens SET token_hash = $2 WHERE jti = $1"
	_, err := s.conn.Exec(context.Background(), query, tokenUUID, tokenHash)
	return err
}

func (s *Storage) Close() error {
	return s.conn.Close(context.Background())
}
