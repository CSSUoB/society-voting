<script lang="ts">
	import BallotEntry from "./ballot-entry.svelte";
	import { type BallotEntry as BallotEntryT, fetching, user } from "../../store";
	import List from "../../lib/list.svelte";
	import Panel from "../../lib/panel.svelte";
	import { createEventDispatcher } from "svelte";
	
	const options = [{
		label: "For",
		symbol: "thumb_up",
		value: 1
	}, {
		label: "Against",
		symbol: "thumb_down",
		value: 2
	}, {
		label: "Abstain",
		symbol: "back_hand",
		value: 0
	}] as { label: string, symbol: string, value: number }[]

	let selectedIndex: number | null;

	const dispatch = createEventDispatcher();
	
	const selectOption = (option: number) => {
		if (selectedIndex == option) {
			selectedIndex = null;		
		} else {
			selectedIndex = option;
		}
		
		dispatch("update", {
			valid: null != selectedIndex,
			choices: null != selectedIndex ? [options[selectedIndex].value] : []
		})
	}
</script>

<Panel title="Your ballot">
	<div class="ballot">
		<p>You have one vote. Select one of the three options on your ballot.</p>

		<List items={options} let:prop={option}>
			<li class="option {option.index == selectedIndex ? 'selected' : null == selectedIndex ? '' : 'not-selected'}" on:click={() => selectOption(option.index)}>
				{#if option.index === selectedIndex}
					<span class="ballot--valid material-symbols-rounded">check_circle</span>
				{:else}
					<span />
				{/if}
				<span class="option-icon material-symbols-rounded">{ option.symbol }</span>
				<span />
				<span class="option-label">{ option.label }</span>
				<span />
			</li>
		</List>
		
		{#if null != selectedIndex}
			<p>You have selected: <b>{options[selectedIndex].label}</b></p>	
		{/if}
	</div>
</Panel>

<style>
	div.ballot {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	li.option {
		padding: 8px 4px;
		display: grid;
		grid-template-columns: 8px 40px 8px auto 1fr auto;
		align-items: center;
		cursor: pointer;
		min-height: 32px;
		user-select: none;
		transition: grid-template-columns 0.1s, margin-left 0.2s, opacity 0.2s;
	}

	li.option.selected {
		grid-template-columns: 28px 40px 8px auto 1fr auto;
		background-color: #eee;
	}

	li.option.not-selected {
		grid-template-columns: 5px 40px 8px auto 1fr auto;
		opacity: 0.6;
	}

	li.option:hover {
		background-color: #eee;
	}

	li.option:not(:last-child) {
		border-bottom: 2px solid #eee;
	}

	span.option-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		color: #1c2e58;
		color: #fff;
		background-color: #1c2e58;
		border-radius: 999em;
		height: 40px;
	}

	span.option-label {
		font-family: "JetBrains Mono", monospace;
		font-weight: bolder;
		position: relative;
	}


	span.ballot--valid {
		color: #00aa00;
	}
</style>