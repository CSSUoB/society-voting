<script lang="ts">
	import Avatar from "$lib/avatar.svelte";
	import Banner from "$lib/banner.svelte";
	import Button from "$lib/button.svelte";
	import Dialog from "$lib/dialog.svelte";
	import List from "$lib/list.svelte";
	import Panel from "$lib/panel.svelte";

	import run from "$lib/assets/run_for_election.svg";
	import { currentElection, elections, error, fetching, user } from "../../../store";
	import { goto } from "$app/navigation";
	import { API } from "$lib/endpoints";
	import { _getCurrentElection, _getElections } from "../../+layout";

	export let data: { id: number };
	$: election = $elections?.find((e) => e.id === data.id);
	$: if (!election) {
		goto("/", { replaceState: true });
	} else if (election.isActive) {
		goto("/vote");
	}

	let buttonText = "Stand for election";
	$: buttonText = `Stand for ${election ? election.roleName : "election"}`;

	const standOrWithdraw = async (id: number, stand: boolean) => {
		$fetching = true;
		const response = await fetch(API.ELECTION_STAND, {
			method: stand ? "POST" : "DELETE",
			body: JSON.stringify({ id }),
		});

		if (response.ok) {
			elections.set(await _getElections());
		} else {
			$error = new Error(await response.text());
		}
		$fetching = false;
	};

	let floorCandidatesInput: HTMLTextAreaElement;
	let startElectionDialog: HTMLDialogElement;
	const startElection = async (id: number) => {
		$fetching = true;
		const extraNames = floorCandidatesInput.value
			.trim()
			.split("\n")
			.filter((x) => x)
			.map((x) => x.trim());
		const response = await fetch(API.ADMIN_ELECTION_START, {
			method: "POST",
			body: JSON.stringify({ id, extraNames }),
		});

		if (response.ok) {
			elections.set(await _getElections());
			currentElection.set(await _getCurrentElection());
			await new Promise((r) => setTimeout(r, 200));
			goto("/stats");
		} else {
			$error = new Error(await response.text());
		}
		$fetching = false;
	};

	let deleteElectionDialog: HTMLDialogElement;
	const deleteElection = async (id: number) => {
		$fetching = true;
		const response = await fetch(API.ADMIN_ELECTION, {
			method: "DELETE",
			body: JSON.stringify({ id }),
		});

		if (response.ok) {
			elections.set(await _getElections());
			goto("/");
		} else {
			$error = new Error(await response.text());
		}
		$fetching = false;
	};
</script>

<svelte:head>
	<title>Electing: {election?.roleName}</title>
</svelte:head>
<Panel title={`Electing: ${election?.roleName}`}>
	<p>{election?.description}</p>
</Panel>

{#if !election?.candidates?.some((c) => c.isMe)}
	<Banner title="Interested in running?" kind="emphasis">
		<img slot="image" src={run} alt="" class="banner-image" />
		<svelte:fragment slot="content">
			<p>
				If you're thinking of running for election, go for it! Being on committee is super
				rewarding, super good fun, and a great way to bolster your CV and give back to the
				community.
			</p>
			<br />
			<Button
				text={buttonText}
				kind="primary"
				on:click={() => standOrWithdraw(election?.id ?? -1, true)}
			/>
		</svelte:fragment>
	</Banner>
{/if}

<Panel title="Candidates" headerIcon="groups">
	<List items={election?.candidates ?? []} let:prop={candidate}>
		<li class="candidate">
			<Avatar name={candidate.name} />
			<p>
				{candidate.name}{#if candidate.isMe}
					<span><small>You</small></span>
				{/if}
			</p>
			{#if candidate.isMe}
				<Button text="Withdraw" on:click={() => standOrWithdraw(election?.id ?? -1, false)} />
			{/if}
		</li>
	</List>
	{#if (election?.candidates?.length ?? 0) === 0}
		<p>There are no candidates currently running in this election</p>
	{/if}
</Panel>

{#if $user.isAdmin}
	<Panel title="Admin stuff" headerIcon="admin_panel_settings">
		<div class="admin-controls">
			<h3>Candidates standing from the floor</h3>
			<textarea
				bind:this={floorCandidatesInput}
				placeholder="Enter each candidate's name in a new line"
			/>
			<Button
				kind="primary"
				text="Save candidates and start election"
				on:click={() => startElectionDialog.showModal()}
			/>
			<Button text="Delete election" on:click={() => deleteElectionDialog.showModal()} />
		</div>
	</Panel>
	<Dialog
		bind:dialog={startElectionDialog}
		title="Confirm candidates and start election?"
		on:submit={() => startElection(election?.id ?? -1)}
	>
		<svelte:fragment slot="actions">
			<Button text="Cancel" />
			<Button text="Start election" kind="emphasis" name="submit" />
		</svelte:fragment>
	</Dialog>
	<Dialog
		bind:dialog={deleteElectionDialog}
		title="Delete election?"
		on:submit={() => deleteElection(election?.id ?? -1)}
	>
		<p>This action cannot be undone.</p>
		<svelte:fragment slot="actions">
			<Button text="Cancel" />
			<Button text="Delete election" kind="emphasis" name="submit" />
		</svelte:fragment>
	</Dialog>
{/if}

<style>
	img.banner-image {
		height: 120px;
		transform: translateY(16px);
	}

	li.candidate {
		padding: 8px 4px;
		display: grid;
		grid-template-columns: auto 1fr auto;
		align-items: center;
		gap: 8px;
	}

	li.candidate:not(:last-child) {
		border-bottom: 2px solid #eee;
	}

	li.candidate p {
		text-overflow: ellipsis;
		overflow: hidden;
	}

	li.candidate p > span {
		margin-left: 8px;
		background: #000;
		color: #fff;
		padding: 0 8px;
		border-radius: 4px;
		text-transform: uppercase;
		font-family: "JetBrains Mono";
		font-weight: bold;
	}

	div.admin-controls {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: 8px;
	}

	textarea {
		border: 2px solid;
		border-radius: 4px;
		padding: 2px 12px;
		width: calc(100% - 24px);
		min-height: 5rem;
		font-family: "Inter", sans-serif;
		resize: vertical;
		margin-bottom: 12px;
	}

	@media only screen and (max-width: 600px) {
		img.banner-image {
			display: none;
		}
	}
</style>
