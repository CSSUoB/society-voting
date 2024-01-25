package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
)

type Election struct {
	bun.BaseModel `json:"-"`

	ID          int    `bun:",pk,autoincrement" json:"id"`
	RoleName    string `json:"roleName"`
	Description string `json:"description"`

	IsActive bool `json:"isActive"`
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
	var candidateModels []*ElectionCandidate
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
	if err := db.NewSelect().Model(res).Where("id = ?", id).Scan(context.Background(), res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get Election model by ID: %w", err)
	}
	return res, nil
}

func GetAllElections(x ...bun.IDB) ([]*Election, error) {
	db := fromVariadic(x)
	var res []*Election
	if err := db.NewSelect().Model(&res).Scan(context.Background(), &res); err != nil {
		return nil, fmt.Errorf("get all Elections: %w", err)
	}
	return res, nil
}

func DeleteElectionByID(electionID int, x ...bun.IDB) error {
	db := fromVariadic(x)
	if _, err := db.NewDelete().Model((*Election)(nil)).Where("id = ?", electionID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete Election: %w", err)
	}
	return nil
}

func GetActiveElection(x ...bun.IDB) (*Election, error) {
	db := fromVariadic(x)
	res := new(Election)
	if count, err := db.NewSelect().Model(res).Where("is_active = 1").ScanAndCount(context.Background(), res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get active Election: %w", err)
	} else if count != 1 {
		return nil, fmt.Errorf("database corrupted: expected 0 active elections, found %d", count)
	}
	return res, nil
}
