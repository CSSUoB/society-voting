<script lang="ts">
	import Button from "$lib/button.svelte";
	import Dialog from "$lib/dialog.svelte";
	import Panel from "$lib/panel.svelte";

	import {
		type ElectionOutcome,
		type ElectionOutcomeResult,
		elections,
		error,
		fetching,
		user,
	} from "../../../store";
	import { goto } from "$app/navigation";
	import { API } from "$lib/endpoints";
	import { _getCurrentElection, _getElections } from "../../+layout";

	export let data: { data: ElectionOutcome };

	$: electionOutcome = data.data;

	$: election = electionOutcome.election;
	$: resultsByRound = electionOutcome.results.reduce((res, cur) => {
		return {
			...res,
			[cur.round]: res[cur.round]
				? [...res[cur.round], cur].sort((a, b) =>
						a.voteCount != b.voteCount ? b.voteCount - a.voteCount : b.isRejected ? -1 : 1,
				  )
				: [cur],
		};
	}, {} as { [key: string]: Array<ElectionOutcomeResult> });
	$: winner = electionOutcome.results
		.filter((e) => e.isElected)
		.map((e) => e.name)
		.join(", ");
	$: date = new Date(electionOutcome.date).toUTCString();

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

	const publishElectionResults = async (id: number, published: boolean) => {
		$fetching = true;
		const response = await fetch(API.ADMIN_ELECTION_PUBLISH, {
			method: "POST",
			body: JSON.stringify({ id, published }),
		});

		if (response.ok) {
			if (electionOutcome) electionOutcome.isPublished = published;
		} else {
			$error = new Error(await response.text());
		}
		$fetching = false;
	};
</script>

<svelte:head>
	<title>Outcome of election for {election.roleName}</title>
</svelte:head>

<Panel title={`Outcome of election for ${election.roleName}`}>
	<p>{election.description}</p>
</Panel>

{#if $user.isAdmin && !electionOutcome.isPublished}
	<Panel kind="emphasis">
		<div class="unpublished-callout">
			<p>
				The results for <b>{election.roleName}</b> have not been published. While they remain unpublished,
				only administrators may view the results.
			</p>
			<div class="button">
				<Button
					kind="primary"
					text="Publish results"
					on:click={() => publishElectionResults(election.id, true)}
				/>
			</div>
		</div>
	</Panel>
{/if}

<Panel title="Results" headerIcon="receipt_long">
	<div class="results">
		<p>
			This election was held at <b>{date}</b>. There were {electionOutcome.ballots}
			ballots cast and {electionOutcome.rounds} rounds in total.
		</p>
		{#each Object.entries(resultsByRound) as [round, results]}
			<h3>Round {round}</h3>
			<table>
				<thead>
					<tr>
						<th>Candidate</th>
						<th>Votes</th>
						<th>Status</th>
					</tr>
				</thead>
				<tbody>
					{#each results as result}
						<tr>
							<td>{result.name}</td>
							<td>{result.voteCount}</td>
							<td>
								{#if result.isRejected}
									Rejected
								{:else if result.isElected}
									Elected
								{:else}
									Proceed
								{/if}
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
			{#if results.some((e) => e.isRejected)}
				<p>
					<b>
						{results
							.filter((e) => e.isRejected)
							.map((e) => e.name)
							.join(", ")}
					</b> was eliminated
				</p>
			{/if}
		{/each}
		<p>
			After {electionOutcome.rounds} rounds of voting, the winning candidate for the post of
			<b>{election.roleName}</b>
			is <b>{winner}</b>
		</p>
		<h2>About these results</h2>
		<p>
			This election was held using the <b
				><a href="https://en.wikipedia.org/wiki/Instant-runoff_voting"
					>instant-runoff voting system</a
				></b
			>.
		</p>
		<p>
			With instant-runoff voting, voters may rank candidates based on preference from first to last,
			the candidate who has a majority of first-preference votes is elected. If there is no clear
			majority, then the candidate(s) with the fewest number of votes is eliminated, and their votes
			are re-assigned according to their next preference. This process is continued until a
			candidate has a majority.
		</p>
		<p>
			Voters may also choose to vote for <b>Re-open nominations (RON)</b>, which is a mechanism for
			voters to not to elect any candidate to the post. If RON wins, then the post remains unfilled
			and the election process starts from the beginning.
		</p>
		<p>
			If an election is tied, then the outcome is decided by random chance, in accordance with Guild
			of Students policy.
		</p>
	</div>
</Panel>

{#if $user.isAdmin}
	<Panel title="Admin stuff" headerIcon="admin_panel_settings">
		<div class="admin-controls">
			<div class="control-group">
				{#if electionOutcome.isPublished}
					<Button
						text="Unpublish results"
						on:click={() => publishElectionResults(election.id, false)}
					/>
				{:else}
					<Button
						text="Publish results"
						on:click={() => publishElectionResults(election.id, true)}
					/>
				{/if}
				<Button
					text="Delete election"
					kind="danger"
					on:click={() => deleteElectionDialog.showModal()}
				/>
			</div>
		</div>
	</Panel>
	<Dialog
		bind:dialog={deleteElectionDialog}
		title="Delete election?"
		on:submit={() => deleteElection(election.id)}
	>
		<p>This will delete this election and the associated results. This action cannot be undone.</p>
		<svelte:fragment slot="actions">
			<Button text="Cancel" />
			<Button text="Delete election" kind="danger" name="submit" />
		</svelte:fragment>
	</Dialog>
{/if}

<style>
	div.unpublished-callout {
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		align-items: center;
		gap: 16px;
	}

	div.unpublished-callout > .button {
		flex-shrink: 0;
	}

	div.control-group {
		display: flex;
		flex-direction: row;
		align-items: flex-start;
		gap: 8px;
	}

	div.results,
	div.admin-controls {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: 8px;
	}

	p {
		margin: 8px 0;
	}

	table {
		background: white;
		font-size: 12pt;
		border-collapse: collapse;
		text-align: left;
		border-radius: 4px;
	}
	table thead th,
	table tfoot th {
		background: #eee;
	}
	table th,
	table td {
		padding: 0.5em;
		border: 1px solid #ddd;
	}
</style>
