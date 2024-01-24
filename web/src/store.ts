import type { Optional } from '$lib/optional';
import { writable } from 'svelte/store'

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
    id: number;
    candidates: Optional<Array<{ name: string, isMe: boolean }>>;
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

export const user = writable<User>(null!);
export const elections = writable<Array<Election> | null>(null);
export const currentElection = writable<CurrentElection | null>(null);
export const fetching = writable<boolean>(false);
export const error = writable<Error | null>(null);