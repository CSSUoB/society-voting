<script lang="ts">
	import Nav from "$lib/nav.svelte";
	import Profile from "./profile.svelte";
	import Upcoming from "./upcoming.svelte";
	import Current from "./current.svelte";
	import Users from "./users.svelte";
	import { user, type User, elections, type Election } from "../store";

	/** @type {import('./$types').PageData} */
	export let data: { user: User; elections: Election[] };

	$: user.set(data.user);
	$: elections.set(data.elections);

	$: currentElection = $elections?.find((e) => e.isActive);
</script>

<div class="container">
	<Nav --location="nav" />
	<span />
	<main>
		<div class="rail">
			<Profile />
			{#if !data.user.admin && currentElection}
				<Current />
			{/if}
			<Upcoming />
			{#if data.user.admin}
				<Users />
			{/if}
		</div>
		<div class="rail">
			<slot />
		</div>
	</main>
	<span />
</div>

<style>
	div.container {
		height: 100svh;
		display: grid;
		grid-template-rows: auto 1fr;
		grid-template-columns: 1fr min(100vw, 1600px) 1fr;
		grid-template-areas:
			"nav nav  nav"
			".   main .";
		background-color: #375db6;
		background-image: url($lib/assets/background.svg);
		background-size: min(1200px, 80vw);
		background-blend-mode: soft-light;
		overflow-y: auto;
	}

	main {
		display: grid;
		grid-template-columns: min(30vw, 0.3 * 1600px) 1fr;
		grid-template-rows: auto 1fr;
		gap: 12px;
		padding: 12px;
	}

	div.rail {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}
</style>
