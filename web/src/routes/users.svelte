<script lang="ts">
	import { goto } from "$app/navigation";
	import Button from "$lib/button.svelte";
	import Panel from "$lib/panel.svelte";
</script>

<Panel title="Users" headerIcon="group">
	<p>Administrator user management tools</p>
	<br />
	<Button kind="emphasis" text="Manage users" on:click={() => goto("/users")} />
</Panel>
