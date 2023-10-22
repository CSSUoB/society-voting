import { writable } from 'svelte/store'

export interface User {
    name: string;
    studentID: string;
}

export const userStore = writable<User>(null!);