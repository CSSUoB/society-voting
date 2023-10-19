package httpcore

import (
	"fmt"
	"git.tdpain.net/codemicro/society-voting/internal/database"
	"github.com/gofiber/fiber/v2"
)

func (endpoints) apiListElections(ctx *fiber.Ctx) error {
	_, isAuthed, err := getSessionAuth(ctx)
	if err != nil {
		return err
	}
	if !isAuthed {
		return fiber.ErrUnauthorized
	}

	elections, err := database.GetAllElections()
	if err != nil {
		return fmt.Errorf("apiListElections get all elections: %w", err)
	}

	return ctx.JSON(elections)
}
