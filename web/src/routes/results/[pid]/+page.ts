import { API } from "$lib/endpoints";
import { error, redirect } from "@sveltejs/kit";

export async function load({ fetch, params }) {
	const id = parseInt(params.pid, 10);
	if (isNaN(id)) {
		throw redirect(302, "/");
	}

	const response = await fetch(`${API.POLL_RESULTS}?id=${id}`);
	if (!response.ok) {
		throw error(response.status, { message: await response.text() });
	}
	const jsonResponse = await response.json();

	return { data: jsonResponse };
}
