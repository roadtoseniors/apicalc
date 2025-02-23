package main

import (
	"context"
	"fmt"
	"os"

	"github.com/roadtoseniors/apicalc/internal/agent/application"
	"github.com/roadtoseniors/apicalc/internal/agent/config"
)

func main() {
	cfg, err := config.NewConfigAg()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	app := application.NewApplication(cfg)

	exitCode := app.Run(context.Background())

	os.Exit(exitCode)
}