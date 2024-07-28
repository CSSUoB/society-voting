package httpcore

import (
	"bufio"
	"crypto/subtle"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/CSSUoB/society-voting/internal/events"
	"github.com/gofiber/fiber/v2"
	"github.com/mattn/go-sqlite3"
)

func (endpoints) apiListElections(ctx *fiber.Ctx) error {
	userID, authStatus := getSessionAuth(ctx)
	if authStatus == authNotAuthed {
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
			for _, cand := range ec.Candidates {
				cand.IsMe = cand.ID == userID
			}
			res = append(res, ec)
		}
	}

	return ctx.JSON(res)
}

func (endpoints) apiElectionsSSE(ctx *fiber.Ctx) error {
	if _, status := getSessionAuth(ctx); status == authNotAuthed {
		return fiber.ErrUnauthorized
	}

	id, receiver := events.NewReceiver(events.TopicElectionStarted, events.TopicElectionEnded)

	ctx.Set("Content-Type", "text/event-stream")
	fr := ctx.Response()
	fr.SetBodyStreamWriter(
		func(w *bufio.Writer) {
			ticker := time.NewTicker(time.Second * 10)
			for {
				select {
				case msg := <-receiver:
					if msg.Topic == events.TopicElectionEnded {
						// we're going to be modifying this msg so let's create a copy and work with that
						{
							// TODO: Refactor away this copying mess
							x := *msg
							y := *(msg.Data.(*events.ElectionEndedData))
							x.Data = &y
							msg = &x
						}
						msg.Data.(*events.ElectionEndedData).Result = ""
					}
					sseData, err := msg.ToSSE()
					if err != nil {
						slog.Error("SSE error", "error", fmt.Errorf("failed to generate SSE event from message: %w", err))
						break
					}
					_, _ = w.Write(sseData)
				case <-ticker.C:
				}

				if err := w.Flush(); err != nil {
					// Client disconnected
					break
				}
			}
			events.CloseReceiver(id)
		},
	)

	return nil
}

func (endpoints) apiGetActiveElectionInformation(ctx *fiber.Ctx) error {
	userID, authStatus := getSessionAuth(ctx)
	if authStatus == authNotAuthed {
		return fiber.ErrUnauthorized
	}

	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiGetActiveElectionInformation start tx: %w", err)
	}
	defer tx.Rollback()

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

	numUsers, err := database.CountUsers(tx)
	if err != nil {
		return fmt.Errorf("apiGetActiveElectionInformation count users: %w", err)
	}

	hasVoted, err := database.HasUserVotedInElection(userID, election.ID, tx)
	if err != nil {
		return fmt.Errorf("apiGetActiveElectionInformation check if user has voted: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiGetActiveElectionInformation commit tx: %w", err)
	}

	// randomise ballot order
	sort.Slice(ballot, func(_, _ int) bool {
		return rand.Intn(2) == 0
	})

	var response = struct {
		Election *database.Election      `json:"election"`
		Ballot   []*database.BallotEntry `json:"ballot"`
		NumUsers int                     `json:"numEligibleVoters"`
		HasVoted bool                    `json:"hasVoted"`
	}{
		Election: election,
		Ballot:   ballot,
		NumUsers: numUsers,
		HasVoted: hasVoted,
	}

	return ctx.JSON(response)
}

func (endpoints) apiVote(ctx *fiber.Ctx) error {
	userID, authStatus := getSessionAuth(ctx)
	if authStatus == authNotAuthed {
		return fiber.ErrUnauthorized
	}

	var request = struct {
		Vote []int  `json:"vote" validate:"unique"`
		Code string `json:"code" validate:"required"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	if subtle.ConstantTimeCompare([]byte(strings.ToUpper(request.Code)), []byte(voteCode)) == 0 {
		return &fiber.Error{
			Code:    fiber.StatusForbidden,
			Message: "Incorrect vote code!",
		}
	}

	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiVote start tx: %w", err)
	}
	defer tx.Rollback()

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

	ballotOptions, err := database.GetAllBallotEntriesForElection(election.ID, tx)
	if err != nil {
		return fmt.Errorf("apiVote get all ballot entries: %w", err)
	}

	for _, id := range request.Vote {
		var found bool
		for _, b := range ballotOptions {
			if b.ID == id {
				found = true
				break
			}
		}
		if !found {
			return &fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: fmt.Sprintf("%d is not a valid ballot option.", id),
			}
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

func (endpoints) apiStandForElection(ctx *fiber.Ctx) error {
	userID, authStatus := getSessionAuth(ctx)
	if authStatus == authNotAuthed {
		return fiber.ErrUnauthorized
	}

	var request = struct {
		ElectionID int `json:"id" validate:"ne=0"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiStandForElection start tx: %w", err)
	}
	defer tx.Rollback()

	user, err := database.GetUser(userID, tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			// User has been deleted
			ctx.Cookie(newSessionTokenDeletionCookie())
			return fiber.ErrUnauthorized
		}
		return fmt.Errorf("apiStandForElection get user: %w", err)
	}

	if user.IsRestricted {
		return &fiber.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can't do that because you're restricted - please speak to a member of committee.",
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
		return fmt.Errorf("apiStandForElection get election with id %d: %w", request.ElectionID, err)
	}

	if election.IsConcluded {
		return &fiber.Error{
			Code:    fiber.StatusConflict,
			Message: "This election has already concluded",
		}
	}

	candidate := &database.Candidate{
		UserID:     user.StudentID,
		ElectionID: election.ID,
	}

	if err := candidate.Insert(tx); err != nil {
		if e2 := errors.Unwrap(err); e2 != nil {
			var e sqlite3.Error
			if errors.As(e2, &e) {
				if e.Code == sqlite3.ErrConstraint {
					return &fiber.Error{
						Code:    fiber.StatusConflict,
						Message: "You're already standing in this election.",
					}
				}
			}
		}
		return fmt.Errorf("apiStandForElection create candidacy: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiStandForElection commit tx: %w", err)
	}

	events.SendEvent(events.TopicUserElectionStand, &events.UserElectionStandData{
		User:     user,
		Election: election,
	})

	ctx.Status(fiber.StatusNoContent)
	return nil
}

