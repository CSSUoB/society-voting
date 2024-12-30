<script lang="ts">
	import Panel from "$lib/panel.svelte";
	import { user, fetching, currentPoll, error } from "../../store";
	import Button from "$lib/button.svelte";
	import { getEndpointForPollType } from "$lib/endpoints";
	import { goto } from "$app/navigation";
	import Dialog from "$lib/dialog.svelte";
	import { _getCurrentPoll } from "../+layout";
	import InstantRunoffBallot from "./instant-runoff-ballot.svelte";
	import ReferendumBallot from "./referendum-ballot.svelte";
	import PollHeader from "$lib/poll-header.svelte";

	let codeInput: HTMLInputElement;
	let votedDialog: HTMLDialogElement;
	let validBallot: boolean;
	let choices: Array<number | undefined>;

	$: if (!$currentPoll) {
		goto("/");
	}
	
	const ballotUpdate = (e: CustomEvent<any>) => {
		choices = e.detail.choices;
		validBallot = e.detail.valid;
	}

	const submit = async () => {
		if (!validBallot || !$currentPoll) return;

		$fetching = true;
		const url = getEndpointForPollType("vote", $currentPoll.poll.pollType.id);
		console.log($currentPoll)
		if (!url) {
			$error = new Error(`Cannot vote for unknown poll type "${$currentPoll.poll.pollType.name}""`);
			$fetching = false;
			return;
		}

		const response = await fetch(url, {
			method: "POST",
			body: JSON.stringify({ id: $currentPoll.poll.id, vote: choices, code: codeInput.value.trim().toUpperCase() }),
		});

		if (!response.ok) {
			$error = new Error(await response.text());
			$fetching = false;
			return;
		}
		currentPoll.set(await _getCurrentPoll());

		$fetching = false;
		votedDialog.showModal();
	};
</script>

{#if $currentPoll}
	<PollHeader prefix="Voting" poll={$currentPoll.poll}></PollHeader>
{/if}

{#if $currentPoll?.poll.election && $currentPoll.ballot}
	<InstantRunoffBallot
		candidates={$currentPoll.ballot}	
		on:update={ballotUpdate}
	></InstantRunoffBallot>
{:else if $currentPoll?.poll.referendum}
	<ReferendumBallot on:update={ballotUpdate}></ReferendumBallot>
{/if}

<Panel title="Submit" kind="emphasis">
	<div class="submit-container">
		<input bind:this={codeInput} placeholder="Enter election code" type="text" />
		<Button
			kind="primary"
			text="Submit vote"
			icon="check"
			on:click={() => submit()}
			disabled={!validBallot}
		/>
	</div>
</Panel>

<Dialog
	title="Congrats, you've voted!"
	bind:dialog={votedDialog}
	on:close={() => goto($user.isAdmin ? "/stats" : "/")}
>
	<div class="dialog-container">
		<img src={`https://cssuob.github.io/resources/dinosaur/tex_ballot.svg`} width="200px" />
		<p>
			Don't forget to grab the special edition <strong>voting TeX sticker</strong> afterwards!
		</p>
	</div>
	<Button slot="actions" text="Close" kind="emphasis" />
</Dialog>

<style>
	.submit-container {
		display: flex;
		flex-direction: column;
		gap: 8px;
		align-items: flex-start;
	}

	.submit-container > input {
		height: 36px;
		padding: 2px 12px;
		border-radius: 8px;
		border: 2px solid #000;
		text-transform: uppercase;
	}

	.dialog-container {
		display: flex;
		flex-direction: column;
		gap: 8px;
		align-items: center;
		max-width: min(100vw, 400px);
		text-align: center;
	}

	.dialog-container img {
		padding-right: 50px;
	}
</style>
