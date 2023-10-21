package httpcore

import (
	"errors"
	"fmt"
	"git.tdpain.net/codemicro/society-voting/internal/database"
	"git.tdpain.net/codemicro/society-voting/internal/events"
	"git.tdpain.net/codemicro/society-voting/internal/instantRunoff"
	"github.com/gofiber/fiber/v2"
	"time"
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

func (endpoints) apiAdminStartElection(ctx *fiber.Ctx) error {
	if !isAdminSession(ctx) {
		return fiber.ErrUnauthorized
	}

	if _, err := database.GetActiveElection(); err == nil {
		return &fiber.Error{
			Code:    fiber.StatusConflict,
			Message: "There is already an election in progress.",
		}
	}

	var request = struct {
		ElectionID int      `json:"id" validate:"required"`
		ExtraNames []string `json:"extraNames" validate:"dive,required"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	election, err := database.GetElection(request.ElectionID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return fiber.ErrNotFound
		}
		return fmt.Errorf("apiAdminStartElection get election by ID: %w", err)
	}

	election.IsActive = true

	if err := election.Update(); err != nil {
		return fmt.Errorf("apiAdminStartElection update election: %w", err)
	}

	var ballot []*database.BallotEntry
	{
		candidates, err := database.GetUsersStandingForElection(request.ElectionID)
		if err != nil {
			return fmt.Errorf("apiAdminStartElection get users standing for election: %w", err)
		}

		var names []string
		for _, candidate := range candidates {
			names = append(names, candidate.Name)
		}

		names = append(names, request.ExtraNames...)

		ballot, err = database.CreateBallot(request.ElectionID, names)
		if err != nil {
			return fmt.Errorf("apiAdminStartElection create ballot: %w", err)
		}
	}

	events.SendEvent(events.TopicElectionStarted, election.ID)

	var response = struct {
		Election *database.Election      `json:"election"`
		Ballot   []*database.BallotEntry `json:"ballot"`
	}{
		Election: election,
		Ballot:   ballot,
	}

	return ctx.JSON(response)
}

func (endpoints) apiAdminStopElection(ctx *fiber.Ctx) error {
	if !isAdminSession(ctx) {
		return fiber.ErrUnauthorized
	}

	election, err := database.GetActiveElection()
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return &fiber.Error{
				Code:    fiber.StatusConflict,
				Message: "There is no active election that you can vote in.",
			}
		}
		return fmt.Errorf("apiVote get active election: %wz", err)
	}

	// Count votes

	votes, err := database.GetAllVotesForElection(election.ID)
	if err != nil {
		return fmt.Errorf("apiAdminStopElection get all votes: %w", err)
	}

	var irVotes []*instantRunoff.Vote
	for _, v := range votes {
		irVotes = append(irVotes, &instantRunoff.Vote{
			RankedChoices: v.Choices,
		})
	}

	ballotEntries, err := database.GetAllBallotEntriesForElection(election.ID)
	if err != nil {
		return fmt.Errorf("apiAdminStopElection get all ballot entries: %w", err)
	}

	ballotNames := make(map[int]string)
	for _, be := range ballotEntries {
		ballotNames[be.ID] = be.Name
	}

	resultText := instantRunoff.Run(irVotes, ballotNames)
	resultText = fmt.Sprintf("ELECTION OF %s BY %d MEMBERS ON %s UTC\n=================================================================\n", election.RoleName, len(votes), time.Now().UTC().Format(time.DateTime)) + resultText

	// Delete votes, ballot and election
	if err := database.DeleteAllVotesForElection(election.ID); err != nil {
		return fmt.Errorf("apiAdminStopElection delete all votes: %w", err)
	}

	if err := database.DeleteBallotForElection(election.ID); err != nil {
		return fmt.Errorf("apiAdminStopElection delete ballot: %w", err)
	}

	if err := election.Delete(); err != nil {
		return fmt.Errorf("apiAdminStopElection delete election: %w", err)
	}

	// Return election results

	var response = struct {
		Result string `json:"result"`
	}{
		Result: resultText,
	}

	return ctx.JSON(response)
}
