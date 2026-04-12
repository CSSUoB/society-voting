package guildScraper

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/CSSUoB/society-voting/internal/config"
	"github.com/carlmjohnson/requests"
)

var (
	activeSessionToken string
	tokenMutex         sync.RWMutex
	tokenInitOnce      sync.Once
)

func getActiveToken() string {
	tokenInitOnce.Do(func() {
		activeSessionToken = config.Get().Guild.SessionToken
	})

	tokenMutex.RLock()
	defer tokenMutex.RUnlock()
	return activeSessionToken
}

func updateActiveToken(newToken string) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()
	activeSessionToken = newToken
}

func AutoRefreshSessionToken(ctx context.Context) {
	timer := time.NewTicker(time.Minute * 10)
	defer timer.Stop()

	slog.Debug("Automatic session token refresh loop started")

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			if err := refreshSessionToken(ctx); err != nil {
				slog.Error("refresh session token", "error", err)
			}
		}
	}
}

func refreshSessionToken(ctx context.Context) error {
	currentToken := getActiveToken()

	return requests.URL("https://www.guildofstudents.com/profile").
		Cookie(".AspNet.SharedCookie", currentToken).
		AddValidator(func(res *http.Response) error {
			for _, cookie := range res.Cookies() {
				if cookie.Name == ".AspNet.SharedCookie" && cookie.Value != currentToken {
					updateActiveToken(cookie.Value)
					slog.Info("Session token refreshed")
					return nil
				}
			}
			slog.Debug("Session token update check performed but not changed")
			return nil
		}).
		Fetch(ctx)
}
