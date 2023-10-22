import { API } from "$lib/endpoints";
import { redirect } from "@sveltejs/kit";

export const ssr = false;

export async function load({ fetch, params }) {
    const response = await fetch(API.ME);
    if (!response.ok) {
        throw redirect(302, "/auth/login");
    }
    return {
        user: await response.json()
    };
}