<script lang="ts">
	import { type Poll } from "../store";
	import Panel from "./panel.svelte";
	import { isElectionPoll, isReferendumPoll } from "./poll";

	export let poll: Poll;
	export let prefix: string = "";
	
	$: title = "";
	$: description = "";

	$: if (isElectionPoll(poll)) {
		title = `Election of ${poll.election.roleName}`
		description = poll.election.description
	} else {
		title = `Referendum on ${poll.referendum.title}`
		description = poll.referendum.description
	}
</script>

<svelte:head>
	<title>{title}</title>
</svelte:head>

<Panel title="{prefix ? prefix + ": " : ""}{title}">
	<p>{description}</p>
</Panel>

{#if isReferendumPoll(poll)} 
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