package httpcore

import (
	"errors"
	"fmt"
	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/gofiber/fiber/v2"
)

func (endpoints) apiMe(ctx *fiber.Ctx) error {
	userID, isAuthed := getSessionAuth(ctx, authRegularUser)
	if !isAuthed {
		return fiber.ErrUnauthorized
	}

	user, err := database.GetUser(userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			// User has been deleted
			ctx.Cookie(newSessionTokenDeletionCookie())
			return fiber.ErrUnauthorized
		}
		return fmt.Errorf("apiVote get user: %w", err)
	}

	return ctx.JSON(user)
}

func (endpoints) apiSetOwnName(ctx *fiber.Ctx) error {
	userID, isAuthed := getSessionAuth(ctx, authRegularUser)
	if !isAuthed {
		return fiber.ErrUnauthorized
	}

	var request = struct {
		Name string `json:"name" validate:"required"`
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
