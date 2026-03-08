package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/nasik90/gophkeeper/internal/common/types"
)

// Store - структура для хранения подключения к БД.
type Store struct {
	conn *sql.DB
}

const dateFormat = "2006-01-02 15:04:05"

// NewStore создает экземпляр структуры Store и применяет миграцию.
func NewStore(conn *sql.DB) (*Store, error) {
	s := &Store{conn: conn}
	//dir := "internal/migrations/client/sqlite"
	// cwd, _ := os.Getwd()
	// dir := filepath.Join(cwd, "..", "..", "internal", "migrations", "client", "sqlite")
	// queryText := `CREATE TABLE secrets_update_info (
	// 	updating_date TEXT NOT NULL,
	// 	update_version INTEGER DEFAULT 0
	// );`
	// conn.Exec(queryText)

	// // Применение миграций
	// migrationErr := goose.Up(conn, dir)
	// //Откат миграций
	// if migrationErr != nil {
	// 	err := goose.Down(conn, dir)
	// 	if err == nil {
	// 		return s, migrationErr
	// 	} else {
	// 		return s, errors.Join(migrationErr, err)
	// 	}
	// }
	return s, nil
}

func (s *Store) Close() error {
	return s.conn.Close()
}

// SaveNewSecret сохраняет в БД новую запись.
func (s *Store) SaveNewSecret(ctx context.Context, secretData *types.SecretData) error {
	qText := `INSERT INTO secrets (guid, key, value, binary_value, version_id, creation_date, updating_date, deletion_mark, to_send, comment) VALUES(?,?,?,?,?,?,?,?,?,?)`
	_, err := s.conn.ExecContext(ctx, qText,
		secretData.Guid,
		secretData.Key,
		secretData.Value,
		secretData.BinaryValue,
		secretData.VersionID,
		secretData.CreationDate.Format(dateFormat),
		secretData.UpdatingDate.Format(dateFormat),
		secretData.DeletionMark,
		secretData.ToSend,
		secretData.Comment)
	return err
}

// UpdateSecret обновляет в БД запись.
func (s *Store) UpdateSecret(ctx context.Context, secretData *types.SecretData) error {
	qText := `UPDATE secrets SET key = @key, value = @value, binary_value = @binary_value, version_id = @version_id, creation_date = @creation_date, updating_date = @updating_date, deletion_mark = @deletion_mark, to_send = @to_send WHERE guid = @guid`
	_, err := s.conn.ExecContext(ctx, qText,
		sql.Named("key", secretData.Key),
		sql.Named("value", secretData.Value),
		sql.Named("binary_value", secretData.BinaryValue),
		sql.Named("version_id", secretData.VersionID),
		sql.Named("creation_date", secretData.CreationDate.Format(dateFormat)),
		sql.Named("updating_date", secretData.UpdatingDate.Format(dateFormat)),
		sql.Named("deletion_mark", secretData.DeletionMark),
		sql.Named("to_send", secretData.ToSend),
		sql.Named("guid", secretData.Guid))
	return err
}

// InsertUpdateSecret вставляет запись, в случае конфликта обновляет.
func (s *Store) InsertUpdateSecret(ctx context.Context, secretData *types.SecretData) error {
	qText := `INSERT OR REPLACE INTO secrets (guid, key, value, binary_value, version_id, creation_date, updating_date, deletion_mark, to_send, comment) VALUES(?,?,?,?,?,?,?,?,?,?)`
	_, err := s.conn.ExecContext(ctx, qText,
		secretData.Guid,
		secretData.Key,
		secretData.Value,
		secretData.BinaryValue,
		secretData.VersionID,
		secretData.CreationDate.Format(dateFormat),
		secretData.UpdatingDate.Format(dateFormat),
		secretData.DeletionMark,
		secretData.ToSend,
		secretData.Comment)
	return err
}

// GetSecret получает запись из БД.
func (s *Store) GetSecret(ctx context.Context, id int) error {
	return nil
}

// GetSecrets получает запись из БД.
func (s *Store) GetSecrets(ctx context.Context, toSend bool) (*[]types.SecretData, error) {
	var (
		result []types.SecretData
		rows   *sql.Rows
		err    error
	)
	qText := `SELECT guid, key, value, binary_value, version_id, creation_date, updating_date, deletion_mark, to_send, comment FROM secrets`
	if toSend {
		qText = qText + ` WHERE to_send = 1`
		rows, err = s.conn.QueryContext(ctx, qText)
	} else {
		rows, err = s.conn.QueryContext(ctx, qText)
	}
	if err != nil {
		return &result, err
	}
	var creationDate, updatingDate string
	for rows.Next() {
		s := new(types.SecretData)
		if err = rows.Scan(&s.Guid, &s.Key, &s.Value, &s.BinaryValue, &s.VersionID, &creationDate, &updatingDate, &s.DeletionMark, &s.ToSend, &s.Comment); err != nil {
			return nil, err
		}
		if s.CreationDate, err = time.Parse(dateFormat, creationDate); err != nil {
			return nil, err
		}
		if s.UpdatingDate, err = time.Parse(dateFormat, updatingDate); err != nil {
			return nil, err
		}
		result = append(result, *s)
	}

	if err = rows.Err(); err != nil {
		return &result, err
	}
	return &result, rows.Close()
}

// SaveDataVersion обновляет в БД версию данных юзера.
func (s *Store) SaveDataVersion(ctx context.Context, dataVersion time.Time) error {
	var recordCount int
	qText := `SELECT COUNT(*) FROM secrets_update_info`
	row := s.conn.QueryRowContext(ctx, qText)
	if row.Err() != nil {
		return row.Err()
	}
	if err := row.Scan(&recordCount); err != nil {
		return err
	}
	if recordCount == 0 {
		qText = `INSERT INTO secrets_update_info (data_version) VALUES (?)`
	} else {
		qText = `UPDATE secrets_update_info SET data_version = ?`
	}
	_, err := s.conn.ExecContext(ctx, qText, dataVersion.Format((dateFormat)))
	return err
}

// GetDataVersion получает из БД версию данных юзера.
func (s *Store) GetDataVersion(ctx context.Context) (time.Time, error) {
	var dataVersion time.Time
	qText := `SELECT data_version FROM secrets_update_info`
	row := s.conn.QueryRowContext(ctx, qText)
	if row.Err() != nil {
		return dataVersion, row.Err()
	}
	if err := row.Scan(&dataVersion); err != nil {
		return dataVersion, row.Err()
	}
	return dataVersion, nil
}
