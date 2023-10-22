package events

import "github.com/CSSUoB/society-voting/internal/database"

const (
	TopicElectionStarted Topic = "election-start"
	TopicElectionEnded   Topic = "election-end"
	TopicVoteReceived    Topic = "vote-received"

	TopicUserElectionStand    Topic = "user-stand"
	TopicUserElectionWithdraw Topic = "user-withdraw"
)

type UserElectionStandData struct {
	User     *database.User     `json:"user"`
	Election *database.Election `json:"election"`
}

type UserElectionWithdrawData struct {
	User     *database.User     `json:"user"`
	Election *database.Election `json:"election"`
	ByForce  bool               `json:"byForce"`
}

type ElectionEndedData struct {
	Election *database.Election `json:"election"`
	Result   string             `json:"result"`
}
