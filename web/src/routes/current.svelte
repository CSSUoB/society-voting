<script lang="ts">
	import { goto } from "$app/navigation";
	import Button from "$lib/button.svelte";
	import Panel from "$lib/panel.svelte";
	import { elections } from "../store";
	$: current = $elections?.find((e) => e.isActive);
</script>

<Panel title="Vote now!" kind="emphasis" headerIcon="how_to_vote">
	<div>
		<h3>Voting is now open for: {current?.roleName}</h3>
		<Button icon="arrow_forward" text="Vote now" kind="primary" on:click={() => goto(`/vote`)} />
	</div>
</Panel>

<style>
	div {
		display: flex;
		flex-direction: column;
		gap: 8px;
		align-items: flex-start;
	}
</style>
