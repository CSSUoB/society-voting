package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/CSSUoB/society-voting/internal/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"
)

//go:embed *.sql
var sqlMigrations embed.FS

var Migrations = migrate.NewMigrations()

var ErrNotFound = errors.New("not found")

type DB struct {
	DB bun.IDB
}

var (
	datab    *bun.DB
	loadOnce = new(sync.Once)
)

func Get() *bun.DB {
	var outerErr error
	loadOnce.Do(func() {
		conf := config.Get().Database

		dsn := conf.DSN + "?_txlock=exclusive&_foreign_keys=on"
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

		datab = b
	})

	if outerErr != nil {
		slog.Error("fatal error when loading configuration", "err", outerErr)
		os.Exit(1)
	}

	return datab
}

func GetTx() (bun.Tx, error) {
	db := Get()
	return db.Begin()
}

func fromVariadic(x []bun.IDB) bun.IDB {
	if len(x) == 0 {
		return Get()
	}
	return x[0]
}

func Migrate(db *bun.DB) error {
	slog.Info("running database migrations")

	if err := Migrations.Discover(sqlMigrations); err != nil {
		return err
	}

	mig := migrate.NewMigrator(db, Migrations)

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
