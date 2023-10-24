import { API } from "$lib/endpoints";
import { redirect } from "@sveltejs/kit";

export async function load({ fetch, params }) {
    return await _getCurrentElection();
}

export const _getCurrentElection = async () => {
    const response = await fetch(API.ELECTION_CURRENT);
    if (!response.ok) {
        throw redirect(302, "/");
    }
    return await response.json();
}