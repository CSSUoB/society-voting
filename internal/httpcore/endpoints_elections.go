package httpcore

import (
	"bufio"
	"crypto/subtle"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/CSSUoB/society-voting/internal/events"
	"github.com/gofiber/fiber/v2"
	"github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
)

func (endpoints) apiListPolls(ctx *fiber.Ctx) error {
	polls, err := database.GetAllPolls()
	if err != nil {
		return fmt.Errorf("apiListPolls get all polls: %w", err)
	}

	type PollWithData struct {
		database.Poll
		Candidates *[]*database.ElectionCandidate `json:"candidates,omitempty"`
	}

	var res []*PollWithData

	userID := ctx.Locals("userID").(string)
	for _, poll := range polls {
		if poll.Election != nil && !poll.IsConcluded {
			if ec, err := poll.Election.WithCandidates(); err != nil {
				return fmt.Errorf("apiListPolls: %w", err)
			} else {
				for _, cand := range ec.Candidates {
					cand.IsMe = cand.ID == userID
				}
				res = append(res, &PollWithData{
					Poll:       *poll,
					Candidates: &ec.Candidates,
				})
			}
		} else {
			res = append(res, &PollWithData{
				Poll: *poll,
			})
		}
	}

	return ctx.JSON(res)
}

func (endpoints) apiPollsSSE(ctx *fiber.Ctx) error {
	id, receiver := events.NewReceiver(events.TopicPollStarted, events.TopicPollEnded)

	ctx.Set("Content-Type", "text/event-stream")
	ctx.Set("Connection", "keep-alive")

	notify := ctx.Context().Done()
	fr := ctx.Response()
	fr.SetBodyStreamWriter(func(w *bufio.Writer) {
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()

	loop:
		for {
			select {
			case <-notify:
				break loop
			case msg := <-receiver:
				if msg.Topic == events.TopicPollEnded {
					// we're going to be modifying this msg so let's create a copy and work with that
					{
						// TODO: Refactor away this copying mess
						x := *msg
						y := *(msg.Data.(*events.PollEndedData))
						x.Data = &y
						msg = &x
					}
					msg.Data.(*events.PollEndedData).Result = ""
				}
				sseData, err := msg.ToSSE()
				if err != nil {
					slog.Error("SSE error", "error", fmt.Errorf("failed to generate SSE event from message: %w", err))
					break
				}
				if _, err := w.Write(sseData); err != nil {
					break loop
				}
			case <-ticker.C:
				// heartbeat
				if _, err := w.Write([]byte(":\n\n")); err != nil {
					break loop
				}
			}

			if err := w.Flush(); err != nil {
				// Client disconnected
				break loop
			}
		}
		events.CloseReceiver(id)
	})

	return nil
}

func (endpoints) apiGetActivePollInformation(ctx *fiber.Ctx) error {
	poll, err := database.GetActivePoll()
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return &fiber.Error{
				Code:    fiber.StatusConflict,
				Message: "There is no active poll.",
			}
		}
		return fmt.Errorf("apiVote get active poll: %wz", err)
	}

	numUsers, err := database.CountUsers()
	if err != nil {
		return fmt.Errorf("apiGetActivePollInformation count users: %w", err)
	}

	userID := ctx.Locals("userID").(string)
	hasVoted, err := database.HasUserVotedInPoll(userID, poll.ID)
	if err != nil {
		return fmt.Errorf("apiGetActivePollInformation check if user has voted: %w", err)
	}

	type BaseResponse struct {
		NumUsers int            `json:"numEligibleVoters"`
		HasVoted bool           `json:"hasVoted"`
		Poll     *database.Poll `json:"poll"`
	}

	baseResponse := BaseResponse{
		NumUsers: numUsers,
		HasVoted: hasVoted,
		Poll:     poll,
	}

	if poll.Election != nil {
		ballot, err := database.GetAllBallotEntriesForElection(poll.ID)
		if err != nil {
			return fmt.Errorf("apiGetActivePollInformation get ballot: %w", err)
		}

		// randomise ballot order
		sort.Slice(ballot, func(_, _ int) bool {
			return rand.Intn(2) == 0
		})

		return ctx.JSON(struct {
			BaseResponse
			Ballot []*database.BallotEntry `json:"ballot"`
		}{
			BaseResponse: baseResponse,
			Ballot:       ballot,
		})
	}

	return ctx.JSON(baseResponse)
}

