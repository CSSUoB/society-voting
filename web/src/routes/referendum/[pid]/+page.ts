import { redirect } from "@sveltejs/kit";

export function load({ params }) {
    const id = parseInt(params.pid, 10);
    if (isNaN(id)) {
        throw redirect(302, "/");
    }
    return {
        id
    }
}