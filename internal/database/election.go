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

	ID       int    `bun:",pk,autoincrement" json:"id"`
	RoleName string `json:"roleName"`
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
