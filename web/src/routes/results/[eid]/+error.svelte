<script lang="ts">
	import Panel from "$lib/panel.svelte";
	import { page } from "$app/stores";
</script>

<svelte:head>
	<title>Error</title>
</svelte:head>

<Panel title="Cannot view election outcome">
	<p>{$page.error?.message} ({$page.status})</p>
</Panel>
