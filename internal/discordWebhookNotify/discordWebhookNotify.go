package discordWebhookNotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/CSSUoB/society-voting/internal/config"
	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/CSSUoB/society-voting/internal/events"
	"github.com/carlmjohnson/requests"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/textproto"
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
			me.Description = fmt.Sprintf("Election ID: %d", data.Election.ID)

			me.Files = append(me.Files, &EmbedFile{
				Filename: data.Election.RoleName + " " + time.Now().UTC().Format("2006-01-02") + ".txt",
				Content:  data.Result,
			})
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
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Colour      uint         `json:"color,omitempty"`
	Files       []*EmbedFile `json:"-"`
}

type EmbedFile struct {
	Filename string
	Content  string
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

	r := requests.URL(conf.URL).Method(http.MethodPost)

	if len(embed.Files) == 0 {
		r.BodyJSON(&req)
	} else {
		buf := new(bytes.Buffer)
		writer := multipart.NewWriter(buf)

		{
			part, _ := writer.CreatePart(textproto.MIMEHeader{"Content-Type": []string{"application/json"}, "Content-Disposition": []string{"form-data; name=\"payload_json\""}})
			dat, err := json.Marshal(&req)
			if err != nil {
				return err
			}
			_, _ = part.Write(dat)
		}

		for i, file := range embed.Files {
			part, _ := writer.CreateFormFile(fmt.Sprintf("files[%d]", i), file.Filename)
			_, _ = part.Write([]byte(file.Content))
		}

		_ = writer.Close()

		r.Header("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
		r.BodyBytes(buf.Bytes())
	}

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
