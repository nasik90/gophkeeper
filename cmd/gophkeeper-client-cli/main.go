package main

import (
	"os"

	"github.com/nasik90/gophkeeper/internal/cli"
)

func main() {
	cmd := cli.RootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

	// appService, err := app.InitService("my_password")
	// if err != nil {
	// 	logger.Log.Fatal("start client app", zap.Error(err))
	// }
	// logger.Log.Info("start client app")

	// secretData := &types.SecretData{}
	// secretData.Key = []byte("nnnn3")
	// secretData.Value = []byte("mmm3")
	// err = appService.CreateNewSecret(context.Background(), secretData)
	// if err != nil {
	// 	logger.Log.Fatal("CreateNewSecret", zap.Error(err))
	// }

	// err = appService.SendSecrets(context.Background())
	// if err != nil {
	// 	logger.Log.Fatal("SendSecrets", zap.Error(err))
	// }

	// secrets, err := appService.GetSecrets(context.Background())
	// if err != nil {
	// 	logger.Log.Fatal("GetSecrets", zap.Error(err))
	// }
	// for _, secret := range *secrets {
	// 	fmt.Println(string(secret.Key), string(secret.Value), secret.Comment)
	// }

	// err = appService.UploadSecrets(context.Background())
	// if err != nil {
	// 	logger.Log.Fatal("UploadSecrets", zap.Error(err))
	// }
}
