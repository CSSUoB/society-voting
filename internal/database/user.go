package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `json:"-"`

	StudentID    string `bun:"id,pk" json:"studentID"`
	Name         string `json:"name"`
	PasswordHash []byte `json:"-"`
}

func (u *User) Insert() error {
	db := Get()

	if _, err := db.DB.NewInsert().Model(u).Exec(context.Background()); err != nil {
		return fmt.Errorf("insert User model: %w", err)
	}

	return nil
}

func (u *User) Update() error {
	db := Get()

	if _, err := db.DB.NewUpdate().Model(u).Where("id = ?", u.StudentID).Exec(context.Background()); err != nil {
		return fmt.Errorf("update user model: %w", err)
	}

	return nil
}

func GetUser(id string) (*User, error) {
	db := Get()
	res := new(User)
	if err := db.DB.NewSelect().Model(res).Where("id = ?", id).Scan(context.Background(), res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get User model by ID: %w", err)
	}
	return res, nil
}
