package database

import (
	"context"
	"database/sql"
	"errors"

	// "database/sql"
	// "errors"
	"fmt"

	"github.com/uptrace/bun"
)

type Referendum struct {
	bun.BaseModel `json:"-"`

	ID          int    `bun:",pk" json:"id"`
	Title       string `bun:",notnull" json:"title"`
	Question    string `bun:",notnull" json:"question"`
	Description string `bun:",notnull" json:"description"`

	Poll *Poll `bun:"rel:belongs-to,join:id=id" json:"-"`
}

func (e *Referendum) Insert(x ...bun.IDB) error {
	db := fromVariadic(x)

	if err := db.NewInsert().Model(e).Returning("id").Scan(context.Background(), &e.ID); err != nil {
		return fmt.Errorf("insert Referendum model: %w", err)
	}

	return nil
}

func (e *Referendum) Update(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewUpdate().Model(e).Where("id = ?", e.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("update Referendum model: %w", err)
	}

	return nil
}

func (e *Referendum) Delete(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewDelete().Model(e).Where("id = ?", e.ID).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete Referendum model: %w", err)
	}

	return nil
}

func GetReferendum(id int, x ...bun.IDB) (*Referendum, error) {
	db := fromVariadic(x)
	res := new(Referendum)
	if err := db.NewSelect().Model(res).Where(`"referendum"."id" = ?`, id).Relation("Poll").Scan(context.Background(), res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get Referendum model by ID: %w", err)
	}
	return res, nil
}
