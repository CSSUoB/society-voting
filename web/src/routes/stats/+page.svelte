<script lang="ts">
	import { goto } from "$app/navigation";
	import Button from "$lib/button.svelte";
	import Dialog from "$lib/dialog.svelte";
	import { API, getEndpointForPollType } from "$lib/endpoints";
	import Panel from "$lib/panel.svelte";
	import PollHeader from "$lib/poll-header.svelte";
	import { _getCurrentPoll, _getPolls } from "../+layout";
	import { error, fetching, currentPoll, polls } from "../../store";

	let electionRunning = true;
	let numberOfVotes = 0;
	let dialog: HTMLDialogElement;

	if (!$currentPoll) {
		goto("/");
	}

	const eventSource = new EventSource(API.ADMIN_POLL_SSE, {
		withCredentials: true,
	});
	eventSource.addEventListener("vote-received", (e) => {
		numberOfVotes = parseInt(e.data, 10);
	});

	const endElection = async () => {
		if (!$currentPoll) return;

		$fetching = true;
		const url = getEndpointForPollType("stop", $currentPoll.poll.pollType.id);
		if (!url) {
			$error = new Error(`Cannot vote for unknown poll type "${$currentPoll.poll.pollType.name}""`);
			$fetching = false;
			return;
		}
		const response = await fetch(url, {
			method: "POST",
		});
		if (!response.ok) {
			$fetching = false;
			$error = new Error(await response.text());
			return;
		}
		electionRunning = false;
		$polls = await _getPolls();
		let id = $currentPoll.poll.id;
		$currentPoll = await _getCurrentPoll();
		$fetching = false;

		goto(`/results/${id}`);
	};
</script>

{#if $currentPoll}
	<PollHeader prefix="Voting" poll={$currentPoll.poll}></PollHeader>
{/if}

{#if electionRunning}
	<Panel title="Admin actions" headerIcon="admin_panel_settings">
		<div class="container">
			<h3>{numberOfVotes} of {$currentPoll?.numEligibleVoters} users have voted so far</h3>
			<p>
				The turnout so far is {(
					(numberOfVotes * 100) /
					($currentPoll?.numEligibleVoters ?? 100)
				).toFixed(2)}%
			</p>
			<Button
				text="End poll and view results"
				kind="primary"
				on:click={() => dialog.showModal()}
			/>
		</div>
	</Panel>
	<Dialog bind:dialog title="End poll and view results?" on:submit={endElection}>
		<svelte:fragment slot="actions">
			<Button text="Cancel" />
			<Button text="End poll" kind="emphasis" name="submit" />
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
