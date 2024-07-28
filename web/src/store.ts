import type { Optional } from "$lib/optional";
import { writable } from "svelte/store";

export interface User {
	name: string;
	studentID: string;
	isAdmin: boolean;
	isRestricted: boolean;
}

export interface Election {
	roleName: string;
	description: string;
	isActive: boolean;
	isConcluded: boolean;
	id: number;
	candidates: Optional<Array<{ name: string; isMe: boolean }>>;
}

export interface BallotEntry {
	id: number;
	name: string;
	isRON: boolean;
}

export interface CurrentElection {
	election: Election;
	ballot: Array<BallotEntry>;
	numEligibleVoters: number;
	hasVoted: boolean;
}

export interface ElectionOutcome {
	election: Election;
	date: number;
	ballots: number;
	rounds: number;
	isPublished: boolean;
	results: Array<ElectionOutcomeResult>;
}

export interface ElectionOutcomeResult {
	name: string;
	round: number;
	voteCount: number;
	isRejected: boolean;
	isElected: boolean;
}

export const user = writable<User>(null!);
export const elections = writable<Array<Election> | null>(null);
export const currentElection = writable<CurrentElection | null>(null);
export const fetching = writable<boolean>(false);
export const error = writable<Error | null>(null);
