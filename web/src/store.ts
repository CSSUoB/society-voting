import { writable } from 'svelte/store'

export interface User {
    name: string;
    studentID: number;
}

export const userStore = writable<User>(null!);