package database

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type ElectionOutcome struct {
	bun.BaseModel `json:"-"`

	ID     int `bun:",pk" json:"id"`
	Rounds int `bun:",notnull" json:"rounds"`

	PollOutcome *PollOutcome             `bun:"rel:belongs-to,join:id=id" json:"-"`
	Results     []*ElectionOutcomeResult `bun:"rel:has-many,join:id=election_outcome_id" json:"results"`
}

func (e *ElectionOutcome) Insert(x ...bun.IDB) error {
	db := fromVariadic(x)

	if err := db.NewInsert().Model(e).Returning("id").Scan(context.Background(), &e.ID); err != nil {
		return fmt.Errorf("insert ElectionOutcome model: %w", err)
	}

	return nil
}

func (e *ElectionOutcome) Update(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewUpdate().Model(e).Where("id = ?", e.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("update ElectionOutcome model: %w", err)
	}

	return nil
}

func (e *ElectionOutcome) Delete(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewDelete().Model(e).Where("id = ?", e.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete ElectionOutcome model: %w", err)
	}

	return nil
}
