<script lang="ts">
	import Panel from "$lib/panel.svelte";
	import { get } from "svelte/store";
	import { userStore, type Election, electionStore } from "../store";
	import Button from "$lib/button.svelte";
	import Dialog from "$lib/dialog.svelte";
	import Input from "$lib/input.svelte";
	import { API } from "$lib/endpoints";
	import { onDestroy } from "svelte";
	import List from "$lib/list.svelte";
	import { _getElections } from "./+layout";

	const user = get(userStore);
	let dialog: HTMLDialogElement;

	let elections: Array<Election>;

	const unsubscribe = electionStore.subscribe((e) => (elections = e));
	onDestroy(unsubscribe);

	const createElection = async (e: CustomEvent<any>) => {
		const response = await fetch(API.ADMIN_ELECTION, {
			method: "POST",
			body: JSON.stringify(e.detail),
		});

		if (response.ok) {
			electionStore.set(await _getElections());
		}
	};
</script>

<Panel title="Upcoming elections" headerIcon="campaign">
	<List items={elections.filter((e) => !e.isActive)} let:prop={election}>
		<li class="election">
			<p>{election.roleName}</p>
			<Button icon="arrow_forward" />
		</li>
	</List>
	{#if user.admin}
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
	}

	li.election:not(:last-child) {
		border-bottom: 1px solid #eee;
	}
</style>
