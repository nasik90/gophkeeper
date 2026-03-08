package cli

import (
	"context"
	"fmt"
	"syscall"

	"github.com/nasik90/gophkeeper/internal/client/app"
	"github.com/nasik90/gophkeeper/internal/common/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/term"
)

// RegisterCommand возвращает команду для регистрации нового пользователя.
func RegisterCommand() *cobra.Command {
	var (
		username string
		password string
	)

	var registerCmd = &cobra.Command{
		Use:   "register",
		Short: "Зарегистрировать нового пользователя",
		Long: `Команда register создает новую учетную запись в системе GophKeeper.
Для безопасности рекомендуется не передавать пароль аргументом, а ввести его интерактивно.`,
		Args: cobra.ExactArgs(0),
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

			// Валидация входных данных
			if err := validateCredentials(username, password); err != nil {
				return err
			}

			// Вызов логики регистрации
			app, err := app.NewApp("")
			if err != nil {
				logger.Log.Fatal("application initializing error", zap.Error(err))
			}
			err = app.Service.RegisterNewUser(context.Background(), username, password)
			if err != nil {
				logger.Log.Fatal("register error", zap.Error(err))
			}
			err = app.StopApp()
			if err != nil {
				logger.Log.Fatal("application stop error", zap.Error(err))
			}
			fmt.Println("✅ Регистрация успешно завершена!")
			return nil
		},
	}

	// Добавляем флаги
	registerCmd.Flags().StringVarP(&username, "username", "u", "", "Имя пользователя для регистрации")
	registerCmd.Flags().StringVarP(&password, "password", "p", "", "Пароль для регистрации (не рекомендуется использовать)")

	return registerCmd
}

// promptPassword запрашивает пароль интерактивно без отображения ввода.
func promptPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}

// validateCredentials проверяет корректность логина и пароля.
func validateCredentials(username, password string) error {
	if username == "" {
		return fmt.Errorf("имя пользователя не может быть пустым")
	}
	if len(username) < 3 {
		return fmt.Errorf("имя пользователя должно содержать минимум 3 символа")
	}
	if password == "" {
		return fmt.Errorf("пароль не может быть пустым")
	}
	if len(password) < 8 {
		return fmt.Errorf("пароль должен содержать минимум 8 символов")
	}
	return nil
}