func (endpoints) apiWithdrawFromElection(ctx *fiber.Ctx) error {
	var (
		electionID int
		userID     string
	)

	actingUserID, authStatus := getSessionAuth(ctx)
	if authStatus == authNotAuthed {
		return fiber.ErrUnauthorized
	}

	if authStatus&authAdminUser != 0 {
		var request = struct {
			ElectionID int    `json:"id" validate:"ne=0"`
			UserID     string `json:"userID"`
		}{}

		if err := parseAndValidateRequestBody(ctx, &request); err != nil {
			return err
		}

		if request.UserID == "" {
			request.UserID = actingUserID
		}

		userID = request.UserID
		electionID = request.ElectionID
	} else {
		var request = struct {
			ElectionID int `json:"id" validate:"ne=0"`
		}{}

		if err := parseAndValidateRequestBody(ctx, &request); err != nil {
			return err
		}

		userID = actingUserID
		electionID = request.ElectionID
	}

	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiWithdrawFromElection start tx: %w", err)
	}
	defer tx.Rollback()

	user, err := database.GetUser(userID, tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			if authStatus&authAdminUser != 0 {
				fmt.Println(userID)
				return &fiber.Error{
					Code:    fiber.StatusNotFound,
					Message: "User with that ID not found.",
				}
			}
			// User has been deleted
			ctx.Cookie(newSessionTokenDeletionCookie())
			return fiber.ErrUnauthorized
		}
		return fmt.Errorf("apiWithdrawFromElection get user: %w", err)
	}

	election, err := database.GetElection(electionID, tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return &fiber.Error{
				Code:    fiber.StatusNotFound,
				Message: "Election with that ID not found",
			}
		}
		return fmt.Errorf("apiWithdrawFromElection get election with id %d: %w", electionID, err)
	}

	candidate := &database.Candidate{
		UserID:     user.StudentID,
		ElectionID: electionID,
	}

	if err := candidate.Delete(tx); err != nil {
		return fmt.Errorf("apiWithdrawFromElection delete candidacy: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiWithdrawFromElection commit tx: %w", err)
	}

	events.SendEvent(events.TopicUserElectionWithdraw, &events.UserElectionWithdrawData{
		ActingUserID: actingUserID,
		User:         user,
		Election:     election,
	})

	ctx.Status(fiber.StatusNoContent)
	return nil
}

func (endpoints) apiGetElectionOutcome(ctx *fiber.Ctx) error {
	userID, authStatus := getSessionAuth(ctx)
	if authStatus == authNotAuthed {
		return fiber.ErrUnauthorized
	}

	electionID := ctx.QueryInt("election_id")

	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiGetElectionOutcome start tx: %w", err)
	}
	defer tx.Rollback()

	user, err := database.GetUser(userID, tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			// User has been deleted
			ctx.Cookie(newSessionTokenDeletionCookie())
			return fiber.ErrUnauthorized
		}
		return fmt.Errorf("apiGetElectionOutcome get user: %w", err)
	}

	electionOutcome, err := database.GetOutcomeForElection(electionID, tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return &fiber.Error{
				Code:    fiber.StatusNotFound,
				Message: "Outcome for election with that ID not found",
			}
		}
		return fmt.Errorf("apiGetElectionOutcome get election outcome for id %d: %w", electionID, err)
	}

	if !electionOutcome.IsPublished && !user.IsAdmin {
		return &fiber.Error{
			Code:    fiber.StatusForbidden,
			Message: "The outcome for this election has not been published yet",
		}
	}

	return ctx.JSON(electionOutcome)
}
