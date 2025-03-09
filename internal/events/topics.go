package events

import "github.com/CSSUoB/society-voting/internal/database"

const (
	TopicPollStarted  Topic = "poll-start"
	TopicPollEnded    Topic = "poll-end"
	TopicVoteReceived Topic = "vote-received"

	TopicUserElectionStand    Topic = "user-stand"
	TopicUserElectionWithdraw Topic = "user-withdraw"
	TopicUserRestricted       Topic = "user-restrict"
	TopicUserDeleted          Topic = "user-delete"
)

type UserElectionStandData struct {
	User     *database.User     `json:"user"`
	Election *database.Election `json:"election"`
}

type UserElectionWithdrawData struct {
	User         *database.User     `json:"user"`
	Election     *database.Election `json:"election"`
	ActingUserID string             `json:"actingUserID,omitempty"`
}

type PollEndedData struct {
	Poll   *database.Poll `json:"poll"`
	Name   string         `json:"name"`
	Result string         `json:"result,omitempty"`
}

type UserRestrictedData struct {
	User         *database.User `json:"user"`
	ActingUserID string         `json:"actingUserID"`
}

type UserDeletedData struct {
	UserID       string `json:"user"`
	ActingUserID string `json:"actingUserID"`
}
