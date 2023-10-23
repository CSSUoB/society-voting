<script lang="ts">
	import Panel from "$lib/panel.svelte";
	import { get } from "svelte/store";
	import { userStore } from "../store";
	import Button from "$lib/button.svelte";
	import Dialog from "$lib/dialog.svelte";
	import Input from "$lib/input.svelte";
	import { API } from "$lib/endpoints";

	const user = get(userStore);
	let dialog: HTMLDialogElement;

	const createElection = async (e: CustomEvent<any>) => {
		const response = await fetch(API.ADMIN_ELECTION, {
			method: "POST",
			body: JSON.stringify(e.detail),
		});

		if (response.ok) {
			alert("Created election!");
		}
	};
</script>

<Panel title="Upcoming elections" headerIcon="campaign">
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
