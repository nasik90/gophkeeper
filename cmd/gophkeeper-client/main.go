package main

import (
	"context"

	api "github.com/nasik90/gophkeeper/internal/client/api"
	"github.com/nasik90/gophkeeper/internal/client/service"
	"github.com/nasik90/gophkeeper/internal/common/logger"
	"go.uber.org/zap"
)

func main() {
	client := api.NewClient("http://localhost:8080")
	service := service.NewService(client)
	err := service.Login(context.Background(), "nasik90", "my_password")
	if err != nil {
		logger.Log.Fatal("login error", zap.Error(err))
	}
	service.CreateNewSecret(context.Background(), "key2", "value2", "comment2")
}
