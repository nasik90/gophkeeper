package cli

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/nasik90/gophkeeper/internal/client/app"
	"github.com/nasik90/gophkeeper/internal/common/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// GetCommand возвращает команду для получения списка секретов.
func GetCommand() *cobra.Command {
	var (
		masterPassword string
	)

	var GetCommandCmd = &cobra.Command{
		Use:   "get",
		Short: "Получение секретов",
		Long:  `Поулчение секретов из локальной базы.`,
		Args:  cobra.ExactArgs(0), // Не принимает аргументов
		RunE: func(cmd *cobra.Command, args []string) error {

			// Мастер пароль
			if masterPassword == "" {
				var err error
				masterPassword, err = promptPassword("Введите мастер пароль: ")
				if err != nil {
					return fmt.Errorf("ошибка ввода пароля: %w", err)
				}
			}
			app, err := app.NewApp(masterPassword)
			if err != nil {
				logger.Log.Fatal("application initializing error", zap.Error(err))
			}
			secrets, err := app.Service.GetSecrets(context.Background())
			if err != nil {
				logger.Log.Fatal("get secrets error", zap.Error(err))
			}
			value := ""
			for _, secret := range *secrets {
				if secret.BinaryValue {
					value = base64.StdEncoding.EncodeToString(secret.Value)
				} else {
					value = string(secret.Value)
				}
				fmt.Println(secret.Guid, string(secret.Key), value, secret.Comment)
			}
			return nil
		},
	}

	GetCommandCmd.Flags().StringVarP(&masterPassword, "masterPassword", "m", "", "Мастер пароль")

	return GetCommandCmd
}
