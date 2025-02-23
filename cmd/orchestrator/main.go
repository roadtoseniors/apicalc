package main

import (
	"context"
	"fmt"
	"os"

	"github.com/roadtoseniors/apicalc/internal/orchestrator/application"
	"github.com/roadtoseniors/apicalc/internal/orchestrator/config"
)

func main() {
	cfg, err := config.NewConfigOrch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	app := application.NewApplication(cfg)
	exitCode := app.Run(ctx)

	os.Exit(exitCode)
}
