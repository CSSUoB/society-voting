package database

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(
		func(ctx context.Context, db *bun.DB) error {
			if _, err := db.NewCreateTable().
				Model((*ElectionOutcome)(nil)).
				ForeignKey("(election_id) REFERENCES elections (id) ON DELETE CASCADE").
				Exec(ctx); err != nil {
				return fmt.Errorf("create ElectionOutcome table: %w", err)
			}

			if _, err := db.NewCreateTable().
				Model((*ElectionOutcomeResult)(nil)).
				ForeignKey("(election_outcome_id) REFERENCES election_outcomes (id) ON DELETE CASCADE").
				Exec(ctx); err != nil {
				return fmt.Errorf("create ElectionOutcomeResult table: %w", err)
			}

			return nil
		},
		func(ctx context.Context, db *bun.DB) error {
			if _, err := db.NewDropTable().Model((*ElectionOutcome)(nil)).Exec(ctx); err != nil {
				return fmt.Errorf("drop ElectionOutcome table: %w", err)
			}

			if _, err := db.NewDropTable().Model((*ElectionOutcomeResult)(nil)).Exec(ctx); err != nil {
				return fmt.Errorf("drop ElectionOutcomeResult table: %w", err)
			}

			return nil
		},
	)
}
