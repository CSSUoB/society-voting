<script lang="ts">
	import Nav from "$lib/nav.svelte";
	import Profile from "./profile.svelte";
	import Upcoming from "./upcoming.svelte";
	import Current from "./current.svelte";
	import Users from "./users.svelte";
	import PresenterMode from "./presenterMode.svelte";
	import {
		user,
		type User,
		elections,
		type Election,
		type CurrentElection,
		currentElection,
	} from "../store";
	import Alert from "$lib/alert.svelte";
	import { API } from "$lib/endpoints";
	import { _getCurrentElection, _getElections } from "./+layout";

	/** @type {import('./$types').PageData} */
	export let data: { user: User; elections: Election[]; currentElection: CurrentElection | null };

	$: user.set(data.user);
	$: elections.set(data.elections);
	$: currentElection.set(data.currentElection);

	let menuOpen = false;
	const toggleMenu = (e: CustomEvent<boolean>) => {
		menuOpen = e.detail;
	};

	const eventSource = new EventSource(API.ELECTION_SSE, {
		withCredentials: true,
	});
	const electionStatusChanged = async () => {
		if (!$user.isAdmin) {
			$elections = await _getElections();
			$currentElection = await _getCurrentElection();
		}
	};
	eventSource.addEventListener("election-start", electionStatusChanged);
	eventSource.addEventListener("election-end", electionStatusChanged);
</script>

<div class="container">
	<Nav --location="nav" on:menuToggle={toggleMenu} />
	<span />
	<main style:left={menuOpen ? "0" : ""}>
		<div class="rail">
			<Profile />
			{#if $currentElection && (!$currentElection.hasVoted || $user.isAdmin)}
				<Current />
			{/if}
			<Upcoming />
			{#if data.user.isAdmin}
				<Users />
				<PresenterMode />
			{/if}
		</div>
		<div class="rail">
			<slot />
		</div>
	</main>
	<span />
	<footer>
		<p>
			Society Voting is open source software licensed under the BSD 2-Clause License. You can
			contribute <a href="https://github.com/CSSUoB/society-voting">here</a>.
		</p>
	</footer>
</div>
<Alert />

<style>
	div.container {
		height: 100vh;
		display: grid;
		grid-template-rows: auto 1fr auto;
		grid-template-columns: 1fr min(100%, 1600px) 1fr;
		grid-template-areas:
			"nav    nav    nav"
			".      main   ."
			"footer footer footer";
		background-color: #375db6;
		background-image: url($lib/assets/background.svg);
		background-size: min(1200px, 80vw);
		background-blend-mode: soft-light;
		overflow-y: auto;
		overflow-x: hidden;
		width: 100vw;
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

	footer {
		grid-area: footer;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 8px;
	}

	footer p {
		color: #fff;
		text-align: center;
	}

	@media only screen and (max-width: 850px) {
		main {
			grid-template-columns: min(100%, 0.3 * 850px) 100%;
			left: calc(-1 * min(100%, 0.3 * 850px) - 12px);
			position: relative;
			transition: left 0.2s ease-out;
		}
	}
</style>
