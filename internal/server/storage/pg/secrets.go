package pg

import (
	"context"
	"database/sql"
	"time"

	stdlibTransactor "github.com/Thiht/transactor/stdlib"
	types "github.com/nasik90/gophkeeper/internal/common/types"
)

// UpdateSecret - создает запись в таблице Secrets.
// Возвращает ID созданной записи.
func (s *Store) SaveNewSecret(ctx context.Context, secretData *types.SecretData, userID int, creationDate time.Time) error {

	row := s.dbGetter(ctx).QueryRowContext(ctx, `INSERT INTO secrets (guid, key, value, comment, binary_value, user_id, creation_date, updating_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		secretData.Guid, secretData.Key, secretData.Value, secretData.Comment, secretData.BinaryValue, userID, creationDate, creationDate)

	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

// UpdateSecret - обновляет запись в таблице Secrets.
func (s *Store) UpdateSecret(ctx context.Context, userID int, secretData *types.SecretData, updatingDate time.Time) error {

	_, err := s.dbGetter(ctx).ExecContext(ctx, `UPDATE secrets SET key = $1, value = $2, version_id = $3, comment = $4, binary_value = $5, updating_date = $6, deletion_mark = $7 WHERE user_id = $8 and guid = $9`,
		secretData.Key, secretData.Value, secretData.VersionID, secretData.Comment, secretData.BinaryValue, updatingDate, secretData.DeletionMark, userID, secretData.Guid)

	return err

}

// GetUserSecretList - возвращает список записей таблицы Secrets.
// Если передать заполненный fromDate, то вернет записи, обновленные с указанной даты
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

func (s *Store) UpdateUserSecretsVersion(ctx context.Context, userID int, newVersion int, updatingDate time.Time) error {

	_, err := s.dbGetter(ctx).ExecContext(ctx,
		`INSERT INTO users_secrets_update_info (user_id, updating_date, update_version)
		VALUES ($1, $2, $3)
	 	ON CONFLICT (user_id)
		DO UPDATE SET updating_date = $2, update_version = $3`,
		userID, updatingDate, newVersion)

	return err

}
