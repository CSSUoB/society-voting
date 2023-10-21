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

func (b *BallotEntry) Insert(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewInsert().Model(b).Exec(context.Background()); err != nil {
		return fmt.Errorf("insert BallotEntry model: %w", err)
	}

	return nil
}

func (b *BallotEntry) Delete(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewDelete().Model(b).Where("id = ?", b.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete BallotEntry model: %w", err)
	}

	return nil
}

func GetAllBallotEntriesForElection(electionID int, x ...bun.IDB) ([]*BallotEntry, error) {
	db := fromVariadic(x)
	var res []*BallotEntry
	if err := db.NewSelect().Model(&res).Where(`"election_id" = ?`, electionID).Scan(context.Background(), &res); err != nil {
		return nil, fmt.Errorf("get all BallotEntrys for election %d: %w", electionID, err)
	}
	return res, nil
}

func DeleteBallotForElection(electionID int, x ...bun.IDB) error {
	db := fromVariadic(x)
	if _, err := db.NewDelete().Model((*BallotEntry)(nil)).Where("election_id = ?", electionID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete BallotEntrys for election: %w", err)
	}
	return nil
}

func CreateBallot(electionID int, names []string, x ...bun.IDB) ([]*BallotEntry, error) {
	var res []*BallotEntry

	b := &BallotEntry{
		ElectionID: electionID,
		Name:       "Re-open nominations (RON)",
		IsRON:      true,
	}
	res = append(res, b)
	if err := b.Insert(x...); err != nil {
		return nil, fmt.Errorf("CreateBallot: %w", err)
	}

	for _, name := range names {
		b := &BallotEntry{
			ElectionID: electionID,
			Name:       name,
		}
		res = append(res, b)
		if err := b.Insert(x...); err != nil {
			return nil, fmt.Errorf("CreateBallot: %w", err)
		}
	}

	return res, nil
}
