package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// LoginCommand возвращает команду для входа.
func LoginCommand() *cobra.Command {
	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Войти в систему",
		Long:  `Команда login позволяет пользователю аутентифицироваться на сервере.`,
		Args:  cobra.ExactArgs(0), // Не принимает аргументов
		Run: func(cmd *cobra.Command, args []string) {
			// Здесь основная логика команды
			fmt.Println("Выполняется логин...")
			// 1. Получить логин/пароль (из флагов, аргументов или интерактивно)
			// 2. Вызвать метод клиента из internal/client
			// 3. Сохранить полученный токен
		},
	}

	// Добавляем специфичные для команды login флаги
	loginCmd.Flags().StringP("username", "u", "", "Имя пользователя")
	loginCmd.Flags().StringP("password", "p", "", "Пароль")
	// Можно сделать обязательным
	// loginCmd.MarkFlagRequired("username")

	return loginCmd
}
