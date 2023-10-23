import { API } from "$lib/endpoints";
import { redirect } from "@sveltejs/kit";

export const ssr = false;

export async function load({ fetch, params }) {
    const [user, elections] = await Promise.all([_getUser(), _getElections()]);

    return {
        user: { ...user, admin: user.studentID === "admin" },
        elections
    };
}

const _getUser = async () => {
    const response = await fetch(API.ME);
    if (!response.ok) {
        throw redirect(302, API.AUTH_LOGIN);
    }
    return await response.json();
}

export const _getElections = async () => {
    const response = await fetch(API.ELECTION);
    if (!response.ok) {
        throw redirect(302, API.AUTH_LOGIN);
    }
    return await response.json();
}