package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/Thiht/transactor"
	stdlibTransactor "github.com/Thiht/transactor/stdlib"
	"github.com/pressly/goose"
)

type Store struct {
	conn       *sql.DB
	dbGetter   stdlibTransactor.DBGetter
	transactor transactor.Transactor
}

func NewStore(conn *sql.DB) (*Store, error) {
	transactor, dbGetter := stdlibTransactor.NewTransactor(
		conn,
		stdlibTransactor.NestedTransactionsSavepoints,
	)
	s := &Store{conn: conn, dbGetter: dbGetter, transactor: transactor}
	fmt.Println(os.Getwd())
	dir := "internal/migrations/server/pg"
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
