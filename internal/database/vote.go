package database

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type Vote struct {
	bun.BaseModel `json:"-"`

	ID      int    `bun:",pk,autoincrement" json:"id"`
	PollID  int    `json:"pollID"`
	UserID  string `json:"userID"`
	Choices []int  `json:"choices"`

	Poll *Poll `bun:"rel:belongs-to,join:poll_id=id" json:"poll"`
}

func (v *Vote) Insert(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewInsert().Model(v).Exec(context.Background()); err != nil {
		return fmt.Errorf("insert Vote model: %w", err)
	}

	return nil
}

func GetAllVotesForPoll(pollID int, x ...bun.IDB) ([]*Vote, error) {
	db := fromVariadic(x)
	var res []*Vote
	if err := db.NewSelect().Model(&res).Where(`"poll_id" = ?`, pollID).Scan(context.Background(), &res); err != nil {
		return nil, fmt.Errorf("get all Votes for poll %d: %w", pollID, err)
	}
	return res, nil
}

func HasUserVotedInPoll(userID string, pollID int, x ...bun.IDB) (bool, error) {
	db := fromVariadic(x)
	count, err := db.NewSelect().Model((*Vote)(nil)).Where("user_id = ? AND poll_id = ?", userID, pollID).Count(context.Background())
	if err != nil {
		return false, fmt.Errorf("check if user %s has voted in %d", userID, pollID)
	}
	return count == 1, nil
}

func DeleteAllVotesForPoll(pollID int, x ...bun.IDB) error {
	db := fromVariadic(x)
	if _, err := db.NewDelete().Model((*Vote)(nil)).Where("poll_id = ?", pollID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete all Votes for poll %d: %w", pollID, err)
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

func CountVotesForElection(pollID int, x ...bun.IDB) (int, error) {
	db := fromVariadic(x)
	n, err := db.NewSelect().Model((*Vote)(nil)).Where("poll_id = ?", pollID).Count(context.Background())
	if err != nil {
		return 0, fmt.Errorf("count Votes for poll: %w", err)
	}
	return n, nil
}
