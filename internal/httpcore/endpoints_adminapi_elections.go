package httpcore

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/CSSUoB/society-voting/internal/events"
	"github.com/CSSUoB/society-voting/internal/instantRunoff"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"time"
)

func (endpoints) apiAdminCreateElection(ctx *fiber.Ctx) error {
	if _, ok := getSessionAuth(ctx, authAdminUser); !ok {
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
	if _, ok := getSessionAuth(ctx, authAdminUser); !ok {
		return fiber.ErrUnauthorized
	}

	var request = struct {
		ElectionID int `json:"id" validate:"required"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiAdminDeleteElection start tx: %w", err)
	}
	defer tx.Rollback()

	if err := database.DeleteCandidatesForElection(request.ElectionID, tx); err != nil {
		return fmt.Errorf("apiAdminDeleteElection delete all candidates: %w", err)
	}

	if err := database.DeleteElectionByID(request.ElectionID, tx); err != nil {
		return fmt.Errorf("apiAdminDeleteElection delete election: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiAdminDeleteElection commit tx: %w", err)
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}

func (endpoints) apiAdminStartElection(ctx *fiber.Ctx) error {
	if _, ok := getSessionAuth(ctx, authAdminUser); !ok {
		return fiber.ErrUnauthorized
	}

	var request = struct {
		ElectionID int      `json:"id" validate:"required"`
		ExtraNames []string `json:"extraNames" validate:"dive,required"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiAdminStartElection start tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := database.GetActiveElection(tx); err == nil {
		return &fiber.Error{
			Code:    fiber.StatusConflict,
			Message: "There is already an election in progress.",
		}
	}

	election, err := database.GetElection(request.ElectionID, tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return &fiber.Error{
				Code:    fiber.StatusNotFound,
				Message: "Election with that ID not found",
			}
		}
		return fmt.Errorf("apiAdminStartElection get election by ID: %w", err)
	}

	election.IsActive = true

	if err := election.Update(tx); err != nil {
		return fmt.Errorf("apiAdminStartElection update election: %w", err)
	}

	var ballot []*database.BallotEntry
	{
		candidates, err := database.GetUsersStandingForElection(request.ElectionID, tx)
		if err != nil {
			return fmt.Errorf("apiAdminStartElection get users standing for election: %w", err)
		}

		var names []string
		for _, candidate := range candidates {
			names = append(names, candidate.Name)
		}

		names = append(names, request.ExtraNames...)

		ballot, err = database.CreateBallot(request.ElectionID, names, tx)
		if err != nil {
			return fmt.Errorf("apiAdminStartElection create ballot: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiAdminStartElection commit tx: %w", err)
	}

	events.SendEvent(events.TopicElectionStarted, election)

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
	if _, ok := getSessionAuth(ctx, authAdminUser); !ok {
		return fiber.ErrUnauthorized
	}

	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiAdminStopElection start tx: %w", err)
	}
	defer tx.Rollback()

	election, err := database.GetActiveElection(tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return &fiber.Error{
				Code:    fiber.StatusConflict,
				Message: "There is no active election to stop.",
			}
		}
		return fmt.Errorf("apiVote get active election: %wz", err)
	}

	// Count votes

	votes, err := database.GetAllVotesForElection(election.ID, tx)
	if err != nil {
		return fmt.Errorf("apiAdminStopElection get all votes: %w", err)
	}

	var irVotes []*instantRunoff.Vote
	for _, v := range votes {
		irVotes = append(irVotes, &instantRunoff.Vote{
			RankedChoices: v.Choices,
		})
	}

	ballotEntries, err := database.GetAllBallotEntriesForElection(election.ID, tx)
	if err != nil {
		return fmt.Errorf("apiAdminStopElection get all ballot entries: %w", err)
	}

	ballotNames := make(map[int]string)
	for _, be := range ballotEntries {
		ballotNames[be.ID] = be.Name
	}

	resultText := instantRunoff.Run(irVotes, ballotNames)
	resultText = fmt.Sprintf("ELECTION OF %s BY %d MEMBERS ON %s UTC\n=================================================================\n\n", election.RoleName, len(votes), time.Now().UTC().Format(time.DateTime)) + resultText

	// Delete votes, ballot and election
	if err := database.DeleteAllVotesForElection(election.ID, tx); err != nil {
		return fmt.Errorf("apiAdminStopElection delete all votes: %w", err)
	}

	if err := database.DeleteBallotForElection(election.ID, tx); err != nil {
		return fmt.Errorf("apiAdminStopElection delete ballot: %w", err)
	}

	if err := database.DeleteCandidatesForElection(election.ID, tx); err != nil {
		return fmt.Errorf("apiAdminStopElection delete candidates: %w", err)
	}

	if err := election.Delete(tx); err != nil {
		return fmt.Errorf("apiAdminStopElection delete election: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiAdminStopElection commit tx: %w", err)
	}

	events.SendEvent(events.TopicElectionEnded, &events.ElectionEndedData{
		Election: election,
		Result:   resultText,
	})

	// Return election results

	var response = struct {
		Result string `json:"result"`
	}{
		Result: resultText,
	}

	return ctx.JSON(response)
}

func (endpoints) apiAdminRunningElectionSSE(ctx *fiber.Ctx) error {
	if _, ok := getSessionAuth(ctx, authAdminUser); !ok {
		return fiber.ErrUnauthorized
	}

	activeElection, err := database.GetActiveElection()
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return &fiber.Error{
				Code:    fiber.StatusConflict,
				Message: "There is no election in progress.",
			}
		}
		return fmt.Errorf("apiAdminRunningElectionSSE get active election: %w", err)
	}

	receiverID, receiver := events.NewReceiver(events.TopicVoteReceived)
	electionEndReceiverID, electionEndReceiver := events.NewReceiver(events.TopicElectionEnded)

	ctx.Set("Content-Type", "text/event-stream")
	fr := ctx.Response()
	fr.SetBodyStreamWriter(
		func(w *bufio.Writer) {
			count, err := database.CountVotesForElection(activeElection.ID)
			if err != nil {
				slog.Error("SSE error", "error", fmt.Errorf("apiAdminRunningElectionSSE count votes: %w", err))
				goto cleanup
			}

		loop:
			for {
				select {
				case msg := <-receiver:
					count += 1
					msg.Data = count
					sseData, err := msg.ToSSE()
					if err != nil {
						slog.Error("SSE error", "error", fmt.Errorf("failed to generate SSE event from message: %w", err))
						continue loop
					}
					_, _ = w.Write(sseData)
					if err := w.Flush(); err != nil {
						// Client disconnected
						break loop
					}
				case <-electionEndReceiver:
					var err error
					for err == nil {
						// if we force the client to disconnect, it'll just try to reconnect.
						// wait for it to disconnect, which will cause an error on Flush
						err = w.Flush()
						time.Sleep(time.Second * 5)
					}
					break loop
				}
			}

		cleanup:
			events.CloseReceiver(receiverID)
			events.CloseReceiver(electionEndReceiverID)
		},
	)

	return nil
}
