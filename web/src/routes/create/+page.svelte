<script lang="ts">
	import { goto } from "$app/navigation";
	import Button from "$lib/button.svelte";
	import Checkbox from "$lib/checkbox.svelte";
	import { getEndpointForPollType, PollTypeId } from "$lib/endpoints";
	import List from "$lib/list.svelte";
	import Panel from "$lib/panel.svelte";

	import { _getCurrentPoll, _getPolls } from "../+layout";
	import { error, fetching, polls } from "../../store";
	import CreateElection from "./create-election.svelte";
	import CreateReferendum from "./create-referendum.svelte";

	const options = [{
		label: "Election",
		description: "Hold an election for a post, using instant-runoff voting.",
		symbol: "people",
		pollTypeId: PollTypeId.ELECTION
	}, {
		label: "Referendum",
		description: "Poll members on a specific question. Members may vote for, against, or abstain.",
		symbol: "rule",
		pollTypeId: PollTypeId.REFERENDUM
	}]
	
	let selectedIndex: number | null = null;
	let stayOnPage: boolean = false;
	let lastPollId: number | null = null;

	const selectOption = (option: number) => {
		lastPollId = null;
		if (selectedIndex == option) {
			selectedIndex = null;		
		} else {
			selectedIndex = option;
		}
	}
	
	const submit = async (e: SubmitEvent) => {
		if ((e.submitter as HTMLButtonElement).name !== "submit" || null === selectedIndex) return;
		e.preventDefault();
		const target = e.target as HTMLFormElement;
		const formData = new FormData(target);
		const option = options[selectedIndex];

		const url = getEndpointForPollType("create", option.pollTypeId);
		if (!url) {
			$error = new Error("Unknown poll type");
			return;
		}
		$fetching = true;

		const response = await fetch(url, {
			method: "POST",
			body: JSON.stringify(Object.fromEntries(formData)),
		});
		
		if (response.ok) {
			const json = await response.json();
			const id = json.id;

			$polls = await _getPolls();
			$fetching = false;

			if (stayOnPage) {
				lastPollId = id;
				target.reset();
				(target.querySelectorAll("input[type=checkbox]")[0] as HTMLInputElement).click(); // re-enable stay on page
			} else {
				goto(`${option.label.toLowerCase()}/${id}`);
			}
		} else {
			$error = new Error(await response.text());
			$fetching = false;
			return;
		}
	};
</script>

<svelte:head>
	<title>Create a new poll</title>
</svelte:head>

<form class="rail" on:submit={submit}>
	<Panel title="Create a new poll" headerIcon="post_add">
		<div class="container">
			<p>Select a poll type</p>
			<List items={options} let:prop={option}>
				<li class="option {option.index == selectedIndex ? 'selected' : null == selectedIndex ? '' : 'not-selected'}" on:click={() => selectOption(option.index)}>
					{#if option.index === selectedIndex}
						<span class="check-mark material-symbols-rounded">check_circle</span>
					{:else}
						<span />
					{/if}
					<span class="option-icon material-symbols-rounded">{ option.symbol }</span>
					<span />
					<div class="option-text">
						<span class="option-label">{ option.label }</span>
						<span class="option-description">{ option.description }</span>
					</div>
				</li>
			</List>
		</div>
	</Panel>

	{#if null !== lastPollId}
		<Panel kind="primary">
			<p>A poll with ID {lastPollId} was successfully created.</p>
		</Panel>
	{/if}

	{#if null !== selectedIndex}
		{#if 0 === selectedIndex}
			<CreateElection />
		{:else if 1 === selectedIndex}
			<CreateReferendum />
		{/if}

		<Panel>
			<div class="controls">
				<Button text="Create" kind="primary" name="submit" />
				<Checkbox label="Keep this page open" bind:value={stayOnPage} />
			</div>
		</Panel>
	{/if}
</form>

<style>
	div.container {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	li.option {
		padding: 8px 4px;
		display: grid;
		grid-template-columns: 8px 40px 8px auto;
		align-items: center;
		cursor: pointer;
		min-height: 32px;
		user-select: none;
		transition: grid-template-columns 0.1s, margin-left 0.2s, opacity 0.2s;
	}
	
	div.option-text {
		display: flex;
		flex-direction: column;
	}

	span.option-description {
		font-family: "Inter", sans-serif;
	}

	span.option-label {
		font-family: "JetBrains Mono", monospace;
		font-weight: bolder;
		position: relative;
	}

	li.option.selected {
		grid-template-columns: 28px 40px 8px auto;
		background-color: rgba(55, 93, 182, 0.1);
	}

	li.option.not-selected {
		grid-template-columns: 5px 40px 8px auto;
		opacity: 0.6;
	}

	li.option:hover {
		background-color: rgba(55, 93, 182, 0.1);
	}

	li.option:not(:last-child) {
		border-bottom: 2px solid #eee;
	}

	span.option-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		color: #fff;
		background-color: #1c2e58;
		border-radius: 999em;
		height: 40px;
	}

	span.check-mark {
		color: #1c2e58;
	}

	div.controls {
		display: flex;
		gap: 0.5rem;
	}
</style>
