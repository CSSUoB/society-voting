package discordWebhookNotify

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/CSSUoB/society-voting/internal/config"
	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/CSSUoB/society-voting/internal/events"
	"github.com/carlmjohnson/requests"
)

func Run() {
	slog.Info("starting Discord webhook event notifier")

	_, receiver := events.NewReceiver(
		events.TopicUserElectionStand,
		events.TopicUserElectionWithdraw,
		events.TopicPollStarted,
		events.TopicPollEnded,
		events.TopicUserRestricted,
		events.TopicUserDeleted,
	)

	for msg := range receiver {
		me := &MinimalEmbed{}
		switch msg.Topic {
		case events.TopicUserElectionStand:
			data := msg.Data.(*events.UserElectionStandData)

			me.Colour = ColourGood
			me.Title = fmt.Sprintf("%s is now standing for election as %s", data.User.Name, data.Election.RoleName)
			me.Description = fmt.Sprintf("User ID: %s\nElection ID: %d", data.User.StudentID, data.Election.ID)

		case events.TopicUserElectionWithdraw:
			data := msg.Data.(*events.UserElectionWithdrawData)

			me.Colour = ColourBad
			me.Title = fmt.Sprintf("%s is no longer standing for election as %s", data.User.Name, data.Election.RoleName)
			me.Description = fmt.Sprintf("User ID: %s\nElection ID: %d", data.User.StudentID, data.Election.ID)

			if data.ActingUserID != data.User.StudentID {
				me.Description += fmt.Sprintf("\n\n**User was removed by admin with ID %s**", data.ActingUserID)
			}

		case events.TopicPollStarted:
			data := msg.Data.(database.Pollable)

			me.Colour = ColourMid
			me.Title = fmt.Sprintf("Voting for %s has started!", data.GetFriendlyTitle())
			me.Description = fmt.Sprintf("Poll ID: %d", data.GetPoll().ID)

		case events.TopicPollEnded:
			data := msg.Data.(*events.PollEndedData)

			me.Colour = ColourGood
			me.Title = fmt.Sprintf("Voting for %s has ended!", data.Name)
			me.Description = fmt.Sprintf("Election ID: %d\n\n```%s```", data.Poll.ID, data.Result)

		case events.TopicUserRestricted:
			data := msg.Data.(*events.UserRestrictedData)

			me.Colour = ColourMid
			me.Title = fmt.Sprintf("%s has been restricted", data.User.Name)
			me.Description = fmt.Sprintf("User ID: %s\n\nRestricted by admin with ID %s", data.User.StudentID, data.ActingUserID)

		case events.TopicUserDeleted:
			data := msg.Data.(*events.UserDeletedData)

			me.Colour = ColourMid
			me.Title = fmt.Sprintf("Account with ID %s has been deleted", data.UserID)
			me.Description = fmt.Sprintf("Deleted by admin with ID %s", data.ActingUserID)
		}

		if err := sendEmbed(me); err != nil {
			slog.Error("failed to send Discord webhook", "error", err)
		}
	}
}

const (
	ColourGood uint = 0x55ff55
	ColourBad  uint = 0xff5555
	ColourMid  uint = 0xffff55
)

type MinimalEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Colour      uint   `json:"color,omitempty"`
}

func sendEmbed(embed *MinimalEmbed) error {
	conf := config.Get().Platform.DiscordWebhook

	type embedWrapper struct {
		Type string `json:"type"`
		*MinimalEmbed
	}

	var req = struct {
		Embeds []*embedWrapper `json:"embeds"`
	}{
		Embeds: []*embedWrapper{
			{Type: "rich", MinimalEmbed: embed},
		},
	}

	r := requests.URL(conf.URL).Method(http.MethodPost).BodyJSON(&req)

	if conf.ThreadID != "" {
		r.Param("thread_id", conf.ThreadID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := r.Fetch(ctx); err != nil {
		return err
	}

	return nil
}
