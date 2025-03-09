package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
)

var _ Pollable = (*Election)(nil)

type Election struct {
	bun.BaseModel `json:"-"`

	ID          int    `bun:",pk" json:"id"`
	RoleName    string `bun:",notnull" json:"roleName"`
	Description string `bun:",notnull" json:"description"`

	Poll *Poll `bun:"rel:belongs-to,join:id=id" json:"-"`
}

type ElectionCandidate struct {
	Name string `json:"name"`
	ID   string `json:"-"`
	IsMe bool   `json:"isMe"`
}

type ElectionWithCandidates struct {
	Election
	Candidates []*ElectionCandidate `json:"candidates"`
}

func (e *Election) GetPoll() *Poll {
	return e.Poll
}

func (e *Election) GetFriendlyTitle() string {
	return e.RoleName
}

func (e *Election) GetElection() *Election {
	return e
}

func (e *Election) GetReferendum() *Referendum {
	return nil
}

func (e *Election) Insert(x ...bun.IDB) error {
	db := fromVariadic(x)

	if err := db.NewInsert().Model(e).Returning("id").Scan(context.Background(), &e.ID); err != nil {
		return fmt.Errorf("insert Election model: %w", err)
	}

	return nil
}

func (e *Election) Update(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewUpdate().Model(e).Where("id = ?", e.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("update Election model: %w", err)
	}

	return nil
}

func (e *Election) Delete(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewDelete().Model(e).Where("id = ?", e.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete Election model: %w", err)
	}

	return nil
}

func (e *Election) WithCandidates(x ...bun.IDB) (*ElectionWithCandidates, error) {
	candidates, err := GetUsersStandingForElection(e.ID)
	if err != nil {
		return nil, fmt.Errorf("populate Election candidates: %w", err)
	}
	candidateModels := make([]*ElectionCandidate, 0)
	for _, cand := range candidates {
		candidateModels = append(candidateModels, &ElectionCandidate{Name: cand.Name, ID: cand.StudentID})
	}
	return &ElectionWithCandidates{
		Election:   *e,
		Candidates: candidateModels,
	}, nil
}

func GetElection(id int, x ...bun.IDB) (*Election, error) {
	db := fromVariadic(x)
	res := new(Election)
	if err := db.NewSelect().Model(res).Where(`"election"."id" = ?`, id).Relation("Poll").Scan(context.Background(), res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get Election model by ID: %w", err)
	}
	return res, nil
}
