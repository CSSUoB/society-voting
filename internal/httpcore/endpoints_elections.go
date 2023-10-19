package httpcore

import (
	"fmt"
	"git.tdpain.net/codemicro/society-voting/internal/database"
	"git.tdpain.net/codemicro/society-voting/internal/events"
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

	for _, election := range elections {
		if err := election.PopulateCandidates(); err != nil {
			return fmt.Errorf("apiListElections: %w", err)
		}
	}

	return ctx.JSON(elections)
}

func (endpoints) apiElectionsSSE(ctx *fiber.Ctx) error {
	_, isAuthed, err := getSessionAuth(ctx)
	if err != nil {
		return err
	}
	if !isAuthed {
		return fiber.ErrUnauthorized
	}

	ctx.Set("Content-Type", "text/event-stream")
	fr := ctx.Response()
	fr.SetBodyStreamWriter(
		events.AsStreamWriter(events.NewReceiver(events.TopicElectionStarted)),
	)

	return nil
}
