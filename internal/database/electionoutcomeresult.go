package database

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type ElectionOutcomeResult struct {
	bun.BaseModel `json:"-"`

	ID                int    `bun:",pk,autoincrement" json:"id"`
	Name              string `bun:",notnull" json:"name"`
	Round             int    `bun:",notnull" json:"round"`
	Votes             int    `bun:",notnull" json:"voteCount"`
	IsRejected        bool   `bun:",notnull" json:"isRejected"`
	IsElected         bool   `bun:",notnull" json:"isElected"`
	ElectionOutcomeID int    `bun:",notnull" json:"-"`

	ElectionOutcome *ElectionOutcome `bun:"rel:belongs-to,join:election_outcome_id=id" json:"-"`
}

func BulkInsertElectionOutcomeResult(entities []*ElectionOutcomeResult, x ...bun.IDB) error {
	db := fromVariadic(x)

	if err := db.NewInsert().Model(&entities).Scan(context.Background()); err != nil {
		return fmt.Errorf("bulk insert ElectionOutcomeResult model: %w", err)
	}

	return nil
}
