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

func (c *Candidate) Insert(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewInsert().Model(c).Exec(context.Background()); err != nil {
		return fmt.Errorf("insert Candidate model: %w", err)
	}

	return nil
}

func (c *Candidate) Delete(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewDelete().Model(c).Where("user_id = ? and election_id = ?", c.UserID, c.ElectionID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete Candidate model: %w", err)
	}

	return nil
}

func GetUsersStandingForElection(electionID int, x ...bun.IDB) ([]*User, error) {
	db := fromVariadic(x)
	var res []*User
	if err := db.NewSelect().Model(&res).Where(`id IN (SELECT "user_id" FROM "candidates" WHERE "election_id" = ?)`, electionID).Scan(context.Background(), &res); err != nil {
		return nil, fmt.Errorf("get all Elections: %w", err)
	}
	return res, nil
}

func DeleteCandidatesForElection(electionID int, x ...bun.IDB) error {
	db := fromVariadic(x)
	if _, err := db.NewDelete().Model((*Candidate)(nil)).Where("election_id = ?", electionID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete Candidates for election: %w", err)
	}
	return nil
}
