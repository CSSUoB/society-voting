import { API } from "$lib/endpoints";
import { redirect } from "@sveltejs/kit";

export async function load({ fetch, params }) {
    const election = await _getCurrentElection();
    if (election.hasVoted) {
        throw redirect(302, "/");
    }
    return election;
}

export const _getCurrentElection = async () => {
    const response = await fetch(API.ELECTION_CURRENT);
    if (!response.ok) {
        throw redirect(302, "/");
    }
    return await response.json();
}