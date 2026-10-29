package main

import (
	"log/slog"
	"os"

	sitebuild "github.com/yash27062004-lab/yash27062004-lab.github.io/internal/build"
	"github.com/yash27062004-lab/yash27062004-lab.github.io/internal/config"
)

func main() {
	logger := slog.New(
		slog.NewTextHandler(
			os.Stderr,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)

	slog.SetDefault(logger)

	if err := run(); err != nil {
		slog.Error("website build failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	generator := sitebuild.New(cfg)

	return generator.Build()
}
