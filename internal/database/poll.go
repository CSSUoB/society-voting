package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
)

type Poll struct {
	bun.BaseModel `json:"-"`

	ID          int  `bun:",pk,autoincrement" json:"id"`
	PollTypeID  int  `bun:",notnull" json:"-"`
	IsActive    bool `bun:",notnull" json:"isActive"`
	IsConcluded bool `bun:",notnull" json:"isConcluded"`

	PollType   *PollType   `bun:"rel:has-one,join:poll_type_id=id" json:"pollType"`
	Election   *Election   `bun:"rel:has-one,join:id=id" json:"election,omitempty"`
	Referendum *Referendum `bun:"rel:has-one,join:id=id" json:"referendum,omitempty"`
}

type Pollable interface {
	GetPoll() *Poll
	GetFriendlyTitle() string
	GetElection() *Election
	GetReferendum() *Referendum
}

func (p *Poll) Insert(x ...bun.IDB) error {
	db := fromVariadic(x)

	if err := db.NewInsert().Model(p).Returning("id").Scan(context.Background(), &p.ID); err != nil {
		return fmt.Errorf("insert Poll model: %w", err)
	}

	return nil
}

func (p *Poll) Update(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewUpdate().Model(p).Where("id = ?", p.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("update Poll model: %w", err)
	}

	return nil
}

func (p *Poll) Delete(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewDelete().Model(p).Where("id = ?", p.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete Poll model: %w", err)
	}

	return nil
}

func GetPoll(id int, x ...bun.IDB) (*Poll, error) {
	db := fromVariadic(x)
	res := new(Poll)
	if err := db.NewSelect().Model(res).Where("id = ?", id).Relation("PollType").Relation("Election").Relation("Referendum").Scan(context.Background(), res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get Poll model by ID: %w", err)
	}
	return res, nil
}

func GetAllPolls(x ...bun.IDB) ([]*Poll, error) {
	db := fromVariadic(x)
	var res []*Poll
	if err := db.NewSelect().Model(&res).Relation("PollType").Relation("Election").Relation("Referendum").Scan(context.Background(), &res); err != nil {
		return nil, fmt.Errorf("get all Polls: %w", err)
	}
	return res, nil
}

func DeletePollByID(pollID int, x ...bun.IDB) error {
	db := fromVariadic(x)
	if _, err := db.NewDelete().Model((*Poll)(nil)).Where("id = ?", pollID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete Poll: %w", err)
	}
	return nil
}

func GetActivePoll(x ...bun.IDB) (*Poll, error) {
	db := fromVariadic(x)
	res := new(Poll)
	if count, err := db.NewSelect().Model(res).Where("is_active = 1").Relation("PollType").Relation("Election").Relation("Referendum").ScanAndCount(context.Background(), res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get active Poll: %w", err)
	} else if count != 1 {
		return nil, fmt.Errorf("database corrupted: expected 0 active elections, found %d", count)
	}
	return res, nil
}
