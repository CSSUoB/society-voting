<script lang="ts">
	import { page } from "$app/stores";
	import Panel from "$lib/panel.svelte";
	import { polls } from "../store";
	import Button from "$lib/button.svelte";
	import List from "$lib/list.svelte";
	import { _getPolls } from "./+layout";
	import { goto } from "$app/navigation";
	import { getFriendlyName } from "$lib/poll";

	// todo: see todo in archive page 
	$: pollsToShow = $polls?.filter((p) => p.isConcluded && (new Date()).getTime() - (new Date(p.date)).getTime() < 60*60*24*30) ?? []
</script>

<Panel title="Archive" headerIcon="inventory_2">
	<List items={pollsToShow} let:prop={poll}>
		<li
			class={`election ${
				$page.url.pathname === `/results/${poll.id}` ? "election--selected" : ""
			}`}
			on:click={() => goto(`/results/${poll.id}`)}
		>
			<p>{getFriendlyName(poll)}</p>
			<Button icon="arrow_forward" on:click={() => goto(`/results/${poll.id}`)} />
		</li>
	</List>
	{#if pollsToShow.length === 0}
		<p>There are no concluded polls within the past 30 days</p>
	{/if}
	<br />
	<Button text="View all" kind="emphasis" on:click={() => goto(`/archive`)} />
</Panel>

<style>
	li.election {
		display: grid;
		grid-template-columns: 1fr auto;
		gap: 8px;
		align-items: center;
		padding: 8px 0;
		cursor: pointer;
		position: relative;
	}

	li.election--selected::after {
		content: "";
		position: absolute;
		top: 0;
		left: -16px;
		height: 100%;
		width: calc(100% + 2 * 16px);
		background-color: rgba(55, 93, 182, 0.1);
	}

	li.election:not(:last-child) {
		border-bottom: 1px solid #eee;
	}
</style>
