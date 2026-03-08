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

	// app, err := app.NewApp("")
	// if err != nil {
	// 	logger.Log.Fatal("application initializing error", zap.Error(err))
	// }
	// err = app.Service.RegisterNewUser(context.Background(), "username", "password")
	// if err != nil {
	// 	logger.Log.Fatal("register error", zap.Error(err))
	// }

}
