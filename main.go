package main

import (
	"fmt"
	"github.com/CSSUoB/society-voting/internal/config"
	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/CSSUoB/society-voting/internal/discordWebhookNotify"
	"github.com/CSSUoB/society-voting/internal/httpcore"
	"log/slog"
	"os"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Unhandled error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	conf := config.Get()
	if conf.Debug {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}

	if err := database.Migrate(database.Get()); err != nil {
		return fmt.Errorf("migrate dataase: %w", err)
	}

	if conf.Platform.DiscordWebhook.URL != "" {
		go discordWebhookNotify.Run()
	} else {
		slog.Warn("discord webhook event notifier disabled")
	}

	return httpcore.ListenAndServe(conf.HTTP.Address())
}
