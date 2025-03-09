package httpcore

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/CSSUoB/society-voting/internal/events"
	"github.com/CSSUoB/society-voting/internal/instantRunoff"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

func createPoll(pollTypeId int, createEntity func(createdPollId int, tx bun.Tx) (any, error)) (any, error) {
	tx, err := database.GetTx()
	if err != nil {
		return nil, fmt.Errorf("start tx: %w", err)
	}
	defer tx.Rollback()

	poll := &database.Poll{
		PollTypeID: pollTypeId,
	}

	if err := poll.Insert(tx); err != nil {
		return nil, fmt.Errorf("create poll: %w", err)
	}

	entity, err := createEntity(poll.ID, tx)
	if err != nil {
		return nil, fmt.Errorf("create entity: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return entity, nil
}

func (endpoints) apiAdminCreateElection(ctx *fiber.Ctx) error {
	var request = struct {
		RoleName    string `validate:"required,max=64"`
		Description string `validate:"max=1024"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	entity, err := createPoll(database.ElectionPollTypeId, func(createdPollId int, tx bun.Tx) (any, error) {
		election := &database.Election{
			ID:          createdPollId,
			RoleName:    request.RoleName,
			Description: request.Description,
		}

		if err := election.Insert(tx); err != nil {
			return nil, err
		}

		return election, nil
	})
	if err != nil {
		return fmt.Errorf("apiAdminCreateElection %w", err)
	}

	return ctx.JSON(entity)
}

func (endpoints) apiAdminCreateReferendum(ctx *fiber.Ctx) error {
	var request = struct {
		Title       string `validate:"required,max=64"`
		Question    string `validate:"required,max=256"`
		Description string `validate:"max=1024"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	entity, err := createPoll(database.ReferendumPollTypeId, func(createdPollId int, tx bun.Tx) (any, error) {
		referendum := &database.Referendum{
			ID:          createdPollId,
			Title:       request.Title,
			Question:    request.Question,
			Description: request.Description,
		}

		if err := referendum.Insert(tx); err != nil {
			return nil, err
		}

		return referendum, nil
	})
	if err != nil {
		return fmt.Errorf("apiAdminCreateReferendum %w", err)
	}

	return ctx.JSON(entity)
}

func (endpoints) apiAdminDeletePoll(ctx *fiber.Ctx) error {
	var request = struct {
		PollID int `json:"id" validate:"required"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	if err := database.DeletePollByID(request.PollID); err != nil {
		return fmt.Errorf("apiAdminDeletePoll delete poll: %w", err)
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}

func startPoll(tx bun.Tx, getEntity func(int, bun.Tx) (database.Pollable, error), pollID int) (database.Pollable, error) {
	if _, err := database.GetActivePoll(tx); err == nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusConflict,
			Message: "There is already a poll in progress.",
		}
	}

	entity, err := getEntity(pollID, tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, &fiber.Error{
				Code:    fiber.StatusNotFound,
				Message: "Poll with that ID not found or is wrong type.",
			}
		}
		return nil, fmt.Errorf("get poll %d: %w", pollID, err)
	}

	poll := entity.GetPoll()
	if poll.IsConcluded {
		return nil, &fiber.Error{
			Code:    fiber.StatusConflict,
			Message: "This poll has already concluded.",
		}
	}

	poll.IsActive = true

	if err := poll.Update(tx); err != nil {
		return nil, fmt.Errorf("startPoll update poll: %w", err)
	}

	return entity, nil
}

func (endpoints) apiAdminStartElection(ctx *fiber.Ctx) error {
	var request = struct {
		ID         int      `json:"id" validate:"required"`
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

	pollable, err := startPoll(tx, func(id int, tx bun.Tx) (database.Pollable, error) {
		return database.GetElection(id, tx)
	}, request.ID)
	if err != nil {
		return fmt.Errorf("apiAdminStartElection %w", err)
	}

	candidates, err := database.GetUsersStandingForElection(request.ID, tx)
	if err != nil {
		return fmt.Errorf("apiAdminStartElection get users standing for election: %w", err)
	}

	var names []string
	for _, candidate := range candidates {
		names = append(names, candidate.Name)
	}

	names = append(names, request.ExtraNames...)

	_, err = database.CreateBallot(request.ID, names, tx)
	if err != nil {
		return fmt.Errorf("apiAdminStartElection create ballot: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiAdminStartElection commit tx: %w", err)
	}

	events.SendEvent(events.TopicPollStarted, pollable)

	ctx.Status(fiber.StatusNoContent)
	return nil
}

func (endpoints) apiAdminStartReferendum(ctx *fiber.Ctx) error {
	var request = struct {
		ID int `json:"id" validate:"required"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiAdminStartReferendum start tx: %w", err)
	}
	defer tx.Rollback()

	pollable, err := startPoll(tx, func(id int, tx bun.Tx) (database.Pollable, error) {
		return database.GetReferendum(id, tx)
	}, request.ID)
	if err != nil {
		return fmt.Errorf("apiAdminStartReferendum %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiAdminStartReferendum commit tx: %w", err)
	}

	events.SendEvent(events.TopicPollStarted, pollable)

	ctx.Status(fiber.StatusNoContent)
	return nil
}

func stopActivePoll(tx bun.Tx, getEntity func(int, bun.Tx) (database.Pollable, error)) (database.Pollable, []*database.Vote, error) {
	poll, err := database.GetActivePoll(tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, nil, &fiber.Error{
				Code:    fiber.StatusConflict,
				Message: "There is no active poll to stop.",
			}
		}
		return nil, nil, fmt.Errorf("get active poll: %w", err)
	}

	pollable, err := getEntity(poll.ID, tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, nil, &fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: "The active poll is not of the correct type.",
			}
		}
		return nil, nil, fmt.Errorf("check type: %w", err)
	}

	poll.IsActive = false
	poll.IsConcluded = true

	if err := poll.Update(tx); err != nil {
		return nil, nil, fmt.Errorf("update poll: %w", err)
	}

	votes, err := database.GetAllVotesForPoll(poll.ID, tx)
	if err != nil {
		return nil, nil, fmt.Errorf("get all votes: %w", err)
	}

	if err := database.DeleteAllVotesForPoll(poll.ID, tx); err != nil {
		return nil, nil, fmt.Errorf("delete all votes: %w", err)
	}

	return pollable, votes, nil
}

