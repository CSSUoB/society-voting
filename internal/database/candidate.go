package database

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
)

type Candidate struct {
	bun.BaseModel `json:"-"`

	UserID     string `json:"userID"`
	ElectionID int    `json:"electionID"`
}

func (c *Candidate) Insert() error {
	db := Get()

	if _, err := db.DB.NewInsert().Model(c).Exec(context.Background()); err != nil {
		return fmt.Errorf("insert Candidate model: %w", err)
	}

	return nil
}

func (c *Candidate) Delete() error {
	db := Get()

	if _, err := db.DB.NewDelete().Model(c).Where("user_id = ? and election_id = ?", c.UserID, c.ElectionID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete Candidate model: %w", err)
	}

	return nil
}

func GetUsersStandingForElection(electionID int) ([]*User, error) {
	db := Get()
	var res []*User
	if err := db.DB.NewSelect().Model(&res).Where(`id IN (SELECT "user_id" FROM "candidates" WHERE "election_id" = ?)`, electionID).Scan(context.Background(), &res); err != nil {
		return nil, fmt.Errorf("get all Elections: %w", err)
	}
	return res, nil
}
