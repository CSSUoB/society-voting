package database

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
)

type BallotEntry struct {
	bun.BaseModel `bun:"ballot_entry" json:"-"`

	ID         int    `bun:",pk,autoincrement" json:"id"`
	ElectionID int    `json:"electionID"`
	Name       string `json:"name"`
	IsRON      bool   `json:"isRON"`
}

func (b *BallotEntry) Insert() error {
	db := Get()

	if _, err := db.DB.NewInsert().Model(b).Exec(context.Background()); err != nil {
		return fmt.Errorf("insert BallotEntry model: %w", err)
	}

	return nil
}

func (b *BallotEntry) Delete() error {
	db := Get()

	if _, err := db.DB.NewDelete().Model(b).Where("id = ?", b.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete BallotEntry model: %w", err)
	}

	return nil
}

func GetAllBallotEntriesForElection(electionID int) ([]*BallotEntry, error) {
	db := Get()
	var res []*BallotEntry
	if err := db.DB.NewSelect().Model(&res).Where(`"election_id" = ?`, electionID).Scan(context.Background(), &res); err != nil {
		return nil, fmt.Errorf("get all BallotEntrys for election %d: %w", electionID, err)
	}
	return res, nil
}

func DeleteBallotForElection(electionID int) error {
	db := Get()
	if _, err := db.DB.NewDelete().Model((*BallotEntry)(nil)).Where("election_id = ?", electionID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete BallotEntrys for election: %w", err)
	}
	return nil
}

func CreateBallot(electionID int, names []string) ([]*BallotEntry, error) {
	var res []*BallotEntry

	b := &BallotEntry{
		ElectionID: electionID,
		Name:       "Re-open nominations (RON)",
		IsRON:      true,
	}
	res = append(res, b)
	if err := b.Insert(); err != nil {
		return nil, fmt.Errorf("CreateBallot: %w", err)
	}

	for _, name := range names {
		b := &BallotEntry{
			ElectionID: electionID,
			Name:       name,
		}
		res = append(res, b)
		if err := b.Insert(); err != nil {
			return nil, fmt.Errorf("CreateBallot: %w", err)
		}
	}

	return res, nil
}