func (endpoints) apiAdminStopElection(ctx *fiber.Ctx) error {
	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiAdminStopElection start tx: %w", err)
	}
	defer tx.Rollback()

	pollable, votes, err := stopActivePoll(tx, func(id int, tx bun.Tx) (database.Pollable, error) {
		return database.GetElection(id, tx)
	})
	if err != nil {
		return fmt.Errorf("apiAdminStopElection stop active poll: %w", err)
	}

	// Count votes

	election := pollable.GetElection()

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

	instantRunoff, err := instantRunoff.Run(irVotes, ballotNames)
	if err != nil {
		return fmt.Errorf("apiAdminStopElection run instant runoff: %w", err)
	}

	resultText := fmt.Sprintf("ELECTION OF %s BY %d MEMBERS ON %s UTC\n=================================================================\n\n", election.RoleName, len(votes), time.Now().UTC().Format(time.DateTime)) + instantRunoff.ResultsAsString()

	pollOutcome, err := database.CreatePollOutcome(election.ID, len(votes), tx)
	if err != nil {
		return fmt.Errorf("apiAdminStopElection create poll outcome: %w", err)
	}

	electionOutcome := &database.ElectionOutcome{
		ID:     pollOutcome.ID,
		Rounds: instantRunoff.Rounds,
	}
	if err := electionOutcome.Insert(tx); err != nil {
		return fmt.Errorf("apiAdminStopElection create election outcome: %w", err)
	}

	var electionOutcomeResults []*database.ElectionOutcomeResult
	for _, tally := range instantRunoff.Tallies {
		electionOutcomeResults = append(electionOutcomeResults, &database.ElectionOutcomeResult{
			Name:              tally.Name,
			Round:             tally.Round,
			Votes:             tally.Count,
			IsRejected:        tally.Eliminated,
			IsElected:         tally.Winner,
			ElectionOutcomeID: electionOutcome.ID,
		})
	}

	if err := database.BulkInsertElectionOutcomeResult(electionOutcomeResults, tx); err != nil {
		return fmt.Errorf("apiAdminStopElection bulk insert election outcome results: %w", err)
	}

	// Delete ballot entries and candidates
	if err := database.DeleteBallotForElection(election.ID, tx); err != nil {
		return fmt.Errorf("apiAdminStopElection delete ballot: %w", err)
	}

	if err := database.DeleteCandidatesForElection(election.ID, tx); err != nil {
		return fmt.Errorf("apiAdminStopElection delete candidates: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiAdminStopElection commit tx: %w", err)
	}

	events.SendEvent(events.TopicPollEnded, &events.PollEndedData{
		Poll:   pollable.GetPoll(),
		Name:   pollable.GetFriendlyTitle(),
		Result: resultText,
	})

	ctx.Status(fiber.StatusNoContent)
	return nil
}

