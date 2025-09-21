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
func (s *Store) SaveNewSecret(ctx context.Context, secretData *types.SecretData, userID int, creationDate time.Time) (int, error) {
	var id int

	row := s.dbGetter(ctx).QueryRowContext(ctx, `INSERT INTO secrets (key, value, comment, binary_value, user_id, creation_date, updating_date) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		secretData.Key, secretData.Value, secretData.Comment, secretData.BinaryValue, userID, creationDate, creationDate)

	if row.Err() != nil {
		return id, row.Err()
	}
	if row.Scan(&id) != nil {
		return id, row.Err()
	}

	return id, nil
}

// UpdateSecret - обновляет запись в таблице Secrets.
func (s *Store) UpdateSecret(ctx context.Context, userID int, secretData *types.SecretData, updatingDate time.Time) error {

	_, err := s.dbGetter(ctx).ExecContext(ctx, `UPDATE secrets SET key = $1, value = $2, version_id = $3, comment = $4, binary_value = $5, updating_date = $6, deletion_mark = $7 WHERE user_id = $8 and id = $9`,
		secretData.Key, secretData.Value, secretData.VersionID, secretData.Comment, secretData.BinaryValue, updatingDate, secretData.DeletionMark, userID, secretData.Id)

	return err

}

// GetUserSecretList - возвращает список записей таблицы Secrets.
// Если передать заполненный fromDate, то вернет записи, обновленные с указанной даты
func (s *Store) GetUserSecretList(ctx context.Context, userID int, secretID int, fromDate time.Time) (*[]types.SecretData, error) {
	var (
		result []types.SecretData
		rows   *sql.Rows
		err    error
	)

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
		s := new(types.SecretData)
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

func (s *Store) GetSecretVersion(ctx context.Context, userID int, secretID int) (int, error) {
	var versionID int
	queryText := `SELECT version_id FROM secrets WHERE user_id = $1 and id = $2 `
	if stdlibTransactor.IsWithinTransaction(ctx) {
		queryText += ` FOR UPDATE`
	}
	row := s.conn.QueryRowContext(ctx, queryText, userID, secretID)
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
