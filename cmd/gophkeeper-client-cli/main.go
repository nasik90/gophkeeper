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
}
