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

func (v *Vote) Insert(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewInsert().Model(v).Exec(context.Background()); err != nil {
		return fmt.Errorf("insert Vote model: %w", err)
	}

	return nil
}

func GetAllVotesForElection(electionID int, x ...bun.IDB) ([]*Vote, error) {
	db := fromVariadic(x)
	var res []*Vote
	if err := db.NewSelect().Model(&res).Where(`"election_id" = ?`, electionID).Scan(context.Background(), &res); err != nil {
		return nil, fmt.Errorf("get all Votes for election %d: %w", electionID, err)
	}
	return res, nil
}

func HasUserVotedInElection(userID string, electionID int, x ...bun.IDB) (bool, error) {
	db := fromVariadic(x)
	count, err := db.NewSelect().Model((*Vote)(nil)).Where("user_id = ? AND election_id = ?", userID, electionID).Count(context.Background())
	if err != nil {
		return false, fmt.Errorf("check if user %s has voted in %d", userID, electionID)
	}
	return count == 1, nil
}

func DeleteAllVotesForElection(electionID int, x ...bun.IDB) error {
	db := fromVariadic(x)
	if _, err := db.NewDelete().Model((*Vote)(nil)).Where("election_id = ?", electionID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete all Votes for election %d: %w", electionID, err)
	}
	return nil
}

func DeleteAllVotesForUser(userID string, x ...bun.IDB) error {
	db := fromVariadic(x)
	if _, err := db.NewDelete().Model((*Vote)(nil)).Where("user_id = ?", userID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete user Votes: %w", err)
	}
	return nil
}

func CountVotesForElection(electionID int, x ...bun.IDB) (int, error) {
	db := fromVariadic(x)
	n, err := db.NewSelect().Model((*Vote)(nil)).Where("election_id = ?", electionID).Count(context.Background())
	if err != nil {
		return 0, fmt.Errorf("count Votes for election: %w", err)
	}
	return n, nil
}
