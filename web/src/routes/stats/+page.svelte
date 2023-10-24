<script lang="ts">
	import Button from "$lib/button.svelte";
	import { API } from "$lib/endpoints";
	import Panel from "$lib/panel.svelte";
	import { fetching, type CurrentElection } from "../../store";

	export let data: CurrentElection;
	let electionRunning = true;
	let results = "";

	const endElection = async () => {
		$fetching = true;
		const response = await fetch(API.ADMIN_ELECTION_STOP, {
			method: "POST",
		});
		if (!response.ok) {
			return;
		}
		electionRunning = false;
		results = (await response.json()).result;
		$fetching = false;
	};

	const downloadResults = (results: string) => {
		const element = document.createElement("a");
		element.setAttribute("href", "data:text/plain;charset=utf-8," + encodeURIComponent(results));
		element.setAttribute(
			"download",
			`${data.election.roleName.toLowerCase().replaceAll(" ", "_")}_results.txt`,
		);

		element.style.display = "none";
		document.body.appendChild(element);

		element.click();

		document.body.removeChild(element);
	};
</script>

<svelte:head>
	<title>Vote for: {data.election.roleName}</title>
</svelte:head>

<Panel title={`Electing: ${data.election.roleName}`}>
	<p>{data.election.description}</p>
</Panel>

{#if electionRunning}
	<Panel title="Admin actions" headerIcon="admin_panel_settings">
		<div class="container">
			<h3>{"n"} users have voted so far</h3>
			<Button text="End election and view results" kind="primary" on:click={endElection} />
		</div>
	</Panel>
{:else}
	<Panel title="Results" headerIcon="receipt_long">
		<div class="container">
			<p>{results}</p>
			<Button text="Download results" kind="primary" on:click={() => downloadResults(results)} />
		</div>
	</Panel>
{/if}

<style>
	div.container {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: 8px;
	}

	p {
		white-space: break-spaces;
	}
</style>
