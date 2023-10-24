import { _getCurrentElection } from "../vote/+page";

export async function load({ fetch, params }) {
    return await _getCurrentElection();
}