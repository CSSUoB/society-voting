package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
)

var _ Pollable = (*Referendum)(nil)

type Referendum struct {
	bun.BaseModel `json:"-"`

	ID          int    `bun:",pk" json:"id"`
	Title       string `bun:",notnull" json:"title"`
	Question    string `bun:",notnull" json:"question"`
	Description string `bun:",notnull" json:"description"`

	Poll *Poll `bun:"rel:belongs-to,join:id=id" json:"-"`
}

func (r *Referendum) GetPoll() *Poll {
	return r.Poll
}

func (r *Referendum) GetFriendlyTitle() string {
	return r.Title
}

func (r *Referendum) GetElection() *Election {
	return nil
}

func (r *Referendum) GetReferendum() *Referendum {
	return r
}

func (r *Referendum) Insert(x ...bun.IDB) error {
	db := fromVariadic(x)

	if err := db.NewInsert().Model(r).Returning("id").Scan(context.Background(), &r.ID); err != nil {
		return fmt.Errorf("insert Referendum model: %w", err)
	}

	return nil
}

func (r *Referendum) Update(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewUpdate().Model(r).Where("id = ?", r.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("update Referendum model: %w", err)
	}

	return nil
}

func (r *Referendum) Delete(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewDelete().Model(r).Where("id = ?", r.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete Referendum model: %w", err)
	}

	return nil
}

func GetReferendum(id int, x ...bun.IDB) (*Referendum, error) {
	db := fromVariadic(x)
	res := new(Referendum)
	if err := db.NewSelect().Model(res).Where(`"referendum"."id" = ?`, id).Relation("Poll").Scan(context.Background(), res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get Referendum model by ID: %w", err)
	}
	return res, nil
}
