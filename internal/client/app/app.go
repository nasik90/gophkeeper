package app

import (
	"database/sql"
	"errors"

	"github.com/nasik90/gophkeeper/internal/client/api"
	"github.com/nasik90/gophkeeper/internal/client/service"
	"github.com/nasik90/gophkeeper/internal/client/settings"
	"github.com/nasik90/gophkeeper/internal/client/storage/sqlite"
	"github.com/nasik90/gophkeeper/internal/common/logger"
	_ "modernc.org/sqlite"
)

func InitService(masterPassword string) (*service.Service, error) {
	// Иницилизируем настройки
	options := parseOptions()
	if err := logger.Initialize(options.LogLevel); err != nil {
		panic(err)
	}
	// Иницилизируем локальное хранилище
	store, err := initStore(options)
	if err != nil {
		return nil, err
	}
	// Иницилизируем API клиент
	client := initApiCleint(options)
	// Иницилизируем слой сервиса
	appService := service.NewService(client, store, masterPassword)

	return appService, err

}

func parseOptions() *settings.Options {
	options := new(settings.Options)
	settings.ParseFlags(options)
	return options
}

func initStore(options *settings.Options) (service.Store, error) {
	var store *sqlite.Store
	//var store service.Store

	if options.DatabaseDSN == "" {
		return store, errors.New("can not initialize store: db settings are empty")
	}
	conn, err := sql.Open("sqlite", options.DatabaseDSN)
	if err != nil {
		return store, err
	}
	store, err = sqlite.NewStore(conn)
	if err != nil {
		return store, err
	}

	return store, nil
}

func initApiCleint(options *settings.Options) *api.Client {
	return api.NewClient(options.BaseURL)
}

// TODO: написать остановку
func StopApp() {

}
