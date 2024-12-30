import type { Optional } from "$lib/optional";
import { writable } from "svelte/store";

export interface User {
	name: string;
	studentID: string;
	isAdmin: boolean;
	isRestricted: boolean;
}

export interface Poll {
	id: number;
	isActive: boolean;
	isConcluded: boolean;
	pollType: PollType;
	election: Optional<Election>;
	referendum: Optional<Referendum>;
	candidates: Optional<Array<{ name: string; isMe: boolean }>>;
}

export interface PollType {
	id: number;
	name: string;
}

export interface Election {
	id: number;
	roleName: string;
	description: string;
}

export interface Referendum {
	title: string;
	question: string;
	description: string;
}

export interface BallotEntry {
	id: number;
	name: string;
	isRON: boolean;
}

export interface CurrentPoll {
	numEligibleVoters: number;
	hasVoted: boolean;
	poll: Poll;
	ballot: Optional<Array<BallotEntry>>;
}

export interface PollOutcome {
	date: string;
	isPublished: boolean;
	poll: Poll;
	ballots: number;
	electionOutcome: Optional<ElectionOutcome>;
	referendumOutcome: Optional<ReferendumOutcome>;
}

export interface ElectionOutcome {
	rounds: number;
	results: Array<ElectionOutcomeResult>;
}

export interface ReferendumOutcome {
	votesFor: number;
	votesAgainst: number;
	votesAbstain: number;
}

export interface ElectionOutcomeResult {
	name: string;
	round: number;
	voteCount: number;
	isRejected: boolean;
	isElected: boolean;
}

export const user = writable<User>(null!);
export const polls = writable<Array<Poll> | null>(null);
export const currentPoll = writable<CurrentPoll | null>(null);
export const fetching = writable<boolean>(false);
export const navigating = writable<boolean>(false);
export const error = writable<Error | null>(null);
