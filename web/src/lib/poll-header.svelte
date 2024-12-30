<script lang="ts">
	import { type Poll } from "../store";
	import Panel from "./panel.svelte";

	export let poll: Poll;
	export let prefix: string = "";
	
	$: title = "";
	$: description = "";

	$: if (poll.election) {
		title = `Election of ${poll.election.roleName}`
		description = poll.election.description
	} else if (poll.referendum) {
		title = poll.referendum.title
		description = poll.referendum.description
	}
</script>

<svelte:head>
	<title>{title}</title>
</svelte:head>

<Panel title="{prefix ? prefix + ": " : ""}{title}">
	<p>{description}</p>
</Panel>

{#if poll.referendum} 
<Panel kind="primary">
	<div class="question">
		<span class="material-symbols-rounded">help_outline</span>
		<p>{poll.referendum.question}</p>
	</div>
</Panel>
{/if}

<style>
	div.question {
		display: flex;
		flex-direction: row;		
		gap: 0.5rem;
		align-items: center;
	}
</style>