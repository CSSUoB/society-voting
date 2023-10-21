package httpcore

import (
	"errors"
	"fmt"
	"git.tdpain.net/codemicro/society-voting/internal/database"
	"git.tdpain.net/codemicro/society-voting/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/mattn/go-sqlite3"
)

func (endpoints) apiMe(ctx *fiber.Ctx) error {
	user, isAuthed, err := getSessionAuth(ctx)
	if err != nil {
		return err
	}
	if !isAuthed {
		return fiber.ErrUnauthorized
	}

	return ctx.JSON(user)
}

func (endpoints) apiSetOwnName(ctx *fiber.Ctx) error {
	user, isAuthed, err := getSessionAuth(ctx)
	if err != nil {
		return err
	}
	if !isAuthed {
		return fiber.ErrUnauthorized
	}

	var request = struct {
		Name string `json:"name" validate:"required"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	user.Name = request.Name

	if err := user.Update(); err != nil {
		return fmt.Errorf("apiSetOwnName update user: %w", err)
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}

func (endpoints) apiStandForElection(ctx *fiber.Ctx) error {
	user, isAuthed, err := getSessionAuth(ctx)
	if err != nil {
		return err
	}
	if !isAuthed {
		return fiber.ErrUnauthorized
	}

	var request = struct {
		ElectionID int `json:"id" validate:"ne=0"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiStandForElection start tx: %w", err)
	}
	defer util.Warn(tx.Rollback())

	election, err := database.GetElection(request.ElectionID, tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return fiber.ErrNotFound
		}
		return fmt.Errorf("apiStandForElection get election with id %d: %w", request.ElectionID, err)
	}

	candidate := &database.Candidate{
		UserID:     user.StudentID,
		ElectionID: election.ID,
	}

	if err := candidate.Insert(tx); err != nil {
		if e2 := errors.Unwrap(err); e2 != nil {
			var e sqlite3.Error
			if errors.As(e2, &e) {
				if e.Code == sqlite3.ErrConstraint {
					return &fiber.Error{
						Code:    fiber.StatusConflict,
						Message: "You're already standing in this election.",
					}
				}
			}
		}
		return fmt.Errorf("apiStandForElection create candidacy: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiStandForElection commit tx: %w", err)
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}

func (endpoints) apiWithdrawFromElection(ctx *fiber.Ctx) error {
	user, isAuthed, err := getSessionAuth(ctx)
	if err != nil {
		return err
	}
	if !isAuthed {
		return fiber.ErrUnauthorized
	}

	var request = struct {
		ElectionID int `json:"id" validate:"ne=0"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	candidate := &database.Candidate{
		UserID:     user.StudentID,
		ElectionID: request.ElectionID,
	}

	if err := candidate.Delete(); err != nil {
		return fmt.Errorf("apiWithdrawFromElection delete candidacy: %w", err)
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}
