<script lang="ts">
	import { page } from "$app/stores";
	import Panel from "$lib/panel.svelte";
	import { user, elections, error } from "../store";
	import Button from "$lib/button.svelte";
	import Dialog from "$lib/dialog.svelte";
	import Input from "$lib/input.svelte";
	import { API } from "$lib/endpoints";
	import List from "$lib/list.svelte";
	import { _getElections } from "./+layout";
	import { goto } from "$app/navigation";

	let dialog: HTMLDialogElement;
	const createElection = async (e: CustomEvent<any>) => {
		const response = await fetch(API.ADMIN_ELECTION, {
			method: "POST",
			body: JSON.stringify(e.detail),
		});

		if (response.ok) {
			elections.set(await _getElections());
			goto(`/election/${$elections?.slice(-1)[0].id}`);
		} else {
			$error = new Error(await response.text());
		}
	};
</script>

<Panel title="Upcoming elections" headerIcon="campaign">
	<List items={$elections?.filter((e) => !e.isActive && !e.isConcluded) ?? []} let:prop={election}>
		<li
			class={`election ${
				$page.url.pathname === `/election/${election.id}` ? "election--selected" : ""
			}`}
			on:click={() => goto(`/election/${election.id}`)}
		>
			<p>{election.roleName}</p>
			<Button icon="arrow_forward" on:click={() => goto(`/election/${election.id}`)} />
		</li>
	</List>
	{#if $user.isAdmin}
		<Button
			kind="emphasis"
			text="Create a new election"
			icon="add"
			on:click={() => dialog.showModal()}
		/>
	{/if}
</Panel>

<Dialog bind:dialog title="Create a new election" on:submit={createElection}>
	<p>Create election with name:</p>
	<Input name="roleName" />
	<p>And description:</p>
	<Input name="description" />
	<svelte:fragment slot="actions">
		<Button text="Cancel" />
		<Button text="Create election" kind="emphasis" name="submit" />
	</svelte:fragment>
</Dialog>

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
