package httpcore

import (
	"fmt"

	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/gofiber/fiber/v2"
)

func (endpoints) apiMe(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*database.User)
	return ctx.JSON(user)
}

func (endpoints) apiSetOwnName(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*database.User)

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

	user.Name = request.Name

	if err := user.Update(tx); err != nil {
		return fmt.Errorf("apiSetOwnName update user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiSetOwnName commit tx: %w", err)
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}
