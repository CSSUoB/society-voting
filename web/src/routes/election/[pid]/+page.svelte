<script lang="ts">
	import Avatar from "$lib/avatar.svelte";
	import Banner from "$lib/banner.svelte";
	import Button from "$lib/button.svelte";
	import Dialog from "$lib/dialog.svelte";
	import List from "$lib/list.svelte";
	import Panel from "$lib/panel.svelte";

	import run from "$lib/assets/run_for_election.svg";
	import { currentPoll, polls, error, fetching, user } from "../../../store";
	import { goto } from "$app/navigation";
	import { API } from "$lib/endpoints";
	import { _getCurrentPoll, _getPolls } from "../../+layout";
	import PollHeader from "$lib/poll-header.svelte";
	import Textarea from "$lib/textarea.svelte";
	import { isElectionPoll } from "$lib/poll";

	export let data: { id: number };
	let buttonText = "Stand for election";
	$: poll = $polls?.find((e) => e.id === data.id);
	$: {
		if (!poll || !isElectionPoll(poll)) {
			goto("/", { replaceState: true });
		} else if (poll.isActive) {
			goto("/vote");
		} else {
			buttonText = `Stand for ${poll.election.roleName ?? "election"}`;
		}
	}

	const standOrWithdraw = async (id: number, stand: boolean) => {
		$fetching = true;
		const response = await fetch(API.ELECTION_STAND, {
			method: stand ? "POST" : "DELETE",
			body: JSON.stringify({ id }),
		});

		if (response.ok) {
			polls.set(await _getPolls());
		} else {
			$error = new Error(await response.text());
		}
		$fetching = false;
	};

	let floorCandidates: string;
	let startElectionDialog: HTMLDialogElement;
	const startElection = async (id: number) => {
		$fetching = true;
		const extraNames = floorCandidates
			.trim()
			.split("\n")
			.filter((x) => x)
			.map((x) => x.trim());
		const response = await fetch(API.ADMIN_ELECTION_START, {
			method: "POST",
			body: JSON.stringify({ id, extraNames }),
		});

		if (response.ok) {
			polls.set(await _getPolls());
			currentPoll.set(await _getCurrentPoll());
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
		const response = await fetch(API.ADMIN_POLL, {
			method: "DELETE",
			body: JSON.stringify({ id }),
		});

		if (response.ok) {
			polls.set(await _getPolls());
			goto("/");
		} else {
			$error = new Error(await response.text());
		}
		$fetching = false;
	};
</script>

{#if poll && isElectionPoll(poll)}
	<PollHeader poll={poll}></PollHeader>

	{#if !poll.candidates.some((c) => c.isMe)}
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
					on:click={() => standOrWithdraw(poll.id ?? -1, true)}
				/>
			</svelte:fragment>
		</Banner>
	{/if}

	<Panel title="Candidates" headerIcon="groups">
		<List items={poll.candidates} let:prop={candidate}>
			<li class="candidate">
				<Avatar name={candidate.name} />
				<p>
					{candidate.name}{#if candidate.isMe}
						<span><small>You</small></span>
					{/if}
				</p>
				{#if candidate.isMe}
					<Button text="Withdraw" on:click={() => standOrWithdraw(poll.id ?? -1, false)} />
				{/if}
			</li>
		</List>
		{#if poll.candidates.length === 0}
			<p>There are no candidates currently running in this election</p>
		{/if}
	</Panel>

	{#if $user.isAdmin}
		<Panel title="Admin stuff" headerIcon="admin_panel_settings">
			<div class="admin-controls">
				<h3>Candidates standing from the floor</h3>
				<Textarea
					bind:value={floorCandidates}
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
			on:submit={() => startElection(poll.id ?? -1)}
		>
			<svelte:fragment slot="actions">
				<Button text="Cancel" />
				<Button text="Start election" kind="emphasis" name="submit" />
			</svelte:fragment>
		</Dialog>
		<Dialog
			bind:dialog={deleteElectionDialog}
			title="Delete election?"
			on:submit={() => deleteElection(poll.id ?? -1)}
		>
			<p>This action cannot be undone.</p>
			<svelte:fragment slot="actions">
				<Button text="Cancel" />
				<Button text="Delete election" kind="danger" name="submit" />
			</svelte:fragment>
		</Dialog>
	{/if}
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

	@media only screen and (max-width: 600px) {
		img.banner-image {
			display: none;
		}
	}
</style>
