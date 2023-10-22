package httpcore

import (
	"fmt"
	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/gofiber/fiber/v2"
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
