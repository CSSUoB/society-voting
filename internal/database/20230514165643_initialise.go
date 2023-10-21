package database

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(
		func(ctx context.Context, db *bun.DB) error {
			if _, err := db.NewCreateTable().Model((*User)(nil)).Exec(ctx); err != nil {
				return fmt.Errorf("create User table: %w", err)
			}

			if _, err := db.NewCreateTable().Model((*Election)(nil)).Exec(ctx); err != nil {
				return fmt.Errorf("create Election table: %w", err)
			}

			if _, err := db.NewRaw(`CREATE TABLE "candidates" ("user_id" VARCHAR, "election_id" INTEGER, PRIMARY KEY ("user_id", "election_id"))`).Exec(ctx); err != nil {
				return fmt.Errorf("create Candidate table: %w", err)
			}

			if _, err := db.NewCreateTable().Model((*BallotEntry)(nil)).Exec(ctx); err != nil {
				return fmt.Errorf("create BallotEntry table: %w", err)
			}

			if _, err := db.NewCreateTable().Model((*Vote)(nil)).Exec(ctx); err != nil {
				return fmt.Errorf("create Vote table: %w", err)
			}

			return nil
		},
		func(ctx context.Context, db *bun.DB) error {
			if _, err := db.NewDropTable().Model((*Vote)(nil)).Exec(ctx); err != nil {
				return fmt.Errorf("drop Vote table: %w", err)
			}

			if _, err := db.NewDropTable().Model((*BallotEntry)(nil)).Exec(ctx); err != nil {
				return fmt.Errorf("drop BallotEntry table: %w", err)
			}

			if _, err := db.NewDropTable().Model((*Candidate)(nil)).Exec(ctx); err != nil {
				return fmt.Errorf("drop Candidate table: %w", err)
			}

			if _, err := db.NewDropTable().Model((*Election)(nil)).Exec(ctx); err != nil {
				return fmt.Errorf("drop Election table: %w", err)
			}

			if _, err := db.NewDropTable().Model((*User)(nil)).Exec(ctx); err != nil {
				return fmt.Errorf("drop User table: %w", err)
			}

			return nil
		},
	)
}
