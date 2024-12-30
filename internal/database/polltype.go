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

func (p *PollType) Insert(x ...bun.IDB) error {
	db := fromVariadic(x)

	if err := db.NewInsert().Model(p).Returning("id").Scan(context.Background(), &p.ID); err != nil {
		return fmt.Errorf("insert PollType model: %w", err)
	}

	return nil
}
