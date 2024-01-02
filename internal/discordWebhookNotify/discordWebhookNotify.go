package discordWebhookNotify

import (
	"context"
	"fmt"
	"github.com/CSSUoB/society-voting/internal/config"
	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/CSSUoB/society-voting/internal/events"
	"github.com/carlmjohnson/requests"
	"log/slog"
	"net/http"
	"time"
)

func Run() {
	slog.Info("starting Discord webhook event notifier")

	_, standReceiver := events.NewReceiver(events.TopicUserElectionStand)
	_, withdrawReceiver := events.NewReceiver(events.TopicUserElectionWithdraw)
	_, electionStart := events.NewReceiver(events.TopicElectionStarted)
	_, electionEnd := events.NewReceiver(events.TopicElectionEnded)

	for {
		me := &MinimalEmbed{}
		select {
		case msg := <-standReceiver:
			data := msg.Data.(*events.UserElectionStandData)

			me.Colour = ColourGood
			me.Title = fmt.Sprintf("%s is now standing for election as %s", data.User.Name, data.Election.RoleName)
			me.Description = fmt.Sprintf("User ID: %s\nElection ID: %d", data.User.StudentID, data.Election.ID)

		case msg := <-withdrawReceiver:
			data := msg.Data.(*events.UserElectionWithdrawData)

			me.Colour = ColourBad
			me.Title = fmt.Sprintf("%s is no longer standing for election as %s", data.User.Name, data.Election.RoleName)
			me.Description = fmt.Sprintf("User ID: %s\nElection ID: %d", data.User.StudentID, data.Election.ID)

			if data.ByForce {
				me.Description += "\n\n**User was removed by an admin**"
			}

		case msg := <-electionStart:
			data := msg.Data.(*database.Election)

			me.Colour = ColourMid
			me.Title = fmt.Sprintf("Voting for %s has started!", data.RoleName)
			me.Description = fmt.Sprintf("Election ID: %d", data.ID)

		case msg := <-electionEnd:
			data := msg.Data.(*events.ElectionEndedData)

			me.Colour = ColourGood
			me.Title = fmt.Sprintf("Voting for %s has ended!", data.Election.RoleName)
			me.Description = fmt.Sprintf("Election ID: %d\n\n```%s```", data.Election.ID, data.Result)
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

	var resp string
	r.ToString(&resp)

	if err := r.Fetch(ctx); err != nil {
		return err
	}

	fmt.Println(resp)

	return nil
}
