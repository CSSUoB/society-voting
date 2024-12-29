package database

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

const ElectionPollTypeId = 1
const ReferendumPollTypeId = 2

type PollType struct {
	bun.BaseModel `json:"-"`

	ID   int    `bun:",pk,autoincrement" json:"id"`
	Name string `bun:",notnull" json:"name"`
}

func (e *PollType) Insert(x ...bun.IDB) error {
	db := fromVariadic(x)

	if err := db.NewInsert().Model(e).Returning("id").Scan(context.Background(), &e.ID); err != nil {
		return fmt.Errorf("insert PollType model: %w", err)
	}

	return nil
}
