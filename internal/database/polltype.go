package database

import (
	"github.com/uptrace/bun"
)

const ElectionPollTypeId = 1
const ReferendumPollTypeId = 2

type PollType struct {
	bun.BaseModel `json:"-"`

	ID   int    `bun:",pk,autoincrement" json:"id"`
	Name string `bun:",notnull" json:"name"`
}
