<script lang="ts">
	import Avatar from "$lib/avatar.svelte";
	import type { Indexed, Optional } from "$lib/optional";
	import type { BallotEntry } from "../../store";
	export let ballot: Array<BallotEntry>;
	export let candidate: Indexed<Optional<BallotEntry>>;
	export let error: Optional<string> = null;
</script>

<li class="candidate">
	<p class="candidate-position">{candidate.index + 1}</p>
	{#if candidate.isRON}
		<span class="candidate-avatar-placeholder--ron">
			<span class="material-symbols-rounded">refresh</span>
		</span>
	{:else if candidate.id}
		<Avatar name={candidate.name} />
	{:else}
		<span class="candidate-avatar-placeholder" />
	{/if}
	<select class="candidate-select" value={candidate.id ?? "-1"} on:change>
		<option value="-1">No-one</option>
		{#each ballot as b}
			<option value={b.id}>{b.name}</option>
		{/each}
	</select>
	<span />
	{#if error}
		<p class="ballot-error">
			<small>{error}</small>
			<span class="ballot--invalid material-symbols-rounded">cancel</span>
		</p>
	{:else}
		<span class="ballot--valid material-symbols-rounded">check_circle</span>
	{/if}
</li>

<style>
	li.candidate {
		padding: 8px 4px;
		display: grid;
		grid-template-columns: 52px auto auto 1fr auto;
		align-items: center;
		gap: 8px;
	}

	li.candidate:not(:last-child) {
		border-bottom: 2px solid #eee;
	}

	p.candidate-position {
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.5rem;
		font-family: "JetBrains Mono", monospace;
		font-weight: bolder;
		background-color: #1c2e58;
		color: #fff;
		height: 48px;
		border-radius: 999em;
		position: relative;
		border: 2px solid #1c2e58;
	}

	span.candidate-avatar-placeholder,
	span.candidate-avatar-placeholder--ron {
		display: block;
		height: 48px;
		border: 2px solid #1c2e58;
		width: 48px;
		border-radius: 100%;
		position: relative;
		background: linear-gradient(
			45deg,
			transparent calc(50% - 1px),
			#1c2e58 calc(50% - 1px),
			#1c2e58 calc(50% + 1px),
			transparent calc(50% + 1px)
		);
	}

	span.candidate-avatar-placeholder--ron {
		background: transparent;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	select.candidate-select {
		height: 48px;
		border: 2px solid #111;
		border-radius: 999em;
		padding: 0 24px;
		max-width: calc(100vw - 24px - 52px - 52px - 16px - 24px - 16px);
	}

	p.ballot-error {
		display: flex;
		gap: 4px;
		align-items: center;
	}

	span.ballot--invalid,
	span.ballot--valid {
		font-size: 32px;
	}

	span.ballot--invalid {
		color: #cc0000;
	}

	span.ballot--valid {
		color: #00aa00;
	}

	@media only screen and (max-width: 850px) {
		li.candidate {
			grid-template-columns: 52px auto auto 1fr;
		}

		p.ballot-error {
			grid-area: 2/1/3/4;
		}
	}
</style>