func apiVote(ctx *fiber.Ctx, fetchPoll func(int, bun.Tx) (*database.Poll, error), validateVote func(int, []int, bun.Tx) error) error {
	userID := ctx.Locals("userID").(string)

	var request = struct {
		ID   int    `json:"id" validate:"required"`
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

	poll, err := fetchPoll(request.ID, tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return &fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: "Poll with that ID not found or is wrong type",
			}
		}
		return fmt.Errorf("apiVote get poll: %w", err)
	}

	if !poll.IsActive {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Poll with that ID is not active",
		}
	}

	hasVotedAlready, err := database.HasUserVotedInPoll(userID, poll.ID, tx)
	if err != nil {
		return fmt.Errorf("apiVote check if user %s has already voted: %w", userID, err)
	}
	if hasVotedAlready {
		return &fiber.Error{
			Code:    fiber.StatusConflict,
			Message: "You have already voted. Go away :)",
		}
	}

	if err := validateVote(request.ID, request.Vote, tx); err != nil {
		return err
	}

	if err := (&database.Vote{
		PollID:  poll.ID,
		UserID:  userID,
		Choices: request.Vote,
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

func (endpoints) apiVoteInElection(ctx *fiber.Ctx) error {
	fetchElection := func(id int, tx bun.Tx) (*database.Poll, error) {
		election, err := database.GetElection(id, tx)
		if err != nil {
			return nil, err
		}
		return election.Poll, nil
	}

	validateElectionVote := func(pollId int, vote []int, tx bun.Tx) error {
		ballotOptions, err := database.GetAllBallotEntriesForElection(pollId, tx)
		if err != nil {
			return fmt.Errorf("get all ballot entries: %w", err)
		}
		for _, id := range vote {
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
		return nil
	}

	return apiVote(ctx, fetchElection, validateElectionVote)
}

func (endpoints) apiVoteInReferendum(ctx *fiber.Ctx) error {
	fetchReferendum := func(id int, tx bun.Tx) (*database.Poll, error) {
		referendum, err := database.GetReferendum(id, tx)
		if err != nil {
			return nil, err
		}
		return referendum.Poll, nil
	}

	validateReferendumVote := func(pollId int, vote []int, tx bun.Tx) error {
		if len(vote) != 1 || !slices.Contains([]int{0, 1, 2}, vote[0]) {
			return &fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid vote",
			}
		}
		return nil
	}

	return apiVote(ctx, fetchReferendum, validateReferendumVote)
}

func (endpoints) apiStandForElection(ctx *fiber.Ctx) error {
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

	user := ctx.Locals("user").(*database.User)

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

	if election.Poll == nil || election.Poll.IsConcluded || election.Poll.IsActive {
		return &fiber.Error{
			Code:    fiber.StatusConflict,
			Message: "You can no longer stand in this election",
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

	actingUserID := ctx.Locals("userID").(string)
	authStatus := ctx.Locals("authStatus").(authType)

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

func (endpoints) apiGetPollOutcome(ctx *fiber.Ctx) error {
	electionID := ctx.QueryInt("id")

	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiGetElectionOutcome start tx: %w", err)
	}
	defer tx.Rollback()

	electionOutcome, err := database.GetOutcomeForPoll(electionID, tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return &fiber.Error{
				Code:    fiber.StatusNotFound,
				Message: "Outcome for poll with that ID not found",
			}
		}
		return fmt.Errorf("apiGetElectionOutcome get election outcome for id %d: %w", electionID, err)
	}

	authStatus := ctx.Locals("authStatus").(authType)
	if !electionOutcome.IsPublished && authStatus&authAdminUser == 0 {
		return &fiber.Error{
			Code:    fiber.StatusForbidden,
			Message: "The outcome for this election has not been published yet",
		}
	}

	return ctx.JSON(electionOutcome)
}
