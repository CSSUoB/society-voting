<script lang="ts">
	import List from "$lib/list.svelte";
	import Panel from "$lib/panel.svelte";
	import { user, type CurrentElection, fetching, currentElection } from "../../store";
	import BallotEntry from "./ballot-entry.svelte";
	import type { BallotEntry as BallotEntryT } from "../../store";
	import Button from "$lib/button.svelte";
	import { API } from "$lib/endpoints";
	import { goto } from "$app/navigation";
	import Dialog from "$lib/dialog.svelte";

	let ballot: Array<BallotEntryT | undefined> = Array.from(Array($currentElection?.ballot.length));
	let errors = Array.from(Array(ballot.length));
	let codeInput: HTMLInputElement;
	let votedDialog: HTMLDialogElement;

	$: if (!$currentElection) {
		goto("/");
	}

	const updateAndValidate = (index: number, changeEvent: Event) => {
		const id: number = parseInt((changeEvent.target as HTMLSelectElement).value);
		const candidate = $currentElection?.ballot.find((c) => c.id === id);
		ballot = ballot.map((b, i) => (i === index ? candidate : b));

		errors = ballot.map((b, i) => {
			if (b === undefined && ballot.slice(i).filter((x) => x).length > 0)
				return "You cannot have gaps in your ranking";
			if (b === undefined) return undefined;
			if (ballot.filter((bb) => bb?.id === b.id).length !== 1)
				return `You cannot rank ${b.isRON ? b.name : b.name.split(" ")[0]} more than once`;
			return undefined;
		});
	};

	const submit = async () => {
		if (errors.filter((x) => x).length > 0) return;
		$fetching = true;
		const votes = ballot.filter((x) => x).map((b) => b?.id);
		const response = await fetch(API.ELECTION_CURRENT_VOTE, {
			method: "POST",
			body: JSON.stringify({ vote: votes, code: codeInput.value.trim().toUpperCase() }),
		});
		if (!response.ok) {
			// show error code
			return;
		}

		$fetching = false;
		votedDialog.showModal();
	};
</script>

<svelte:head>
	<title>Vote for: {$currentElection?.election.roleName}</title>
</svelte:head>

<Panel title={`Electing: ${$currentElection?.election.roleName}`}>
	<p>{$currentElection?.election.description}</p>
</Panel>

{#if !$user.admin}
	<Panel title="Your ballot">
		<p>There are {($currentElection?.ballot.length ?? 1) - 1} candidates on the ballot.</p>
		<p>Rank candidates in order of your choice. You do not need to rank all candidates.</p>
		<List items={ballot} let:prop={candidate}>
			<BallotEntry
				ballot={$currentElection?.ballot ?? []}
				{candidate}
				error={errors[candidate.index]}
				on:change={(e) => updateAndValidate(candidate.index, e)}
			/>
		</List>
	</Panel>

	<Panel title="Submit" kind="emphasis">
		<div class="submit-container">
			<input bind:this={codeInput} placeholder="Enter election code" type="text" />
			<Button
				kind="primary"
				text="Submit vote"
				icon="check"
				on:click={() => submit()}
				disabled={errors.filter((x) => x).length > 0}
			/>
		</div>
	</Panel>
{/if}

<Dialog title="Congrats, you've voted!" bind:dialog={votedDialog} on:close={() => goto("/")}>
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
