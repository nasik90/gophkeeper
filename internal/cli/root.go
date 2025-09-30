package cli

import (
	"github.com/spf13/cobra"
)

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
	rootCmd.AddCommand(CreateCommand())
	rootCmd.AddCommand(SyncCommand())
	rootCmd.AddCommand(GetCommand())
	rootCmd.AddCommand(EditCommand())

	return rootCmd
}
