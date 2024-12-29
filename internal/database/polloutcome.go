package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type PollOutcome struct {
	bun.BaseModel `json:"-"`

	ID          int       `bun:",pk,autoincrement" json:"id"`
	PollID      int       `bun:",notnull,unique" json:"-"`
	Date        time.Time `bun:",notnull,default:current_timestamp" json:"date"`
	Ballots     int       `bun:",notnull" json:"ballots"`
	IsPublished bool      `bun:",notnull" json:"isPublished"`

	Poll              *Poll              `bun:"rel:belongs-to,join:poll_id=id" json:"poll"`
	ElectionOutcome   *ElectionOutcome   `bun:"rel:has-one,join:id=id" json:"electionOutcome,omitempty"`
	ReferendumOutcome *ReferendumOutcome `bun:"rel:has-one,join:id=id" json:"referendumOutcome,omitempty"`
}

func (e *PollOutcome) Insert(x ...bun.IDB) error {
	db := fromVariadic(x)

	if err := db.NewInsert().Model(e).Returning("id").Scan(context.Background(), &e.ID); err != nil {
		return fmt.Errorf("insert PollOutcome model: %w", err)
	}

	return nil
}

func (e *PollOutcome) Update(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewUpdate().Model(e).Where("id = ?", e.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("update PollOutcome model: %w", err)
	}

	return nil
}

func (e *PollOutcome) Delete(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewDelete().Model(e).Where("id = ?", e.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete PollOutcome model: %w", err)
	}

	return nil
}

func GetOutcomeForPoll(id int, x ...bun.IDB) (*PollOutcome, error) {
	db := fromVariadic(x)
	res := new(PollOutcome)
	if err := db.NewSelect().Model(res).
		Where("poll_id = ?", id).
		Relation("Poll").
		Relation("Poll.PollType").
		Relation("Poll.Election").
		Relation("Poll.Referendum").
		Relation("ElectionOutcome").
		Relation("ElectionOutcome.Results").
		Relation("ReferendumOutcome").
		Scan(context.Background()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get PollOutcome model by poll ID: %w", err)
	}
	return res, nil
}

func PublishPollOutcome(id int, x ...bun.IDB) error {
	db := fromVariadic(x)
	if _, err := db.NewUpdate().Table("poll_outcomes").Set("published = ?", true).Where("poll_id = ?", id).Exec(context.Background()); err != nil {
		return fmt.Errorf("publish PollOutcome: %w", err)
	}

	return nil
}
