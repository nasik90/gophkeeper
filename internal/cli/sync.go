package cli

import (
	"context"
	"fmt"

	"github.com/nasik90/gophkeeper/internal/client/app"
	"github.com/nasik90/gophkeeper/internal/common/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// SyncCommand возвращает команду для синхронизации данных на клиенте с сервером.
func SyncCommand() *cobra.Command {

	var SyncCommandCmd = &cobra.Command{
		Use:   "sync",
		Short: "Синхронизация данных",
		Long:  `Синхронизация данных между локальными клиентом и сервером.`,
		Args:  cobra.ExactArgs(0), // Не принимает аргументов
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := app.NewApp("")
			if err != nil {
				logger.Log.Fatal("application initializing error", zap.Error(err))
			}
			err = app.Service.SendSecrets(context.Background())
			if err != nil {
				logger.Log.Fatal("send secrets error", zap.Error(err))
			}
			err = app.Service.UploadSecrets(context.Background())
			if err != nil {
				logger.Log.Fatal("upload secrets error", zap.Error(err))
			}
			err = app.StopApp()
			if err != nil {
				logger.Log.Fatal("application stop error", zap.Error(err))
			}
			fmt.Println("✅ Синхронизация успешно завершена!")
			return nil
		},
	}

	return SyncCommandCmd
}
