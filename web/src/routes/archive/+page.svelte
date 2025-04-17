<script lang="ts">
	import { goto } from "$app/navigation";
	import Button from "$lib/button.svelte";
	import List from "$lib/list.svelte";
	import Panel from "$lib/panel.svelte";
	import { getFriendlyName } from "$lib/poll";
	import { polls } from "../../store";

	const dateFormat = new Intl.DateTimeFormat("en-GB", {
		dateStyle: "medium",
		timeStyle: "short",
	});
	$: archivedPolls =
		$polls?.filter((p) => p.isConcluded).sort((a, b) => b.date.getTime() - a.date.getTime()) ?? [];
</script>

<svelte:head>
	<title>Archive</title>
</svelte:head>
<Panel title="Archive" headerIcon="inventory_2">
	{#if archivedPolls.length === 0}
		<p>There are no archived polls to show</p>
	{:else}
		<p>All polls, including those from over 30 days ago</p>
		<List items={archivedPolls} let:prop={poll}>
			<li class="poll">
				<div class="poll-details">
					<p><b>{getFriendlyName(poll)}</b></p>
					<p>{dateFormat.format(poll.date)}</p>
				</div>
				<Button
					icon="receipt_long"
					text="View results"
					on:click={() => goto(`/results/${poll.id}`)}
				/>
			</li>
		</List>
	{/if}
</Panel>

<style>
	li.poll {
		padding: 8px 0;
		display: grid;
		grid-template-columns: 1fr auto;
		align-items: center;
		gap: 8px;
	}

	li.poll:not(:last-child) {
		border-bottom: 2px solid #eee;
	}

	div.poll-details {
		display: flex;
		justify-content: space-between;
	}

	@media only screen and (max-width: 850px) {
		div.poll-details {
			justify-content: flex-start;
			flex-direction: column;
		}
	}
</style>
