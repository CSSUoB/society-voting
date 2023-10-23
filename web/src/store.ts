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

export const user = writable<User>(null!);
export const elections = writable<Array<Election>>([]);
