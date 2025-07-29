package pg

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/nasik90/gophkeeper/internal/app/storage"

	"github.com/pressly/goose"
)

type Store struct {
	conn *sql.DB
}

func NewStore(conn *sql.DB) (*Store, error) {
	s := &Store{conn: conn}
	dir := "internal/migrations/pg"
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

func (s Store) Close() error {
	return s.conn.Close()
}

// UpdateSecret - создает запись в таблице Secrets.
// Возвращает версию данных пользователя.
func (s *Store) SaveNewSecret(ctx context.Context, key, value, login string) (int, error) {

	userID, err := s.getUserID(ctx, login)
	if err != nil {
		return 0, err
	}

	creationDate := time.Now()

	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, `INSERT INTO secrets (key, value, user_id, creation_date, updating_date) VALUES ($1, $2, $3, $4, $5)`,
		key, value, userID, creationDate, creationDate); err != nil {
		return 0, err
	}

	updateVersion, err := storeUserSecretsUpdateVersion(ctx, tx, userID, creationDate)
	if err != nil {
		return updateVersion, err
	}

	return updateVersion, tx.Commit()
}

// UpdateSecret - обновляет запись в таблице Secrets.
// Возвращает версию данных пользователя.
func (s *Store) UpdateSecret(ctx context.Context, login string, secretData *storage.SecretData) (int, error) {

	userID, err := s.getUserID(ctx, login)
	if err != nil {
		return 0, err
	}

	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var versionID int

	row := tx.QueryRowContext(ctx, `
		SELECT version_id 
		FROM secrets
		WHERE id = $1 FOR UPDATE
	`, secretData.Id)
	if row.Err() != nil {
		return 0, row.Err()
	}
	if row.Scan(&versionID) != nil {
		return 0, row.Err()
	}
	// Если переданная с клиента версия отличается от той, что в БД,
	// то это значит кто-то модифицировал запись, надо прекратить дальнейшую обработу и вернуть ошибку
	if versionID != secretData.VersionID {
		return 0, storage.ErrVersionIdNotTrue
	}

	versionID = versionID + 1
	updatingDate := time.Now()

	if _, err = tx.ExecContext(ctx, `UPDATE secrets SET key = $1, value = $2, user_id = $3 , version_id = $4, updating_date = $5, deletion_mark = $6`,
		secretData.Key, secretData.Value, userID, versionID, updatingDate, secretData.DeletionMark); err != nil {
		return 0, err
	}

	updateVersion, err := storeUserSecretsUpdateVersion(ctx, tx, userID, updatingDate)
	if err != nil {
		return updateVersion, err
	}

	return updateVersion, tx.Commit()

}

func (s *Store) GetSecretData(ctx context.Context, login string, secretID int) (string, string, error) {
	var key, value string
	userID, err := s.getUserID(ctx, login)
	if err != nil {
		return key, value, err
	}
	row := s.conn.QueryRowContext(ctx, `SELECT key, value FROM secrets WHERE id = $1 and user_id = $2`, secretID, userID)
	if row.Err() != nil {
		return key, value, row.Err()
	}
	err = row.Scan(&key, &value)
	return key, value, err
}

// GetSecretData - возвращает список записей таблицы Secrets.
// Если передать заполненный fromDate, то вернет записи, обновленные с указанной даты
func (s *Store) GetUserSecretList(ctx context.Context, login string, secretID int, fromDate time.Time) (*[]storage.SecretData, error) {
	var (
		result []storage.SecretData
		rows   *sql.Rows
	)
	userID, err := s.getUserID(ctx, login)
	if err != nil {
		return &result, err
	}

	queryText := `SELECT id, key, value, version_id, creation_date, updating_date, deletion_mark FROM secrets WHERE id = $1 and user_id = $2`
	if !fromDate.IsZero() {
		queryText = queryText + "updating_date >= $3"
		rows, err = s.conn.QueryContext(ctx, queryText, secretID, userID, fromDate)
	} else {
		rows, err = s.conn.QueryContext(ctx, queryText, secretID, userID)
	}

	if err != nil {
		return &result, err
	}
	for rows.Next() {
		s := new(storage.SecretData)
		if err := rows.Scan(&s.Id, &s.Key, &s.Value, &s.VersionID, &s.CreationDate, &s.UpdatingDate, &s.DeletionMark); err != nil {
			return nil, err
		}
		result = append(result, *s)
	}

	if err := rows.Err(); err != nil {
		return &result, err
	}
	return &result, rows.Close()
}

func (s *Store) GetUserSecretsUpdateVersion(ctx context.Context, login string) (int, error) {
	var updateVersion int
	userID, err := s.getUserID(ctx, login)
	if err != nil {
		return updateVersion, err
	}

	row := s.conn.QueryRowContext(ctx, `SELECT update_version FROM users_secrets_update_info WHERE user_id = $1`, userID)
	if row.Err() != nil {
		return updateVersion, row.Err()
	}
	err = row.Scan(&updateVersion)
	return updateVersion, err
}

func storeUserSecretsUpdateVersion(ctx context.Context, tx *sql.Tx, userID int, updatingDate time.Time) (int, error) {

	var updateVersion int

	row := tx.QueryRowContext(ctx, `
		SELECT version_id 
		FROM users_secrets_update_info
		WHERE user_id = $1 FOR UPDATE
	`, userID)
	if row.Err() != nil {
		return updateVersion, row.Err()
	}
	if row.Scan(&updateVersion) != nil {
		return updateVersion, row.Err()
	}

	updateVersion++

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO users_secrets_update_info (user_id, updating_date, update_version)
		VALUES ($1, $2, $3)
	 	ON CONFLICT (user_id)
		DO UPDATE SET updating_date = $2, update_version = $3`,
		userID, updatingDate, updateVersion); err != nil {
		return updateVersion, err
	}

	return updateVersion, nil
}
