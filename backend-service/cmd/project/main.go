package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"semen_project/internal/app"
	"semen_project/internal/config"
)

func main() {
	slog.Info("application is starting...")

	if err := godotenv.Load(); err != nil {
		slog.Error("no .env file found", "error", err)
	}
	cfg := &config.Config{}
	err := config.Load(cfg)
	if err != nil {
		slog.Error("can not load app config", "error", err)
		panic(err)
	}

	slog.Info("configuration loaded")

	if err := app.Run(cfg); err != nil {
		slog.Error("fail to start application", "error", err)
		os.Exit(1)
	}
}
