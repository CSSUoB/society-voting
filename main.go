package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/CSSUoB/society-voting/internal/config"
	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/CSSUoB/society-voting/internal/discordWebhookNotify"
	"github.com/CSSUoB/society-voting/internal/events"
	"github.com/CSSUoB/society-voting/internal/httpcore"
)

var exitingForRestart = false

func main() {
	if err := run(); err != nil {
		slog.Error("Unhandled error", "error", err)
		os.Exit(1)
	}

	if exitingForRestart {
		os.Exit(153)
	}
}

func run() error {
	slog.Info("starting society-voting")

	conf := config.Get()
	if conf.Debug {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}

	if err := database.Migrate(database.Get()); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	if conf.Platform.DiscordWebhook.URL != "" {
		go discordWebhookNotify.Run()
	} else {
		slog.Warn("discord webhook event notifier disabled")
	}

	httpcore.InitialiseSigner(conf.Platform.SessionSigningToken)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if os.Getenv("SOCIETY_VOTING_RESTART_ENABLED") != "" {
		slog.Info("restart shim enabled")
		_, pollEndReceiver := events.NewReceiver(events.TopicPollEnded)
		go func() {
			<-pollEndReceiver
			slog.Info("poll ended - restarting")
			exitingForRestart = true
			cancel()
		}()
	}

	return httpcore.ListenAndServe(ctx, conf.HTTP.Address())
}
