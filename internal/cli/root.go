package cli

import (
	"github.com/nasik90/gophkeeper/cmd/gophkeeper-client-cli/settings"
	"github.com/nasik90/gophkeeper/internal/client/api"
	"github.com/nasik90/gophkeeper/internal/client/service"
	"github.com/spf13/cobra"
)

var appService *service.Service

// NewRootCommand создает и возвращает корневую команду
func RootCommand() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "gophkeeper",
		Short: "GophKeeper - безопасное хранилище ваших паролей и данных",
		Long:  `GophKeeper позволяет безопасно хранить и синхронизировать ваши пароли, заметки и банковские карты.`,
		// PersistentPreRun можно использовать для инициализации конфига, логгера и т.д.
		// Эта функция выполнится ДЛЯ ЛЮБОЙ команды.
	}

	// Здесь регистрируем Persistent Flags (общие для всех команд)
	// Например, флаг для конфигурационного файла
	rootCmd.PersistentFlags().StringP("config", "c", "", "Конфигурационный файл (default is $HOME/.gophkeeper.yaml)")

	// Здесь добавляем все дочерние команды
	rootCmd.AddCommand(LoginCommand())
	rootCmd.AddCommand(RegisterCommand())
	rootCmd.AddCommand(CreateSecretCommand())
	// rootCmd.AddCommand(NewGetCommand())
	// rootCmd.AddCommand(NewListCommand())

	return rootCmd
}

func initService() {
	// Иницилизируем настройки
	options := parseOptions()
	// Иницилизируем локальное хранилище
	store := initStore(options)
	// Иницилизируем API клиент
	client := initApiCleint(options)
	// Иницилизируем слой сервиса
	appService = service.NewService(client, store)

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
