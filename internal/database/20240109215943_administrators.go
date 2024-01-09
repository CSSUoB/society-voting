package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
	"log/slog"
)

func init() {
	Migrations.MustRegister(
		func(ctx context.Context, db *bun.DB) error {
			if _, err := db.NewRaw(`ALTER TABLE "users" ADD COLUMN "is_admin" BOOLEAN;`).Exec(ctx); err != nil {
				var e sqlite3.Error
				if errors.As(err, &e) {
					if err.(sqlite3.Error).Code == 1 {
						// assume this is due to a duplicate column
						return nil
					}
				}
				return fmt.Errorf("alter User table: %w", err)
			}

			slog.Warn("If you have just upgraded society-voting from a previous version, you will now need to manually edit the database to create some administrator users.")

			return nil
		},
		func(ctx context.Context, db *bun.DB) error {
			return nil
		},
	)
}
