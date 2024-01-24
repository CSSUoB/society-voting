import { API } from "$lib/endpoints";
import { redirect } from "@sveltejs/kit";

export async function load({ fetch, params }) {
    const response = await fetch(API.ADMIN_USER);
    if (!response.ok) {
        throw redirect(302, "/");
    }
}