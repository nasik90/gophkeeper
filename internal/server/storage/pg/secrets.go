package pg

import (
	"context"
	"database/sql"
	"time"

	stdlibTransactor "github.com/Thiht/transactor/stdlib"
	types "github.com/nasik90/gophkeeper/internal/common/types"
)

// LoadSecret - создает/обновляет запись в таблице Secrets.
func (s *Store) LoadSecret(ctx context.Context, secretData *types.SecretData, userID int) error {
	qText := `INSERT INTO secrets (guid, key, value, comment, binary_value, user_id, creation_date, updating_date) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (guid)
			DO UPDATE SET key = $2, value = $3, comment = $4, binary_value = $5, user_id = $6, creation_date = $7, updating_date = $8`
	row := s.dbGetter(ctx).QueryRowContext(ctx, qText,
		secretData.Guid,
		secretData.Key,
		secretData.Value,
		secretData.Comment,
		secretData.BinaryValue,
		userID,
		secretData.CreationDate,
		secretData.UpdatingDate)

	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

// GetUserSecretList - возвращает список записей таблицы Secrets.
// Если передать заполненный fromDate, то вернет записи, обновленные с указанной даты.
func (s *Store) GetUserSecretList(ctx context.Context, userID int, fromDate time.Time) (*[]types.SecretData, error) {
	var (
		result []types.SecretData
		rows   *sql.Rows
		err    error
	)

	queryText := `SELECT guid, key, value, version_id, creation_date, updating_date, deletion_mark FROM secrets WHERE user_id = $1`
	if !fromDate.IsZero() {
		queryText = queryText + " and updating_date >= $2"
		rows, err = s.conn.QueryContext(ctx, queryText, userID, fromDate)
	} else {
		rows, err = s.conn.QueryContext(ctx, queryText, userID)
	}

	if err != nil {
		return &result, err
	}
	for rows.Next() {
		s := new(types.SecretData)
		if err := rows.Scan(&s.Guid, &s.Key, &s.Value, &s.VersionID, &s.CreationDate, &s.UpdatingDate, &s.DeletionMark); err != nil {
			return nil, err
		}
		result = append(result, *s)
	}

	if err := rows.Err(); err != nil {
		return &result, err
	}
	return &result, rows.Close()
}

// GetUserSecretsVersion - получает текущую версию данных секретов пользователя.
func (s *Store) GetUserSecretsVersion(ctx context.Context, userID int) (int, error) {
	var updateVersion int
	queryText := `SELECT update_version FROM users_secrets_update_info WHERE user_id = $1`
	if stdlibTransactor.IsWithinTransaction(ctx) {
		queryText += ` FOR UPDATE`
	}
	row := s.conn.QueryRowContext(ctx, queryText, userID)
	if row.Err() != nil {
		return updateVersion, row.Err()
	}
	err := row.Scan(&updateVersion)
	return updateVersion, err
}

// GetSecretVersion - получает текущую версию данных секретА пользователя.
func (s *Store) GetSecretVersion(ctx context.Context, userID int, guid string) (int, error) {
	var versionID int
	queryText := `SELECT version_id FROM secrets WHERE user_id = $1 and guid = $2 `
	if stdlibTransactor.IsWithinTransaction(ctx) {
		queryText += ` FOR UPDATE`
	}
	row := s.conn.QueryRowContext(ctx, queryText, userID, guid)
	if row.Err() != nil {
		return versionID, row.Err()
	}
	err := row.Scan(&versionID)
	return versionID, err
}

// UpdateUserSecretsVersion - обновляет версию данных секретов пользователя.
func (s *Store) UpdateUserSecretsVersion(ctx context.Context, userID int, newVersion int, updatingDate time.Time) error {

	_, err := s.dbGetter(ctx).ExecContext(ctx,
		`INSERT INTO users_secrets_update_info (user_id, updating_date, update_version)
		VALUES ($1, $2, $3)
	 	ON CONFLICT (user_id)
		DO UPDATE SET updating_date = $2, update_version = $3`,
		userID, updatingDate, newVersion)

	return err

}
