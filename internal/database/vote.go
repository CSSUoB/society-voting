package database

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
)

type Vote struct {
	bun.BaseModel `json:"-"`

	ID         int    `bun:",pk,autoincrement" json:"id"`
	ElectionID int    `json:"electionID"`
	UserID     string `json:"userID"`
	Choices    []int  `json:"choices"`
}

func (v *Vote) Insert() error {
	db := Get()

	if _, err := db.DB.NewInsert().Model(v).Exec(context.Background()); err != nil {
		return fmt.Errorf("insert Vote model: %w", err)
	}

	return nil
}

func GetAllVotesForElection(electionID int) ([]*Vote, error) {
	db := Get()
	var res []*Vote
	if err := db.DB.NewSelect().Model(&res).Where(`"election_id" = ?`, electionID).Scan(context.Background(), &res); err != nil {
		return nil, fmt.Errorf("get all Votes for election %d: %w", electionID, err)
	}
	return res, nil
}

func HasUserVotedInElection(userID string, electionID int) (bool, error) {
	db := Get()
	count, err := db.DB.NewSelect().Model((*Vote)(nil)).Where("user_id = ? AND election_id = ?", userID, electionID).Count(context.Background())
	if err != nil {
		return false, fmt.Errorf("check if user %s has voted in %d", userID, electionID)
	}
	return count == 1, nil
}

func DeleteAllVotesForElection(electionID int) error {
	db := Get()
	if _, err := db.DB.NewDelete().Model((*Vote)(nil)).Where("election_id = ?", electionID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete all Votes for election %d: %w", electionID, err)
	}
	return nil
}
