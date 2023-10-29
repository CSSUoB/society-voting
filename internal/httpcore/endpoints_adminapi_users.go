package httpcore

import (
	"errors"
	"fmt"
	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/gofiber/fiber/v2"
	"sort"
)

func (endpoints) apiAdminDeleteUser(ctx *fiber.Ctx) error {
	if _, isAuthed := getSessionAuth(ctx, authAdminUser); !isAuthed {
		return fiber.ErrUnauthorized
	}

	var request = struct {
		UserID string `json:"userID" validate:"ne=0"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiAdminDeleteUser start tx: %w", err)
	}
	defer tx.Rollback()

	if err := database.DeleteUser(request.UserID, tx); err != nil {
		return fmt.Errorf("apiAdminDeleteUser delete user: %w", err)
	}

	if err := database.DeleteAllCandidatesForUser(request.UserID, tx); err != nil {
		return fmt.Errorf("apiAdminDeleteUser delete user candidates: %w", err)
	}

	if err := database.DeleteAllVotesForUser(request.UserID, tx); err != nil {
		return fmt.Errorf("apiAdminDeleteUser delete user votes: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiAdminDeleteUser commit tx: %w", err)
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}

func (endpoints) apiAdminListUsers(ctx *fiber.Ctx) error {
	if _, isAuthed := getSessionAuth(ctx, authAdminUser); !isAuthed {
		return fiber.ErrUnauthorized
	}

	users, err := database.GetAllUsers()
	if err != nil {
		return fmt.Errorf("apiAdminListUsers get all users: %w", err)
	}

	// place restricted users at the top of the list
	sort.Slice(users, func(i, j int) bool {
		if users[i].IsRestricted != users[j].IsRestricted {
			return users[i].IsRestricted
		}
		return users[i].Name < users[j].Name
	})

	return ctx.JSON(users)
}

func (endpoints) apiAdminToggleRestrictUser(ctx *fiber.Ctx) error {
	if _, isAuthed := getSessionAuth(ctx, authAdminUser); !isAuthed {
		return fiber.ErrUnauthorized
	}

	var request = struct {
		UserID string `json:"userID" validate:"ne=0"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiAdminToggleRestrictUser start tx: %w", err)
	}
	defer tx.Rollback()

	user, err := database.GetUser(request.UserID, tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return &fiber.Error{
				Code:    fiber.StatusNotFound,
				Message: "User with that ID not found",
			}
		}
		return fmt.Errorf("apiAdminToggleRestrictUser get user: %w", err)
	}

	user.IsRestricted = !user.IsRestricted

	if err := user.Update(tx); err != nil {
		return fmt.Errorf("apiAdminToggleRestrictUser update user: %w", err)
	}

	if user.IsRestricted {
		if err := database.DeleteAllCandidatesForUser(user.StudentID, tx); err != nil {
			return fmt.Errorf("apiAdminToggleRestrictUser delete user candidacies: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiAdminToggleRestrictUser commit tx: %w", err)
	}

	return ctx.JSON(map[string]bool{"isRestricted": user.IsRestricted})
}
