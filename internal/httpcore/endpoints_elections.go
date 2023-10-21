package httpcore

import (
	"errors"
	"fmt"
	"git.tdpain.net/codemicro/society-voting/internal/database"
	"git.tdpain.net/codemicro/society-voting/internal/events"
	"git.tdpain.net/codemicro/society-voting/internal/util"
	"github.com/gofiber/fiber/v2"
)

func (endpoints) apiListElections(ctx *fiber.Ctx) error {
	if _, ok := getSessionAuth(ctx, authAdminUser|authRegularUser); !ok {
		return fiber.ErrUnauthorized
	}

	elections, err := database.GetAllElections()
	if err != nil {
		return fmt.Errorf("apiListElections get all elections: %w", err)
	}

	var res []*database.ElectionWithCandidates

	for _, election := range elections {
		if ec, err := election.WithCandidates(); err != nil {
			return fmt.Errorf("apiListElections: %w", err)
		} else {
			res = append(res, ec)
		}
	}

	return ctx.JSON(res)
}

func (endpoints) apiElectionsSSE(ctx *fiber.Ctx) error {
	if _, ok := getSessionAuth(ctx, authAdminUser|authRegularUser); !ok {
		return fiber.ErrUnauthorized
	}

	ctx.Set("Content-Type", "text/event-stream")
	fr := ctx.Response()
	fr.SetBodyStreamWriter(
		events.AsStreamWriter(events.NewReceiver(events.TopicElectionStarted)),
	)

	return nil
}

func (endpoints) apiGetActiveElectionInformation(ctx *fiber.Ctx) error {
	if _, ok := getSessionAuth(ctx, authAdminUser|authRegularUser); !ok {
		return fiber.ErrUnauthorized
	}

	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiGetActiveElectionInformation start tx: %w", err)
	}
	defer util.Warn(tx.Rollback())

	election, err := database.GetActiveElection(tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return &fiber.Error{
				Code:    fiber.StatusConflict,
				Message: "There is no active election.",
			}
		}
		return fmt.Errorf("apiVote get active election: %wz", err)
	}

	ballot, err := database.GetAllBallotEntriesForElection(election.ID, tx)
	if err != nil {
		return fmt.Errorf("apiGetActiveElectionInformation get ballot: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiGetActiveElectionInformation commit tx: %w", err)
	}

	var response = struct {
		Election *database.Election      `json:"election"`
		Ballot   []*database.BallotEntry `json:"ballot"`
	}{
		Election: election,
		Ballot:   ballot,
	}

	return ctx.JSON(response)
}

func (endpoints) apiVote(ctx *fiber.Ctx) error {
	userID, isAuthed := getSessionAuth(ctx, authRegularUser)
	if !isAuthed {
		return fiber.ErrUnauthorized
	}

	var request = struct {
		Vote []int  `json:"vote" validate:"unique"`
		Code string `json:"code" validate:"required"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiVote start tx: %w", err)
	}
	defer util.Warn(tx.Rollback())

	user, err := database.GetUser(userID, tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			// User has been deleted
			ctx.Cookie(newSessionTokenDeletionCookie())
			return fiber.ErrUnauthorized
		}
		return fmt.Errorf("apiVote get user: %w", err)
	}

	election, err := database.GetActiveElection(tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return &fiber.Error{
				Code:    fiber.StatusConflict,
				Message: "There is no active election that you can vote in.",
			}
		}
		return fmt.Errorf("apiVote get active election: %wz", err)
	}

	hasVotedAlready, err := database.HasUserVotedInElection(user.StudentID, election.ID, tx)
	if err != nil {
		return fmt.Errorf("apiVote check if user %s has already voted: %w", user.StudentID, err)
	}

	if hasVotedAlready {
		return &fiber.Error{
			Code:    fiber.StatusConflict,
			Message: "You have already voted. Go away :)",
		}
	}

	if err := (&database.Vote{
		ElectionID: election.ID,
		UserID:     user.StudentID,
		Choices:    request.Vote,
	}).Insert(tx); err != nil {
		return fmt.Errorf("apiVote insert user vote: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiVote commit tx: %w", err)
	}

	events.SendEvent(events.TopicVoteReceived, nil)

	ctx.Status(fiber.StatusNoContent)
	return nil
}
