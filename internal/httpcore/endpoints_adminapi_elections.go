package httpcore

import (
	"fmt"
	"git.tdpain.net/codemicro/society-voting/internal/database"
	"github.com/gofiber/fiber/v2"
)

func (endpoints) apiAdminCreateElection(ctx *fiber.Ctx) error {
	if !isAdminSession(ctx) {
		return fiber.ErrUnauthorized
	}

	var request = struct {
		RoleName    string `validate:"required"`
		Description string
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	election := &database.Election{
		RoleName:    request.RoleName,
		Description: request.Description,
	}

	if err := election.Insert(); err != nil {
		return fmt.Errorf("apiAdminCreateElection create election: %w", err)
	}

	return ctx.JSON(election)
}

func (endpoints) apiAdminDeleteElection(ctx *fiber.Ctx) error {
	if !isAdminSession(ctx) {
		return fiber.ErrUnauthorized
	}

	var request = struct {
		ElectionID int `json:"id" validate:"required"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	if err := database.DeleteCandidatesForElection(request.ElectionID); err != nil {
		return fmt.Errorf("apiAdminDeleteElection delete all candidates: %w", err)
	}

	if err := database.DeleteElectionByID(request.ElectionID); err != nil {
		return fmt.Errorf("apiAdminDeleteElection delete election: %w", err)
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}
