<script lang="ts">
	import Button from "$lib/button.svelte";
	import Dialog from "$lib/dialog.svelte";
	import Panel from "$lib/panel.svelte";

	import { currentPoll, polls, error, fetching, user } from "../../../store";
	import { goto } from "$app/navigation";
	import { API } from "$lib/endpoints";
	import { _getCurrentPoll, _getPolls } from "../../+layout";
	import PollHeader from "$lib/poll-header.svelte";
	import { isReferendumPoll } from "$lib/poll";

	export let data: { id: number };
	$: poll = $polls?.find((e) => e.id === data.id);
	$: if (!poll || !isReferendumPoll(poll)) {
		goto("/", { replaceState: true });
	} else if (poll.isActive) {
		goto("/vote");
	}

	let startReferendumDialog: HTMLDialogElement;
	const startReferendum = async (id: number) => {
		$fetching = true;
		const response = await fetch(API.ADMIN_REFERENDUM_START, {
			method: "POST",
			body: JSON.stringify({ id }),
		});

		if (response.ok) {
			polls.set(await _getPolls());
			currentPoll.set(await _getCurrentPoll());
			await new Promise((r) => setTimeout(r, 200));
			goto("/stats");
		} else {
			$error = new Error(await response.text());
		}
		$fetching = false;
	};

	let deleteReferendumDialog: HTMLDialogElement;
	const deleteReferendum = async (id: number) => {
		$fetching = true;
		const response = await fetch(API.ADMIN_POLL, {
			method: "DELETE",
			body: JSON.stringify({ id }),
		});

		if (response.ok) {
			polls.set(await _getPolls());
			goto("/");
		} else {
			$error = new Error(await response.text());
		}
		$fetching = false;
	};
</script>

{#if poll}
	<PollHeader poll={poll}></PollHeader>

	{#if $user.isAdmin}
		<Panel title="Admin stuff" headerIcon="admin_panel_settings">
			<div class="admin-controls">
				<Button
					kind="primary"
					text="Start referendum"
					on:click={() => startReferendumDialog.showModal()}
				/>
				<Button text="Delete referendum" on:click={() => deleteReferendumDialog.showModal()} />
			</div>
		</Panel>
		<Dialog
			bind:dialog={startReferendumDialog}
			title="Start referendum?"
			on:submit={() => startReferendum(poll.id ?? -1)}
		>
			<svelte:fragment slot="actions">
				<Button text="Cancel" />
				<Button text="Start referendum" kind="emphasis" name="submit" />
			</svelte:fragment>
		</Dialog>
		<Dialog
			bind:dialog={deleteReferendumDialog}
			title="Delete referendum?"
			on:submit={() => deleteReferendum(poll.id ?? -1)}
		>
			<p>This action cannot be undone.</p>
			<svelte:fragment slot="actions">
				<Button text="Cancel" />
				<Button text="Delete referendum" kind="danger" name="submit" />
			</svelte:fragment>
		</Dialog>
	{/if}
{/if}

<style>
	div.admin-controls {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: 8px;
	}
</style>