func (endpoints) apiAdminStopReferendum(ctx *fiber.Ctx) error {
	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiAdminStopReferendum start tx: %w", err)
	}
	defer tx.Rollback()

	pollable, votes, err := stopActivePoll(tx, func(id int, tx bun.Tx) (database.Pollable, error) {
		return database.GetReferendum(id, tx)
	})
	if err != nil {
		return fmt.Errorf("apiAdminStopElection stop active poll: %w", err)
	}

	referendum := pollable.GetReferendum()

	var votesAbstain, votesFor, votesAgainst int
	for _, vote := range votes {
		if len(vote.Choices) != 1 {
			// invalid ballot
			continue
		}
		for _, value := range vote.Choices {
			switch value {
			case 0:
				votesAbstain++
			case 1:
				votesFor++
			case 2:
				votesAgainst++
			}
		}
	}

	resultText := fmt.Sprintf("REFERENDUM ON %s BY %d MEMBERS ON %s UTC\n=================================================================\n\nQuestion: %s\n\nFor: %d\nAgainst: %d\nAbstain: %d", referendum.Title, len(votes), time.Now().UTC().Format(time.DateTime), referendum.Question, votesFor, votesAgainst, votesAbstain)

	pollOutcome, err := database.CreatePollOutcome(referendum.ID, len(votes), tx)
	if err != nil {
		return fmt.Errorf("apiAdminStopReferendum create poll outcome: %w", err)
	}

	referendumOutcome := &database.ReferendumOutcome{
		ID:           pollOutcome.ID,
		VotesAbstain: votesAbstain,
		VotesFor:     votesFor,
		VotesAgainst: votesAgainst,
	}
	if err := referendumOutcome.Insert(tx); err != nil {
		return fmt.Errorf("apiAdminStopElection create election outcome: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiAdminStopElection commit tx: %w", err)
	}

	events.SendEvent(events.TopicPollEnded, &events.PollEndedData{
		Poll:   pollable.GetPoll(),
		Name:   pollable.GetFriendlyTitle(),
		Result: resultText,
	})

	ctx.Status(fiber.StatusNoContent)
	return nil
}

func (endpoints) apiAdminPublishPollOutcome(ctx *fiber.Ctx) error {
	var request = struct {
		PollID    int   `json:"id" validate:"required"`
		Published *bool `json:"published" validate:"required"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiAdminPublishElectionOutcome start tx: %w", err)
	}
	defer tx.Rollback()

	pollOutcome, err := database.GetOutcomeForPoll(request.PollID, tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return &fiber.Error{
				Code:    fiber.StatusNotFound,
				Message: "There is no outcome for that election",
			}
		}
		return fmt.Errorf("apiAdminPublishElectionOutcome get election outcome: %w", err)
	}

	pollOutcome.IsPublished = *request.Published

	if err := pollOutcome.Update(tx); err != nil {
		return fmt.Errorf("apiAdminPublishElectonOutcome update election outcome: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiAdminPublishElectonOutcome commit tx: %w", err)
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}

func (endpoints) apiAdminRunningPollSSE(ctx *fiber.Ctx) error {
	activeElection, err := database.GetActivePoll()
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

	ctx.Set("Content-Type", "text/event-stream")
	ctx.Set("Connection", "keep-alive")

	doneNotify := ctx.Context().Done()
	fr := ctx.Response()
	fr.SetBodyStreamWriter(
		func(w *bufio.Writer) {
			ticker := time.NewTicker(time.Second * 10)
			defer ticker.Stop()

			count, err := database.CountVotesForElection(activeElection.ID)
			if err != nil {
				slog.Error("SSE error", "error", fmt.Errorf("apiAdminRunningElectionSSE count votes: %w", err))
				goto cleanup
			}

			// send initial vote count for late visitors
			{
				msg := events.Message{
					Topic: events.TopicVoteReceived,
					Data:  count,
				}
				sseData, err := msg.ToSSE()
				if err != nil {
					slog.Error("SSE error", "error", fmt.Errorf("failed to generate SSE event from message: %w", err))
					goto cleanup
				}

				if _, err := w.Write(sseData); err != nil {
					goto cleanup
				}
				if err := w.Flush(); err != nil {
					goto cleanup
				}
			}

		loop:
			for {
				select {
				case <-doneNotify:
					break loop
				case msg := <-receiver:
					count += 1
					msg.Data = count
					sseData, err := msg.ToSSE()
					if err != nil {
						slog.Error("SSE error", "error", fmt.Errorf("failed to generate SSE event from message: %w", err))
						continue loop
					}

					if _, err := w.Write(sseData); err != nil {
						break loop
					}
					if err := w.Flush(); err != nil {
						// Client disconnected
						break loop
					}
				case <-ticker.C:
					// heartbeat
					if _, err := w.Write([]byte(":\n\n")); err != nil {
						break loop
					}
					if err := w.Flush(); err != nil {
						break loop
					}
				}
			}

		cleanup:
			events.CloseReceiver(receiverID)
		},
	)

	return nil
}
