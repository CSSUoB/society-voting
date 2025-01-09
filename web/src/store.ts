import { writable } from "svelte/store";

export interface User {
	name: string;
	studentID: string;
	isAdmin: boolean;
	isRestricted: boolean;
}

export interface PollType {
	id: number;
	name: string;
}

export type Poll = ElectionPoll | ReferendumPoll;

export interface ElectionPoll {
    id: number;
    isActive: boolean;
    isConcluded: boolean;
    pollType: PollType; 
    election: Election;
	candidates: Array<{ name: string; isMe: boolean }>;
}

export interface ReferendumPoll {
    id: number;
    isActive: boolean;
    isConcluded: boolean;
    pollType: PollType; 
    referendum: Referendum;
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
	ballot?: Array<BallotEntry>;
}

export type PollOutcome = ElectionPollOutcome | ReferendumPollOutcome;

export interface ElectionPollOutcome {
	date: string;
	isPublished: boolean;
	poll: ElectionPoll;
	ballots: number;
	electionOutcome: ElectionOutcome;
}

export interface ReferendumPollOutcome {
	date: string;
	isPublished: boolean;
	poll: ReferendumPoll;
	ballots: number;
	referendumOutcome: ReferendumOutcome;
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
