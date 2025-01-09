<script lang="ts">
	import { goto } from "$app/navigation";
	import Button from "$lib/button.svelte";
	import Panel from "$lib/panel.svelte";
	import { getFriendlyName } from "$lib/poll";
	import { currentPoll, polls, user } from "../store";
	$: current = $polls?.find((p) => p.isActive);
</script>

<Panel title="Vote now!" kind="emphasis" headerIcon="how_to_vote">
	<div class="container">
		<h3>Voting is now open for: {current && getFriendlyName(current)}</h3>
		<div class="actions">
			{#if !$currentPoll?.hasVoted}
				<Button
					icon="arrow_forward"
					text="Vote now"
					kind="primary"
					on:click={() => goto(`/vote`)}
				/>
			{/if}
			{#if $user.isAdmin}
				<Button text="Manage poll" kind="emphasis" on:click={() => goto(`/stats`)} />
			{/if}
		</div>
	</div>
</Panel>

<style>
	div.container,
	div.actions {
		display: flex;
		flex-direction: column;
		gap: 8px;
		align-items: flex-start;
	}

	div.actions {
		flex-direction: row;
	}
</style>
