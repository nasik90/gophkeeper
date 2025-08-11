package pg

import (
	"context"
	"database/sql"
	"time"

	stdlibTransactor "github.com/Thiht/transactor/stdlib"
	"github.com/nasik90/gophkeeper/internal/server/storage"
)

// UpdateSecret - создает запись в таблице Secrets.
// Возвращает ID созданной записи.
func (s *Store) SaveNewSecret(ctx context.Context, key, value, login string, creationDate time.Time) (int, error) {

	userID, err := s.getUserID(ctx, login)
	if err != nil {
		return 0, err
	}

	var id int

	row := s.dbGetter(ctx).QueryRowContext(ctx, `INSERT INTO secrets (key, value, user_id, creation_date, updating_date) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		key, value, userID, creationDate, creationDate)

	if row.Err() != nil {
		return id, row.Err()
	}
	if row.Scan(&id) != nil {
		return id, row.Err()
	}

	return id, nil
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
		WHERE user_id = $1 and key = $2 FOR UPDATE
	`, userID, secretData.Key)
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

	if _, err = tx.ExecContext(ctx, `UPDATE secrets SET value = $1, version_id = $2, updating_date = $3, deletion_mark = $4 WHERE user_id = $5 and key = $6`,
		secretData.Value, versionID, updatingDate, secretData.DeletionMark, userID, secretData.Key); err != nil {
		return 0, err
	}

	updateVersion, err := delete_storeUserSecretsUpdateVersion(ctx, tx, userID, updatingDate)
	if err != nil {
		return updateVersion, err
	}

	return updateVersion, tx.Commit()

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

	queryText := `SELECT key, value, version_id, creation_date, updating_date, deletion_mark FROM secrets WHERE id = $1 and user_id = $2`
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
		if err := rows.Scan(&s.Key, &s.Value, &s.VersionID, &s.CreationDate, &s.UpdatingDate, &s.DeletionMark); err != nil {
			return nil, err
		}
		result = append(result, *s)
	}

	if err := rows.Err(); err != nil {
		return &result, err
	}
	return &result, rows.Close()
}

func (s *Store) GetUserSecretsVersion(ctx context.Context, login string) (int, error) {
	var updateVersion int
	userID, err := s.getUserID(ctx, login)
	if err != nil {
		return updateVersion, err
	}
	queryText := `SELECT update_version FROM users_secrets_update_info WHERE user_id = $1`
	if stdlibTransactor.IsWithinTransaction(ctx) {
		queryText += ` FOR UPDATE`
	}
	row := s.conn.QueryRowContext(ctx, queryText, userID)
	if row.Err() != nil {
		return updateVersion, row.Err()
	}
	err = row.Scan(&updateVersion)
	return updateVersion, err
}

func delete_storeUserSecretsUpdateVersion(ctx context.Context, tx *sql.Tx, userID int, updatingDate time.Time) (int, error) {

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

func (s *Store) UpdateUserSecretsVersion(ctx context.Context, login string, newVersion int, updatingDate time.Time) error {

	userID, err := s.getUserID(ctx, login)
	if err != nil {
		return err
	}

	_, err = s.dbGetter(ctx).ExecContext(ctx,
		`INSERT INTO users_secrets_update_info (user_id, updating_date, update_version)
		VALUES ($1, $2, $3)
	 	ON CONFLICT (user_id)
		DO UPDATE SET updating_date = $2, update_version = $3`,
		userID, updatingDate, newVersion)

	return err

}
