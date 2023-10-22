<script lang="ts">
	import Button from "$lib/button.svelte";
	import Dialog from "$lib/dialog.svelte";
	import { API } from "$lib/endpoints";
	import Input from "$lib/input.svelte";
	import Panel from "$lib/panel.svelte";
	import { userStore, type User } from "../store";
	import { onDestroy } from "svelte";

	let user: User;

	const unsubscribe = userStore.subscribe((u) => (user = u));
	onDestroy(unsubscribe);

	let dialog: HTMLDialogElement;
	const updateName = async (e: CustomEvent<any>) => {
		const response = await fetch(API.ME_NAME, {
			method: "PUT",
			body: JSON.stringify(e.detail),
		});

		if (response.ok) {
			$: userStore.set({ ...user, name: e.detail.name });
		}
	};
</script>

<Panel title="Hi {user.name.split(' ')[0]}!" headerIcon="waving_hand">
	<Button
		slot="header-action"
		text="Use a different name"
		kind="emphasis"
		on:click={() => dialog.showModal()}
	/>
	<p>Welcome to CSS' voting system! View, stand for, and vote in currently running elections.</p>
	<br />

	<span class="log-out">
		<p>Not you?</p>
		<Button text="Log out" kind="inline" />
	</span>
</Panel>
<Dialog bind:dialog title="Use a different name" on:submit={updateName}>
	<p>Enter a new name to use</p>
	<Input value={user.name} name="name" />
	<svelte:fragment slot="actions">
		<Button text="Cancel" />
		<Button text="Update name" kind="emphasis" name="submit" />
	</svelte:fragment>
</Dialog>

<style>
	span.log-out {
		display: flex;
		align-items: center;
		gap: 8px;
	}
</style>
