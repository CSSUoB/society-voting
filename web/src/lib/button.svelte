<script lang="ts">
	import type { Optional } from "$lib/optional";

	export let text: Optional<string> = null;
	export let icon: Optional<string> = null;
	export let kind: "flat" | "inline" | "emphasis" | "primary" = "flat";
</script>

<button class={`${kind} ${icon && !text ? "icon-only" : ""}`}>
	{#if text}
		{text}
	{/if}
	{#if icon}
		<span class="material-symbols-rounded">{icon}</span>
	{/if}
</button>

<style>
	button {
		height: 36px;
		padding: 0 16px;
		font-family: "JetBrains Mono", monospace;
		font-weight: bold;
		border-radius: 999em;
		border: 2px solid;
		display: flex;
		align-items: center;
		gap: 4px;
	}

	button.icon-only {
		width: 36px;
		padding: 0;
		display: grid;
		align-items: center;
	}

	button.flat {
		border: 1px solid #ddd;
	}

	button.primary {
		background-color: #03e421;
	}
</style>
