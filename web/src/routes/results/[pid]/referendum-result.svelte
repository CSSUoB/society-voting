<script lang="ts">
	import { type ReferendumOutcome } from "../../../store";
	import Panel from "../../../lib/panel.svelte";

	export let date: string;
	export let ballotsCast: number;
	export let referendumOutcome: ReferendumOutcome;

	$: date = new Intl.DateTimeFormat("en-GB", {
		dateStyle: "full",
		timeStyle: "long",
	}).format(new Date(date));
	$: passed = referendumOutcome.votesFor > referendumOutcome.votesAgainst;
</script>

<Panel title="Results" headerIcon="receipt_long">
	<div class="results">
		<p>
			This referendum was held on <b>{date}</b>. 
			{#if ballotsCast === 1}
				There was 1 ballot cast.
			{:else}
				There were {ballotsCast} ballots cast.
			{/if}
		</p>
		<table>
			<thead>
				<tr>
					<th>Option</th>
					<th>Votes</th>
				</tr>
			</thead>
			<tbody>
					<tr>
						<td>For</td>
						<td>
							{referendumOutcome.votesFor}
						</td>
					</tr>
					<tr>
						<td>Against</td>
						<td>
							{referendumOutcome.votesAgainst}
						</td>
					</tr>
			</tbody>
		</table>
		{#if referendumOutcome.votesAbstain > 0}
			<p>
				{#if referendumOutcome.votesAbstain === 1}
					There was 1 abstention.
				{:else}
					There were {referendumOutcome.votesAbstain}	abstentions.
				{/if}
			</p>
		{/if}
		{#if passed}
			<div class="interpretation">
				<span class="interpretation-icon material-symbols-rounded">thumb_up</span>
				<span class="interpretation-text">
					As the number of votes in support of this referendum exceeds the number 
					of votes in opposition, this referendum has <b>passed</b>.
				</span>
			</div>
		{:else}
			<div class="interpretation">
				<span class="interpretation-icon material-symbols-rounded">thumb_down</span>
				<span class="interpretation-text">
					As the number of votes in support of this referendum does not exceed the number 
					of votes in opposition, this referendum has <b>failed</b>.
				</span>
			</div>
		{/if}
	</div>
</Panel>

<style>
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

	div.interpretation {
		display: grid;
		grid-template-columns: 40px auto;
		gap: 8px;
		align-items: flex-start;
		min-height: 32px;
		width: 100%;
	}

	span.interpretation-text {
		font-family: "Inter", sans-serif;
		align-self: center;
	}

	span.interpretation-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		color: #000;
		background-color: #ddd;
		border-radius: 999em;
		height: 40px;
	}
</style>
