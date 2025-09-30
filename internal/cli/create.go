package cli

import (
	"context"
	"fmt"

	"github.com/nasik90/gophkeeper/internal/client/app"
	"github.com/nasik90/gophkeeper/internal/common/logger"
	"github.com/nasik90/gophkeeper/internal/common/types"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// LoginCommand возвращает команду для входа.
func CreateCommand() *cobra.Command {
	var (
		key            string
		valueS         string
		valueB         []byte
		masterPassword string
		comment        string
		binaryValue    bool
	)

	var CreateCommandCmd = &cobra.Command{
		Use:   "create",
		Short: "Создание секрета",
		Long:  `Команда create позволяет создать секрет.`,
		Args:  cobra.ExactArgs(0), // Не принимает аргументов
		RunE: func(cmd *cobra.Command, args []string) error {

			// ключ секрета
			if key == "" {
				fmt.Print("Введите ключ секрета: ")
				_, err := fmt.Scanln(&key)
				if err != nil {
					return fmt.Errorf("ошибка ввода ключа секрета: %w", err)
				}
			}

			// значение секрета
			if !binaryValue && valueS == "" {
				var err error
				valueS, err = promptPassword("Введите значение секрета: ")
				if err != nil {
					return fmt.Errorf("ошибка ввода значения секрета: %w", err)
				}
			}

			// Комментарий
			if comment == "" {
				var err error
				_, err = fmt.Scanln(&comment)
				if err != nil {
					return fmt.Errorf("ошибка ввода комментария: %w", err)
				}
			}

			// Мастер пароль
			if masterPassword == "" {
				var err error
				masterPassword, err = promptPassword("Введите мастер пароль: ")
				if err != nil {
					return fmt.Errorf("ошибка ввода пароля: %w", err)
				}
			}

			var value []byte
			if binaryValue {
				value = valueB
			} else {
				value = []byte(valueS)
			}

			appService, _ := app.InitService(masterPassword)
			secretData := &types.SecretData{Key: []byte(key), Value: value, Comment: comment, BinaryValue: binaryValue}
			err := appService.CreateNewSecret(context.Background(), secretData)
			if err != nil {
				logger.Log.Fatal("create secret error", zap.Error(err))
			}

			fmt.Println("✅ Секрет успешно создан!")
			return nil
		},
	}

	// Добавляем флаги
	CreateCommandCmd.Flags().StringVarP(&key, "key", "k", "", "Ключ секрета")
	CreateCommandCmd.Flags().StringVarP(&valueS, "valueS", "v", "", "Значение секрета (строка)")
	CreateCommandCmd.Flags().StringVarP(&masterPassword, "masterPassword", "m", "", "Мастер пароль")
	CreateCommandCmd.Flags().StringVarP(&comment, "comment", "o", "", "Комментарий")
	CreateCommandCmd.Flags().BoolVarP(&binaryValue, "binaryValue", "b", false, "Значение бинарное")
	CreateCommandCmd.Flags().BytesBase64VarP(&valueB, "valueB", "s", []byte(""), "Значение секрета (binary base64)")
	// Можно сделать обязательным
	// loginCmd.MarkFlagRequired("username")

	return CreateCommandCmd
}
