package util

import (
	"log/slog"
)

func Warn(err error) {
	if err != nil {
		slog.Warn("error in deferred function", "error", err)
	}
}
