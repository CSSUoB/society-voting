import type { Optional } from '$lib/optional';
import { writable } from 'svelte/store'

export interface User {
    name: string;
    studentID: string;
    admin: boolean;
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
}

export const user = writable<User>(null!);
export const elections = writable<Array<Election> | null>(null);
