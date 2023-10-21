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

type ElectionWithCandidates struct {
	Election
	Candidates []string `json:"candidates"`
}

func (e *Election) Insert() error {
	db := Get()

	if err := db.DB.NewInsert().Model(e).Returning("id").Scan(context.Background(), &e.ID); err != nil {
		return fmt.Errorf("insert Election model: %w", err)
	}

	return nil
}

func (e *Election) Update() error {
	db := Get()

	if _, err := db.DB.NewUpdate().Model(e).Where("id = ?", e.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("update Election model: %w", err)
	}

	return nil
}

func (e *Election) Delete() error {
	db := Get()

	if _, err := db.DB.NewDelete().Model(e).Where("id = ?", e.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete Election model: %w", err)
	}

	return nil
}

func (e *Election) WithCandidates() (*ElectionWithCandidates, error) {
	candidates, err := GetUsersStandingForElection(e.ID)
	if err != nil {
		return nil, fmt.Errorf("populate Election candidates: %w", err)
	}
	var candidateStrings []string
	for _, cand := range candidates {
		candidateStrings = append(candidateStrings, cand.Name)
	}
	return &ElectionWithCandidates{
		Election:   *e,
		Candidates: candidateStrings,
	}, nil
}

func GetElection(id int) (*Election, error) {
	db := Get()
	res := new(Election)
	if err := db.DB.NewSelect().Model(res).Where("id = ?", id).Scan(context.Background(), res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get Election model by ID: %w", err)
	}
	return res, nil
}

func GetAllElections() ([]*Election, error) {
	db := Get()
	var res []*Election
	if err := db.DB.NewSelect().Model(&res).Scan(context.Background(), &res); err != nil {
		return nil, fmt.Errorf("get all Elections: %w", err)
	}
	return res, nil
}

func DeleteElectionByID(electionID int) error {
	db := Get()
	if _, err := db.DB.NewDelete().Model((*Election)(nil)).Where("id = ?", electionID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete Election: %w", err)
	}
	return nil
}

func GetActiveElection() (*Election, error) {
	db := Get()
	res := new(Election)
	if count, err := db.DB.NewSelect().Model(res).Where("is_active = 1").ScanAndCount(context.Background(), res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get active Election: %w", err)
	} else if count != 1 {
		return nil, fmt.Errorf("database corrupted: expected 0 active elections, found %d", count)
	}
	return res, nil
}
