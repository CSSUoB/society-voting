<script lang="ts">
	import Button from "$lib/button.svelte";
	import Dialog from "$lib/dialog.svelte";
	import Panel from "$lib/panel.svelte";

	import {
		polls,
		error,
		fetching,
		user,
		type PollOutcome,
	} from "../../../store";
	import { goto } from "$app/navigation";
	import { API } from "$lib/endpoints";
	import { _getCurrentPoll, _getPolls } from "../../+layout";
	import ElectionResult from "./election-result.svelte";
	import PollHeader from "$lib/poll-header.svelte";
	import ReferendumResult from "./referendum-result.svelte";
	import { getFriendlyName, isElectionPollOutcome } from "$lib/poll";

	export let data: { data: PollOutcome };

	$: pollOutcome = data.data;

	let deletePollDialog: HTMLDialogElement;
	const deletePoll = async (id: number) => {
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

	const publishPollResults = async (id: number, published: boolean) => {
		$fetching = true;
		const response = await fetch(API.ADMIN_POLL_PUBLISH, {
			method: "POST",
			body: JSON.stringify({ id, published }),
		});

		if (response.ok) {
			pollOutcome.isPublished = published;
		} else {
			$error = new Error(await response.text());
		}
		$fetching = false;
	};

</script>

<svelte:head>
	<title>Poll outcome: {getFriendlyName(pollOutcome.poll)}</title>
</svelte:head>

<PollHeader poll={pollOutcome.poll}></PollHeader>

{#if $user.isAdmin && !pollOutcome.isPublished}
	<Panel kind="emphasis">
		<div class="unpublished-callout">
			<p>
				The results for this poll have not been published. While they remain unpublished,
				only administrators may view the results.
			</p>
			<Button
				kind="primary"
				text="Publish results"
				on:click={() => publishPollResults(pollOutcome.poll.id, true)}
			/>
		</div>
	</Panel>
{/if}

{#if isElectionPollOutcome(pollOutcome)}
	<ElectionResult date={pollOutcome.date} ballotsCast={pollOutcome.ballots} electionName={pollOutcome.poll.election.roleName} electionOutcome={pollOutcome.electionOutcome} />
{:else}
	<ReferendumResult date={pollOutcome.date} ballotsCast={pollOutcome.ballots} referendumOutcome={pollOutcome.referendumOutcome} />
{/if}

{#if $user.isAdmin}
	<Panel title="Admin stuff" headerIcon="admin_panel_settings">
		<div class="admin-controls">
			<div class="control-group">
				{#if pollOutcome.isPublished}
					<Button
						text="Unpublish results"
						on:click={() => publishPollResults(pollOutcome.poll.id, false)}
					/>
				{:else}
					<Button
						text="Publish results"
						on:click={() => publishPollResults(pollOutcome.poll.id, true)}
					/>
				{/if}
				<Button
					text="Delete poll"
					kind="danger"
					on:click={() => deletePollDialog.showModal()}
				/>
			</div>
		</div>
	</Panel>
	<Dialog
		bind:dialog={deletePollDialog}
		title="Delete poll?"
		on:submit={() => deletePoll(pollOutcome.poll.id)}
	>
		<p>This will delete this poll and the associated results. This action cannot be undone.</p>
		<svelte:fragment slot="actions">
			<Button text="Cancel" />
			<Button text="Delete poll" kind="danger" name="submit" />
		</svelte:fragment>
	</Dialog>
{/if}

<style>
	div.unpublished-callout {
		display: grid;
		grid-template-columns: 1fr auto;
		align-items: center;
		gap: 16px;
	}

	div.control-group {
		display: flex;
		flex-direction: row;
		align-items: flex-start;
		gap: 8px;
	}

	div.admin-controls {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: 16px;
	}
</style>
