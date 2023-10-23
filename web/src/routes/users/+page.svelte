<script lang="ts">
	import Button from "$lib/button.svelte";
	import { API } from "$lib/endpoints";
	import List from "$lib/list.svelte";
	import Panel from "$lib/panel.svelte";
	import type { User } from "../../store";

	export let data: { users: Array<User> };

	const removeUser = async (userID: string) => {
		const response = await fetch(API.ADMIN_USER_DELETE, {
			method: "DELETE",
			body: JSON.stringify({ userID }),
		});

		if (response.ok) {
			data = { ...data, users: data.users.filter((u) => u.studentID !== userID) };
		}
	};
</script>

<svelte:head>
	<title>Users</title>
</svelte:head>

<Panel title="Manage users" headerIcon="admin_panel_settings">
	<Button slot="header-action" icon="search" kind="emphasis" text="Search users" />
	<List items={data.users} let:prop={user}>
		<li class="user">
			<p>{user.studentID}</p>
			<p>{user.name}</p>
			<Button
				icon="person_remove"
				text="Delete user"
				on:click={removeUser.bind(null, user.studentID)}
			/>
		</li>
	</List>
</Panel>

<style>
	li.user {
		padding: 8px 4px;
		display: grid;
		grid-template-columns: 150px 1fr auto;
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
</style>
