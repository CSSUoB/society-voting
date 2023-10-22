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
	User     *database.User
	Election *database.Election
}

type UserElectionWithdrawData struct {
	User     *database.User
	Election *database.Election
	ByForce  bool
}

type ElectionEndedData struct {
	Election *database.Election
	Result   string
}
