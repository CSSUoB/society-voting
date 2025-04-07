<script lang="ts">
	import { goto } from "$app/navigation";
	import Button from "$lib/button.svelte";
	import List from "$lib/list.svelte";
	import Panel from "$lib/panel.svelte";
	import { getFriendlyName } from "$lib/poll";
	import { polls } from "../../store";

	$: archivedPolls = $polls?.filter((p) => p.isConcluded);
</script>

<svelte:head>
	<title>Users</title>
</svelte:head>
<Panel title="Archive" headerIcon="inventory_2">
    <p>
        All polls, including those from over 30 days ago
    </p>
	<List items={archivedPolls} let:prop={poll}>
		<li class="poll">
			<p>{getFriendlyName(poll)}</p>
            <p>{(new Date(poll.date)).toUTCString()}</p> <!-- todo: probably a better idea to have the right type in the first place -->
			<Button
				icon="receipt_long"
				text="View results"
                on:click={() => goto(`/results/${poll.id}`)}
			/>
		</li>
	</List>
</Panel>

<style>
	li.poll {
		padding: 8px 4px;
		display: grid;
		grid-template-columns: 1fr auto auto;
		align-items: center;
		gap: 8px;
	}

	li.poll:not(:last-child) {
		border-bottom: 2px solid #eee;
	}

	li.poll p {
		text-overflow: ellipsis;
		overflow: hidden;
	}
</style>
