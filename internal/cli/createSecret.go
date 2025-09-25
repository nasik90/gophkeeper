package cli

import (
	"context"
	"fmt"

	"github.com/nasik90/gophkeeper/internal/common/logger"
	"github.com/nasik90/gophkeeper/internal/common/types"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// LoginCommand возвращает команду для входа.
func CreateSecretCommand() *cobra.Command {
	var (
		key   string
		value string
	)

	var CreateSecretCommandCmd = &cobra.Command{
		Use:   "createSecret",
		Short: "Создание секрета",
		Long:  `Команда createSecret позволяет создать секрет.`,
		Args:  cobra.ExactArgs(0), // Не принимает аргументов
		RunE: func(cmd *cobra.Command, args []string) error {
			// Если username не указан, запросим его
			if key == "" {
				fmt.Print("Введите имя пользователя: ")
				_, err := fmt.Scanln(&key)
				if err != nil {
					return fmt.Errorf("ошибка ввода имени пользователя: %w", err)
				}
			}

			// Если пароль не указан флагом, запросим его интерактивно
			if value == "" {
				var err error
				value, err = promptPassword("Введите пароль: ")
				if err != nil {
					return fmt.Errorf("ошибка ввода пароля: %w", err)
				}
			}
			initService()
			secretData := &types.SecretData{Key: key, Value: value}
			err := appService.CreateNewSecret(context.Background(), secretData)
			if err != nil {
				logger.Log.Fatal("create secret error", zap.Error(err))
			}

			fmt.Println("✅ Секрет успешно создан!")
			return nil
		},
	}

	// Добавляем флаги
	CreateSecretCommandCmd.Flags().StringVarP(&key, "key", "k", "", "Имя")
	CreateSecretCommandCmd.Flags().StringVarP(&value, "value", "v", "", "Пароль")
	// Можно сделать обязательным
	// loginCmd.MarkFlagRequired("username")

	return CreateSecretCommandCmd
}
