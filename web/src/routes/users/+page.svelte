<script lang="ts">
	import Button from "$lib/button.svelte";
	import Dialog from "$lib/dialog.svelte";
	import { API } from "$lib/endpoints";
	import List from "$lib/list.svelte";
	import type { Optional } from "$lib/optional";
	import Panel from "$lib/panel.svelte";
	import { error, fetching, type User } from "../../store";

	export let data: { users: Array<User> };

	let userToModify: Optional<User>;
	let removeUserDialog: HTMLDialogElement;
	const confirmRemoveUser = (user: User) => {
		userToModify = user;
		removeUserDialog.showModal();
	};
	const removeUser = async (userID: string) => {
		$fetching = true;
		const response = await fetch(API.ADMIN_USER_DELETE, {
			method: "DELETE",
			body: JSON.stringify({ userID }),
		});

		if (response.ok) {
			data = { ...data, users: data.users.filter((u) => u.studentID !== userID) };
		} else {
			$error = new Error(await response.text());
		}
		$fetching = false;
	};

	let restrictUserDialog: HTMLDialogElement;
	const confirmRestrictUser = (user: User) => {
		if (user.isRestricted) return toggleUserRestriction(user.studentID, user.isRestricted);
		userToModify = user;
		restrictUserDialog.showModal();
	};
	const toggleUserRestriction = async (userID: string, isCurrentlyRestricted: boolean) => {
		$fetching = true;
		const response = await fetch(API.ADMIN_USER_RESTRICT, {
			method: "POST",
			body: JSON.stringify({ userID }),
		});

		const j = await response.json();

		if (response.ok) {
			data = {
				...data,
				users: data.users.map((u) => {
					if (u.studentID === userID) {
						u.isRestricted = j.isRestricted;
					}
					return u;
				}),
			};
		} else {
			$error = new Error(await response.text());
		}
		$fetching = false;
	};
</script>

<svelte:head>
	<title>Users</title>
</svelte:head>

<Panel title="Manage users" headerIcon="admin_panel_settings">
	<Button slot="header-action" icon="search" kind="emphasis" text="Search users" />
	<List items={data.users} let:prop={user}>
		<li class="user" class:restricted={user.isRestricted}>
			<p>{user.studentID}</p>
			<p>
				{user.name}{#if user.isRestricted} <span class="pill"><small>Restricted</small></span>{/if}
			</p>
			<Button
				icon={user.isRestricted ? "check" : "block"}
				text={user.isRestricted ? "Unrestrict user" : "Restrict user"}
				on:click={confirmRestrictUser.bind(null, user)}
			/>
			<Button
				icon="person_remove"
				text="Delete user"
				on:click={confirmRemoveUser.bind(null, user)}
			/>
		</li>
	</List>
</Panel>

<Dialog
	bind:dialog={removeUserDialog}
	title={`Delete "${userToModify?.name}"?`}
	on:submit={() => userToModify && removeUser(userToModify?.studentID)}
>
	<p>
		Once deleted, they will have to sign up again and re-run for any elections they are currently
		contesting.
	</p>
	<svelte:fragment slot="actions">
		<Button text="Cancel" />
		<Button text="Delete user" kind="emphasis" name="submit" />
	</svelte:fragment>
</Dialog>

<Dialog
	bind:dialog={restrictUserDialog}
	title={`Restrict "${userToModify?.name}"?`}
	on:submit={() =>
		userToModify && toggleUserRestriction(userToModify?.studentID, userToModify.isRestricted)}
>
	<p>
		This will remove the user from all elections they are currently standing in and make them unable
		to stand for any other elections or change their name.
		<br />
		Are you sure this is what you want to do?
	</p>
	<svelte:fragment slot="actions">
		<Button text="Cancel" />
		<Button text="Restrict user" kind="emphasis" name="submit" />
	</svelte:fragment>
</Dialog>

<style>
	li.user {
		padding: 8px 4px;
		display: grid;
		grid-template-columns: 150px 1fr auto auto;
		align-items: center;
		gap: 8px;
	}

	li.user:not(:last-child) {
		border-bottom: 2px solid #eee;
	}

	li.user p {
		text-overflow: ellipsis;
		overflow: hidden;
	}

	li.user.restricted {
		background-color: rgba(255, 0, 0, 0.2);
	}

	li.user p > span.pill {
		margin-left: 8px;
		background: rgba(255, 0, 0, 0.5);
		color: #fff;
		padding: 0 8px;
		border-radius: 4px;
		text-transform: uppercase;
		font-family: "JetBrains Mono";
		font-weight: bold;
	}
</style>
