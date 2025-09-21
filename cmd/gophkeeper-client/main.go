package main

import (
	"context"

	"github.com/nasik90/gophkeeper/cmd/gophkeeper-client/settings"
	"github.com/nasik90/gophkeeper/internal/client/api"
	"github.com/nasik90/gophkeeper/internal/client/service"
	"github.com/nasik90/gophkeeper/internal/common/logger"
	"go.uber.org/zap"
)

func main() {
	// Иницилизируем настройки
	options := parseOptions()
	// Иницилизируем локальное хранилище
	store := initStore(options)
	// Иницилизируем API клиент
	client := initApiCleint(options)
	// Иницилизируем слой сервиса
	service := service.NewService(client, store)

	// Получим логин/пароль из CLI/TUI
	login := "nasik90"
	password := "my_password"
	// Залогинимся на сервере
	err := service.Login(context.Background(), login, password)
	if err != nil {
		logger.Log.Fatal("login error", zap.Error(err))
	}

	// Проверим наличие токена для шифрования данных, если токена нет, то попросим его ввести (мастер-пароль данных)

	// По команде сделаем graceful shutdown

}

func parseOptions() *settings.Options {
	options := new(settings.Options)
	settings.ParseFlags(options)
	return options
}

func initStore(options *settings.Options) service.Store {
	return nil
}

func initApiCleint(options *settings.Options) *api.Client {
	return api.NewClient(options.BaseURL)
}
