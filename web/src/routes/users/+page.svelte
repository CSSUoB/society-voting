<script lang="ts">
	import Button from "$lib/button.svelte";
	import List from "$lib/list.svelte";
	import Panel from "$lib/panel.svelte";

	let users = [{ name: "", studentId: 0 }];
</script>

<Panel title="Manage users" headerIcon="admin_panel_settings">
	<Button slot="header-action" icon="search" kind="emphasis" />
	<List items={users} let:prop={user}>
		<li class="candidate">
			<p>{user.studentId}</p>
			<p>{user.name}</p>
			<Button icon="person_remove" text="Delete user" />
		</li>
	</List>
</Panel>

<style>
	li.candidate {
		padding: 8px 4px;
		display: grid;
		grid-template-columns: 150px 1fr auto;
		align-items: center;
		gap: 8px;
	}

	li.candidate:not(:last-child) {
		border-bottom: 2px solid #eee;
	}

	li.candidate p {
		text-overflow: ellipsis;
		overflow: hidden;
	}
</style>
