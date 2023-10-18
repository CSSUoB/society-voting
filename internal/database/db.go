package database

import (
	"context"
	"database/sql"
	"errors"
	"git.tdpain.net/codemicro/society-voting/internal/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"
	"log/slog"
	"os"
	"sync"
	"time"
)

var Migrations = migrate.NewMigrations()

var ErrNotFound = errors.New("not found")

type DB struct {
	DB *bun.DB
}

var (
	datab    *DB
	loadOnce = new(sync.Once)
)

func Get() *DB {
	var outerErr error
	loadOnce.Do(func() {
		conf := config.Get().Database

		dsn := conf.DSN
		slog.Info("connecting to database")
		db, err := sql.Open("sqlite3", dsn)
		if err != nil {
			outerErr = err
			return
		}

		db.SetMaxOpenConns(1) // https://github.com/mattn/go-sqlite3/issues/274#issuecomment-191597862

		b := bun.NewDB(db, sqlitedialect.New())

		if config.Get().Debug {
			b.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
		}

		datab = &DB{b}
	})

	if outerErr != nil {
		slog.Error("fatal error when loading configuration", "err", outerErr)
		os.Exit(1)
	}

	return datab
}

func (db *DB) Migrate() error {
	slog.Info("running database migrations")

	mig := migrate.NewMigrator(db.DB, Migrations)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := mig.Init(ctx); err != nil {
		return err
	}

	group, err := mig.Migrate(ctx)
	if err != nil {
		return err
	}

	if group.IsZero() {
		slog.Info("no migrations applied (database up to date)")
	} else {
		slog.Info("migrations applied")
	}

	return nil
}
