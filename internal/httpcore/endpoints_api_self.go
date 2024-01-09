package httpcore

import (
	"errors"
	"fmt"
	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/gofiber/fiber/v2"
)

func (endpoints) apiMe(ctx *fiber.Ctx) error {
	userID, authStatus := getSessionAuth(ctx)
	if authStatus == authNotAuthed {
		return fiber.ErrUnauthorized
	}

	var user *database.User

	if userID == "admin" {
		user = &database.User{
			StudentID: "admin",
			Name:      "Administrator",
		}
	} else {
		u, err := database.GetUser(userID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				// User has been deleted
				ctx.Cookie(newSessionTokenDeletionCookie())
				return fiber.ErrUnauthorized
			}
			return fmt.Errorf("apiVote get user: %w", err)
		}
		user = u
	}

	return ctx.JSON(user)
}

func (endpoints) apiSetOwnName(ctx *fiber.Ctx) error {
	userID, authStatus := getSessionAuth(ctx)
	if authStatus == authNotAuthed {
		return fiber.ErrUnauthorized
	}

	var request = struct {
		Name string `json:"name" validate:"required,max=64"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiSetOwnName start tx: %w", err)
	}
	defer tx.Rollback()

	user, err := database.GetUser(userID, tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			// User has been deleted
			ctx.Cookie(newSessionTokenDeletionCookie())
			return fiber.ErrUnauthorized
		}
		return fmt.Errorf("apiSetOwnName get user: %w", err)
	}

	if user.IsRestricted {
		return &fiber.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can't do that because you're restricted - please speak to a member of committee.",
		}
	}

	user.Name = request.Name

	if err := user.Update(tx); err != nil {
		return fmt.Errorf("apiSetOwnName update user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiSetOwnName commit tx: %w")
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}
