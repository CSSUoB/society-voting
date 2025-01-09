<script lang="ts">
	import { type ElectionOutcome, type ElectionOutcomeResult, user } from "../../../store";
	import Button from "../../../lib/button.svelte";
	import Panel from "../../../lib/panel.svelte";

	export let date: string;
	export let electionName: string;
	export let ballotsCast: number;
	export let electionOutcome: ElectionOutcome;

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
	$: date = new Intl.DateTimeFormat("en-GB", {
		dateStyle: "full",
		timeStyle: "long",
	}).format(new Date(date));
	$: isVotesShown = $user.isAdmin;

	const toggleVoteVisibility = () => (isVotesShown = !isVotesShown);
</script>

<Panel title="Results" headerIcon="receipt_long">
	<div slot="header-action">
		{#if isVotesShown}
			<Button
				icon="visibility_off"
				kind="emphasis"
				text="Hide votes"
				on:click={toggleVoteVisibility}
			/>
		{:else}
			<Button
				icon="visibility"
				kind="emphasis"
				text="Reveal votes"
				on:click={toggleVoteVisibility}
			/>
		{/if}
	</div>
	<div class="results">
		<p>
			This election was held on <b>{date}</b>. 
			{#if ballotsCast === 1}
				There was 1 ballot cast
			{:else}
				There were {ballotsCast} ballots cast
			{/if}
			and {electionOutcome.rounds} rounds in total.
		</p>
		{#each Object.entries(resultsByRound) as [round, results]}
			<h3>Round {round}</h3>
			<table>
				<thead>
					<tr>
						<th>Candidate</th>
						<th>Votes</th>
						<th>Outcome</th>
					</tr>
				</thead>
				<tbody>
					{#each results as result}
						<tr>
							<td>{result.name}</td>
							{#if isVotesShown}
								<td>
									{result.voteCount}
								</td>
							{:else}
								<td
									class="hidden-icon"
									on:click={toggleVoteVisibility}
									title="Click to reveal votes"
								>
									<span class="material-symbols-rounded">visibility_off</span>
								</td>
							{/if}
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
			<b>{electionName}</b>
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

<style>
	span.material-symbols-rounded {
		font-size: 30px;
	}
	
	div.results {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: 16px;
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
	
	td.hidden-icon {
		color: #555;
		cursor: pointer;
	}

	td.hidden-icon > span {
		font-size: 0.9em;
	}
</style>
