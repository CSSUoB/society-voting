import { API } from "$lib/endpoints";
import { redirect } from "@sveltejs/kit";

export const ssr = false;

export async function load({ fetch, params }) {
    const [user, polls, currentPoll] = await Promise.all([_getUser(), _getPolls(), _getCurrentPoll()]);

    return {
        user,
        polls,
        currentPoll
    };
}

const _getUser = async () => {
    const response = await fetch(API.ME);
    if (!response.ok) {
        throw redirect(302, API.AUTH_LOGIN);
    }
    return await response.json();
}

export const _getPolls = async () => {
    const response = await fetch(API.POLL);
    if (!response.ok) {
        throw redirect(302, API.AUTH_LOGIN);
    }
    return await response.json();
}

export const _getCurrentPoll = async () => {
    const response = await fetch(API.POLL_CURRENT);
    if (!response.ok) {
        return null;
    }
    return await response.json();
}