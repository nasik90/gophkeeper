package cli

import (
	"context"
	"fmt"

	"github.com/nasik90/gophkeeper/internal/client/app"
	"github.com/nasik90/gophkeeper/internal/common/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// LoginCommand возвращает команду для входа.
func LoginCommand() *cobra.Command {
	var (
		username string
		password string
	)

	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Войти в систему",
		Long:  `Команда login позволяет пользователю аутентифицироваться на сервере.`,
		Args:  cobra.ExactArgs(0), // Не принимает аргументов
		RunE: func(cmd *cobra.Command, args []string) error {
			// Если username не указан, запросим его
			if username == "" {
				fmt.Print("Введите имя пользователя: ")
				_, err := fmt.Scanln(&username)
				if err != nil {
					return fmt.Errorf("ошибка ввода имени пользователя: %w", err)
				}
			}

			// Если пароль не указан флагом, запросим его интерактивно
			if password == "" {
				var err error
				password, err = promptPassword("Введите пароль: ")
				if err != nil {
					return fmt.Errorf("ошибка ввода пароля: %w", err)
				}
			}
			app, err := app.NewApp("")
			if err != nil {
				logger.Log.Fatal("application initializing error", zap.Error(err))
			}
			err = app.Service.Login(context.Background(), username, password)
			if err != nil {
				logger.Log.Fatal("login error", zap.Error(err))
			}
			err = app.StopApp()
			if err != nil {
				logger.Log.Fatal("application stop error", zap.Error(err))
			}
			fmt.Println("✅ Аутентификация успешно завершена!")
			return nil
		},
	}

	// Добавляем специфичные для команды login флаги
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "Имя пользователя")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Пароль для аутентификации (не рекомендуется использовать)")
	// Можно сделать обязательным
	// loginCmd.MarkFlagRequired("username")

	return loginCmd
}
