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

	let electionDialog: HTMLDialogElement;
	let referendumDialog: HTMLDialogElement;
	const createElection = async (e: CustomEvent<any>) => {
		const response = await fetch(API.ADMIN_ELECTION, {
			method: "POST",
			body: JSON.stringify(e.detail),
		});

		if (response.ok) {
			polls.set(await _getPolls());
			goto(`/election/${$polls?.slice(-1)[0].id}`);
		} else {
			$error = new Error(await response.text());
		}
	};

	const createReferendum = async (e: CustomEvent<any>) => {
		const response = await fetch(API.ADMIN_REFERENDUM, {
			method: "POST",
			body: JSON.stringify(e.detail),
		});

		if (response.ok) {
			polls.set(await _getPolls());
			goto(`/referendum/${$polls?.slice(-1)[0].id}`);
		} else {
			$error = new Error(await response.text());
		}
	};
</script>

<Panel title="Upcoming" headerIcon="campaign">
		<List items={$polls?.filter((p) => !p.isActive && !p.isConcluded) ?? []} let:prop={poll}>
			{#if poll.election}
			<li
				class={`poll ${
					$page.url.pathname === `/election/${poll.id}` ? "poll--selected" : ""
				}`}
				on:click={() => goto(`/election/${poll.id}`)}
			>
				<p>{poll.election.roleName}</p>
				<Button icon="arrow_forward" on:click={() => goto(`/election/${poll.id}`)} />
			</li>
			{:else if poll.referendum}
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
			text="New election"
			icon="add"
			on:click={() => electionDialog.showModal()}
		/>

		<Button
			kind="emphasis"
			text="New referendum"
			icon="add"
			on:click={() => referendumDialog.showModal()}
		/>
	</div>
	{/if}
</Panel>

<Dialog bind:dialog={electionDialog} title="Create a new election" on:submit={createElection}>
	<p>Create election with name:</p>
	<Input name="roleName" />
	<p>And description:</p>
	<Input name="description" />
	<svelte:fragment slot="actions">
		<Button text="Cancel" />
		<Button text="Create election" kind="emphasis" name="submit" />
	</svelte:fragment>
</Dialog>

<Dialog bind:dialog={referendumDialog} title="Create a new referendum" on:submit={createReferendum}>
	<p>Create referendum with title:</p>
	<Input name="title" />
	<p>And question:</p>
	<Input name="question" />
	<p>And description:</p>
	<Input name="description" />
	<svelte:fragment slot="actions">
		<Button text="Cancel" />
		<Button text="Create referendum" kind="emphasis" name="submit" />
	</svelte:fragment>
</Dialog>

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
