package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/pressly/goose"

	"github.com/Thiht/transactor"
	stdlibTransactor "github.com/Thiht/transactor/stdlib"
)

type Store struct {
	conn       *sql.DB
	dbGetter   stdlibTransactor.DBGetter
	transactor transactor.Transactor
}

//	func NewStore(conn *sql.DB) (*Store, error) {
//		s := &Store{conn: conn}
func NewStore(conn *sql.DB) (*Store, transactor.Transactor, error) {
	transactor, dbGetter := stdlibTransactor.NewTransactor(
		conn,
		stdlibTransactor.NestedTransactionsSavepoints,
	)
	s := &Store{conn: conn, dbGetter: dbGetter, transactor: transactor}
	//dir := "C:/golang_projects/gophkeeper/internal/migrations/server/pg"
	fmt.Println(os.Getwd())
	dir := "/internal/migrations/server/pg"
	//dir := "C:\golang_project\gophkeeper\internal\migrations\server\pg"
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
