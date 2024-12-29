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
</script>

<Panel title="Results" headerIcon="receipt_long">
	<div class="results">
		<p>
			This referendum was held on <b>{date}</b>. There were {ballotsCast} ballots cast.
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
							{referendumOutcome.votesAbstain}
						</td>
					</tr>
			</tbody>
		</table>
		{#if referendumOutcome.votesAbstain > 0}
			<p>
				There were {referendumOutcome.votesAbstain}	abstentions.
			</p>
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
</style>
