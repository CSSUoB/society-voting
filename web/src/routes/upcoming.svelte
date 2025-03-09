<script lang="ts">
	import { page } from "$app/stores";
	import Panel from "$lib/panel.svelte";
	import { user, polls, error } from "../store";
	import Button from "$lib/button.svelte";
	import Dialog from "$lib/dialog.svelte";
	import Input from "$lib/input.svelte";
	import { API } from "$lib/endpoints";
	import List from "$lib/list.svelte";
	import { _getPolls } from "./+layout";
	import { goto } from "$app/navigation";
	import { isElectionPoll } from "$lib/poll";
</script>

<Panel title="Upcoming" headerIcon="campaign">
		<List items={$polls?.filter((p) => !p.isActive && !p.isConcluded) ?? []} let:prop={poll}>
			{#if isElectionPoll(poll)}
			<li
				class={`poll ${
					$page.url.pathname === `/election/${poll.id}` ? "poll--selected" : ""
				}`}
				on:click={() => goto(`/election/${poll.id}`)}
			>
				<p>{poll.election.roleName}</p>
				<Button icon="arrow_forward" on:click={() => goto(`/election/${poll.id}`)} />
			</li>
			{:else}
			<li
				class={`poll ${
					$page.url.pathname === `/referendum/${poll.id}` ? "poll--selected" : ""
				}`}
				on:click={() => goto(`/referendum/${poll.id}`)}
			>
				<p>{poll.referendum.title}</p>
				<Button icon="arrow_forward" on:click={() => goto(`/referendum/${poll.id}`)} />
			</li>
			{/if}
		</List>
	{#if $user.isAdmin}
	<div class="create-poll">
		<Button
			kind="emphasis"
			text="Create poll"
			icon="add"
			on:click={() => goto("/create")}
		/>
	</div>
	{/if}
</Panel>

<style>
	li.poll {
		display: grid;
		grid-template-columns: 1fr auto;
		gap: 8px;
		align-items: center;
		padding: 8px 0;
		cursor: pointer;
		position: relative;
	}

	li.poll--selected::after {
		content: "";
		position: absolute;
		top: 0;
		left: -16px;
		height: 100%;
		width: calc(100% + 2 * 16px);
		background-color: rgba(55, 93, 182, 0.1);
	}

	li.poll:not(:last-child) {
		border-bottom: 1px solid #eee;
	}
	
	div.create-poll {
		display: flex;
		flex-direction: row;
		align-items: flex-start;
		gap: 8px;
	}
</style>
