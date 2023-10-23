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
    candidates: null;
}

export const userStore = writable<User>(null!);
export const electionStore = writable<Array<Election>>([]);
