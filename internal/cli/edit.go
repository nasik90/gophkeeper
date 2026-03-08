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

// EditCommand возвращает команду для редактирования секретов.
func EditCommand() *cobra.Command {
	var (
		id             string
		key            string
		valueS         string
		valueB         []byte
		masterPassword string
		comment        string
		deletionMark   bool
		binaryValue    bool
	)

	var EditCommandCmd = &cobra.Command{
		Use:   "edit",
		Short: "Редактирование секрета",
		Long:  `Команда edit позволяет редактировать секрет.`,
		Args:  cobra.ExactArgs(0), // Не принимает аргументов
		RunE: func(cmd *cobra.Command, args []string) error {

			// Если id не указан, запросим его
			if id == "" {
				fmt.Print("Введите id секрета: ")
				_, err := fmt.Scanln(&id)
				if err != nil {
					return fmt.Errorf("ошибка ввода id секрета: %w", err)
				}
			}

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

			app, err := app.NewApp(masterPassword)
			if err != nil {
				logger.Log.Fatal("application initializing error", zap.Error(err))
			}
			secretData := &types.SecretData{Key: []byte(key), Value: value, Comment: comment, BinaryValue: binaryValue}
			err = app.Service.EditSecret(context.Background(), secretData)
			if err != nil {
				logger.Log.Fatal("edit secret error", zap.Error(err))
			}
			err = app.StopApp()
			if err != nil {
				logger.Log.Fatal("application stop error", zap.Error(err))
			}
			fmt.Println("✅ Секрет успешно изменен!")
			return nil
		},
	}

	// Добавляем флаги
	EditCommandCmd.Flags().StringVarP(&id, "id", "i", "", "Идентификатор секрета")
	EditCommandCmd.Flags().StringVarP(&key, "key", "k", "", "Ключ секрета")
	EditCommandCmd.Flags().StringVarP(&valueS, "valueS", "v", "", "Значение секрета (строка)")
	EditCommandCmd.Flags().StringVarP(&masterPassword, "masterPassword", "m", "", "Мастер пароль")
	EditCommandCmd.Flags().StringVarP(&comment, "comment", "o", "", "Комментарий")
	EditCommandCmd.Flags().BoolVarP(&deletionMark, "deletionMark", "d", false, "Пометка удаления")
	EditCommandCmd.Flags().BoolVarP(&binaryValue, "binaryValue", "b", false, "Значение бинарное")
	EditCommandCmd.Flags().BytesBase64VarP(&valueB, "valueB", "s", []byte(""), "Значение секрета (binary base64)")
	// Можно сделать обязательным
	// EditCommandCmd.MarkFlagRequired("id")
	// EditCommandCmd.MarkFlagRequired("masterPassword")

	return EditCommandCmd
}
