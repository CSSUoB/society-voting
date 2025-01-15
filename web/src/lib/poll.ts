import type { ElectionPoll, ElectionPollOutcome, Poll, PollOutcome, ReferendumPoll, ReferendumPollOutcome } from "../store"

export function isElectionPoll(poll: Poll): poll is ElectionPoll {
    return poll.pollType.name === "Election";
}

export function isReferendumPoll(poll: Poll): poll is ReferendumPoll {
    return poll.pollType.name === "Referendum";
}

export function isElectionPollOutcome(pollOutcome: PollOutcome): pollOutcome is ElectionPollOutcome {
    return pollOutcome.poll.pollType.name === "Election";
}

export function isReferendumPollOutcome(pollOutcome: PollOutcome): pollOutcome is ReferendumPollOutcome {
    return pollOutcome.poll.pollType.name === "Referendum";
}

export function getFriendlyName(poll: Poll): string {
    if (isElectionPoll(poll)) {
        return poll.election.roleName;
    } else {
        return poll.referendum.title;
    }
}