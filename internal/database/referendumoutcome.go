package database

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type ReferendumOutcome struct {
	bun.BaseModel `json:"-"`

	ID           int `bun:",pk" json:"id"`
	VotesFor     int `bun:",notnull" json:"votesFor"`
	VotesAgainst int `bun:",notnull" json:"votesAgainst"`
	VotesAbstain int `bun:",notnull" json:"votesAbstain"`

	PollOutcome *PollOutcome `bun:"rel:belongs-to,join:id=id" json:"-"`
}

func (e *ReferendumOutcome) Insert(x ...bun.IDB) error {
	db := fromVariadic(x)

	if err := db.NewInsert().Model(e).Returning("id").Scan(context.Background(), &e.ID); err != nil {
		return fmt.Errorf("insert ReferendumOutcome model: %w", err)
	}

	return nil
}

func (e *ReferendumOutcome) Update(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewUpdate().Model(e).Where("id = ?", e.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("update ReferendumOutcome model: %w", err)
	}

	return nil
}

func (e *ReferendumOutcome) Delete(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewDelete().Model(e).Where("id = ?", e.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete ReferendumOutcome model: %w", err)
	}

	return nil
}
