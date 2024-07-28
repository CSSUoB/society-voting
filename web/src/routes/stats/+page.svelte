<script lang="ts">
	import { goto } from "$app/navigation";
	import Button from "$lib/button.svelte";
	import Dialog from "$lib/dialog.svelte";
	import { API } from "$lib/endpoints";
	import Panel from "$lib/panel.svelte";
	import { _getCurrentElection, _getElections } from "../+layout";
	import { error, fetching, currentElection, elections } from "../../store";

	let electionRunning = true;
	let numberOfVotes = 0;
	let dialog: HTMLDialogElement;

	if (!$currentElection) {
		goto("/");
	}

	const eventSource = new EventSource(API.ADMIN_ELECTION_SSE, {
		withCredentials: true,
	});
	eventSource.addEventListener("vote-received", (e) => {
		numberOfVotes = parseInt(e.data, 10);
	});

	const endElection = async () => {
		$fetching = true;
		const response = await fetch(API.ADMIN_ELECTION_STOP, {
			method: "POST",
		});
		if (!response.ok) {
			$fetching = false;
			$error = new Error(await response.text());
			return;
		}
		electionRunning = false;
		let electionId = (await response.json()).election.id;
		$elections = await _getElections();
		$currentElection = await _getCurrentElection();
		$fetching = false;
		goto(`/results/${electionId}`);
	};
</script>

<svelte:head>
	<title>Vote for: {$currentElection?.election.roleName}</title>
</svelte:head>

<Panel title={`Electing: ${$currentElection?.election.roleName}`}>
	<p>{$currentElection?.election.description}</p>
</Panel>

{#if electionRunning}
	<Panel title="Admin actions" headerIcon="admin_panel_settings">
		<div class="container">
			<h3>{numberOfVotes} of {$currentElection?.numEligibleVoters} users have voted so far</h3>
			<p>
				The turnout so far is {(
					(numberOfVotes * 100) /
					($currentElection?.numEligibleVoters ?? 100)
				).toFixed(2)}%
			</p>
			<Button
				text="End election and view results"
				kind="primary"
				on:click={() => dialog.showModal()}
			/>
		</div>
	</Panel>
	<Dialog bind:dialog title="End election and view results?" on:submit={endElection}>
		<svelte:fragment slot="actions">
			<Button text="Cancel" />
			<Button text="End election" kind="emphasis" name="submit" />
		</svelte:fragment>
	</Dialog>
{/if}

<style>
	div.container {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: 8px;
	}

	p {
		white-space: break-spaces;
	}
</style>
