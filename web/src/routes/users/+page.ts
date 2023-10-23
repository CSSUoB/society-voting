import { API } from "$lib/endpoints";
import { redirect } from "@sveltejs/kit";

export async function load({ fetch, params }) {
    return {
        users: await getUsers(),
    };
}

const getUsers = async () => {
    const response = await fetch(API.ADMIN_USER);
    if (!response.ok) {
        throw redirect(302, "/");
    }
    return await response.json();
}