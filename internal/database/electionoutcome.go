package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type ElectionOutcome struct {
	bun.BaseModel `json:"-"`

	ID          int       `bun:",pk,autoincrement" json:"id"`
	ElectionID  int       `bun:",notnull,unique" json:"-"`
	Date        time.Time `bun:",notnull,default:current_timestamp" json:"date"`
	Ballots     int       `bun:",notnull" json:"ballots"`
	Rounds      int       `bun:",notnull" json:"rounds"`
	IsPublished bool      `bun:",notnull" json:"isPublished"`

	Election *Election                `bun:"rel:belongs-to,join:election_id=id" json:"election"`
	Results  []*ElectionOutcomeResult `bun:"rel:has-many,join:id=election_outcome_id" json:"results"`
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

func GetOutcomeForElection(id int, x ...bun.IDB) (*ElectionOutcome, error) {
	db := fromVariadic(x)
	res := new(ElectionOutcome)
	if err := db.NewSelect().Model(res).Relation("Election").Where("election_id = ?", id).Relation("Results").Scan(context.Background()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get ElectionOutcome model by election ID: %w", err)
	}
	return res, nil
}

func PublishElectionOutcome(id int, x ...bun.IDB) error {
	db := fromVariadic(x)
	if _, err := db.NewUpdate().Table("election_outcomes").Set("published = ?", true).Where("election_id = ?", id).Exec(context.Background()); err != nil {
		return fmt.Errorf("publish ElectionOutcome: %w", err)
	}

	return nil
}
