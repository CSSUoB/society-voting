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

func (r *ReferendumOutcome) Insert(x ...bun.IDB) error {
	db := fromVariadic(x)

	if err := db.NewInsert().Model(r).Returning("id").Scan(context.Background(), &r.ID); err != nil {
		return fmt.Errorf("insert ReferendumOutcome model: %w", err)
	}

	return nil
}

func (r *ReferendumOutcome) Update(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewUpdate().Model(r).Where("id = ?", r.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("update ReferendumOutcome model: %w", err)
	}

	return nil
}

func (r *ReferendumOutcome) Delete(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewDelete().Model(r).Where("id = ?", r.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete ReferendumOutcome model: %w", err)
	}

	return nil
}
