package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nasik90/gophkeeper/internal/common/types"
	"github.com/pressly/goose"
)

// Store - структура для хранения подключения к БД.
type Store struct {
	conn *sql.DB
}

// NewStore создает экземпляр структуры Store и применяет миграцию.
func NewStore(conn *sql.DB) (*Store, error) {
	s := &Store{conn: conn}
	dir := "internal/migrations/client/sqlite"
	// Применение миграций
	migrationErr := goose.Up(conn, dir)
	//Откат миграций
	if migrationErr != nil {
		err := goose.Down(conn, dir)
		if err == nil {
			return s, migrationErr
		} else {
			return s, errors.Join(migrationErr, err)
		}
	}
	return s, nil
}

// SaveNewSecret сохраняет в БД новую запись.
func (s *Store) SaveNewSecret(ctx context.Context, secretData types.SecretData) error {
	return nil
}

// UpdateSecret обновляет в БД запись.
func (s *Store) UpdateSecret(ctx context.Context, secretData types.SecretData) error {
	return nil
}

// GetSecret получает запись из БД.
func (s *Store) GetSecret(ctx context.Context, id int) error {
	return nil
}

// GetSecrets получает запись из БД.
func (s *Store) GetSecrets(ctx context.Context) error {
	return nil
}

// GetSecretsToken получает из БД токен для шифровки и дешифровки секретов.
func (s *Store) GetSecretsToken(ctx context.Context) (string, error) {
	return "", nil
}

// SaveSecretsToken записывает в БД токен для шифровки и дешифровки секретов.
func (s *Store) SaveSecretsToken(ctx context.Context, token string) error {
	return nil
}
