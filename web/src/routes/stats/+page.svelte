<script lang="ts">
	import { goto } from "$app/navigation";
	import Button from "$lib/button.svelte";
	import { API } from "$lib/endpoints";
	import Panel from "$lib/panel.svelte";
	import { error, fetching, currentElection } from "../../store";

	let electionRunning = true;
	let results = "";
	let numberOfVotes = 0;

	$: if (!$currentElection) {
		goto("/");
	}

	const eventSource = new EventSource(API.ADMIN_ELECTION_SSE, {
		withCredentials: true,
	});
	eventSource.addEventListener("vote-received", (e) => {
		numberOfVotes += parseInt(e.data, 10);
	});

	const endElection = async () => {
		$fetching = true;
		const response = await fetch(API.ADMIN_ELECTION_STOP, {
			method: "POST",
		});
		if (!response.ok) {
			$fetching = false;
			$error = new Error(await response.text());
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
			`${$currentElection?.election.roleName.toLowerCase().replaceAll(" ", "_")}_results.txt`,
		);

		element.style.display = "none";
		document.body.appendChild(element);

		element.click();

		document.body.removeChild(element);
	};
</script>

<svelte:head>
	<title>Vote for: {$currentElection?.election.roleName}</title>
</svelte:head>

<Panel title={`Electing: ${$currentElection?.election.roleName}`}>
	<p>{$currentElection?.election.description}</p>
</Panel>

{#if electionRunning}
	<Panel title="Admin actions" headerIcon="admin_panel_settings">
		<div class="container">
			<h3>{numberOfVotes} of {$currentElection?.numEligibleVoters} users have voted so far</h3>
			<p>
				The turnout so far is {(numberOfVotes * 100) /
					($currentElection?.numEligibleVoters ?? 100)}%
			</p>
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
