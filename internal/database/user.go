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
	PasswordHash string `json:"-"`
	IsRestricted bool   `json:"isRestricted"`
	IsAdmin      bool   `json:"isAdmin"`
}

func (u *User) Insert(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewInsert().Model(u).Exec(context.Background()); err != nil {
		return fmt.Errorf("insert User model: %w", err)
	}

	return nil
}

func (u *User) Update(x ...bun.IDB) error {
	db := fromVariadic(x)

	if _, err := db.NewUpdate().Model(u).Where("id = ?", u.StudentID).Exec(context.Background()); err != nil {
		return fmt.Errorf("update user model: %w", err)
	}

	return nil
}

func GetUser(id string, x ...bun.IDB) (*User, error) {
	db := fromVariadic(x)
	res := new(User)
	if err := db.NewSelect().Model(res).Where("id = ?", id).Scan(context.Background(), res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get User model by ID: %w", err)
	}
	return res, nil
}

func DeleteUser(id string, x ...bun.IDB) error {
	db := fromVariadic(x)
	if _, err := db.NewDelete().Model((*User)(nil)).Where("id = ?", id).Exec(context.Background()); err != nil {
		return fmt.Errorf("delete User: %w", err)
	}
	return nil
}

func CountUsers(x ...bun.IDB) (int, error) {
	db := fromVariadic(x)
	n, err := db.NewSelect().Model((*User)(nil)).Count(context.Background())
	if err != nil {
		return 0, fmt.Errorf("count Users: %w", err)
	}
	return n, nil
}

func GetAllUsers(x ...bun.IDB) ([]*User, error) {
	db := fromVariadic(x)
	var res []*User
	if err := db.NewSelect().Model(&res).Scan(context.Background(), &res); err != nil {
		return nil, fmt.Errorf("get all Users: %w", err)
	}
	return res, nil
}
