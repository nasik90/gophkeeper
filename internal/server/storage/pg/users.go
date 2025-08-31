package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nasik90/gophkeeper/internal/server/storage"
)

func (s *Store) SaveNewUser(ctx context.Context, login, password string) error {
	_, err := s.dbGetter(ctx).ExecContext(ctx, `INSERT INTO users (login, password) VALUES ($1, $2)`, login, password)
	err = saveNewUserCheckInsertError(err)
	return err
}

func saveNewUserCheckInsertError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		if pgErr.ConstraintName == "login_ukey" {
			return storage.ErrUserNotUnique
		}
	}
	return err
}

func (s *Store) UserIsValid(ctx context.Context, login, password string) (bool, error) {
	rows, err := s.dbGetter(ctx).ExecContext(ctx, `SELECT FROM users WHERE login = $1 and password = $2`, login, password)
	if err != nil {
		return false, err
	}
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected > 0 {
		return true, nil
	}
	return false, nil
}

func (s *Store) getUserID(ctx context.Context, login string) (int, error) {
	row := s.dbGetter(ctx).QueryRowContext(ctx, `SELECT id FROM users WHERE login = $1`, login)
	var userID int
	if err := row.Scan(&userID); err != nil {
		return 0, err
	}
	return userID, nil
}
