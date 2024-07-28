<script lang="ts">
	import { page } from "$app/stores";
	import Panel from "$lib/panel.svelte";
	import { elections } from "../store";
	import Button from "$lib/button.svelte";
	import List from "$lib/list.svelte";
	import { _getElections } from "./+layout";
	import { goto } from "$app/navigation";
</script>

<Panel title="Past elections" headerIcon="inventory_2">
	<List items={$elections?.filter((e) => e.isConcluded) ?? []} let:prop={election}>
		<li
			class={`election ${
				$page.url.pathname === `/results/${election.id}` ? "election--selected" : ""
			}`}
			on:click={() => goto(`/results/${election.id}`)}
		>
			<p>{election.roleName}</p>
			<Button icon="arrow_forward" on:click={() => goto(`/results/${election.id}`)} />
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
