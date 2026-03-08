package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Thiht/transactor"
	stdlibTransactor "github.com/Thiht/transactor/stdlib"
	"github.com/pressly/goose"
)

type Store struct {
	conn     *sql.DB
	dbGetter stdlibTransactor.DBGetter
}

func NewStore(conn *sql.DB) (*Store, transactor.Transactor, error) {
	transactor, dbGetter := stdlibTransactor.NewTransactor(
		conn,
		stdlibTransactor.NestedTransactionsSavepoints,
	)
	s := &Store{conn: conn, dbGetter: dbGetter}
	fmt.Println(os.Getwd())
	//dir := "internal/migrations/server/pg"
	cwd, _ := os.Getwd()
	dir := filepath.Join(cwd, "..", "..", "internal", "migrations", "server", "pg")
	// Применение миграций
	migrationErr := goose.Up(conn, dir)
	//Откат миграций
	if migrationErr != nil {
		err := goose.Down(conn, dir)
		if err == nil {
			return s, transactor, migrationErr
		} else {
			return s, transactor, errors.Join(migrationErr, err)
		}
	}
	return s, transactor, nil
}

func (s Store) Close() error {
	return s.conn.Close()
}
