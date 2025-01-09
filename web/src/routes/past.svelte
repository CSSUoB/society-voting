<script lang="ts">
	import { page } from "$app/stores";
	import Panel from "$lib/panel.svelte";
	import { polls } from "../store";
	import Button from "$lib/button.svelte";
	import List from "$lib/list.svelte";
	import { _getPolls } from "./+layout";
	import { goto } from "$app/navigation";
	import { getFriendlyName } from "$lib/poll";
</script>

<Panel title="Archive" headerIcon="inventory_2">
	<List items={$polls?.filter((p) => p.isConcluded) ?? []} let:prop={poll}>
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
