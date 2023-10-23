import { writable } from 'svelte/store'

export interface User {
    name: string;
    studentID: string;
    admin: boolean;
}

export const userStore = writable<User>(null!);