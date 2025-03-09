<script lang="ts">
	import type { Optional } from "$lib/optional";

	export let text: Optional<string> = null;
	export let icon: Optional<string> = null;
	export let kind: "flat" | "inline" | "emphasis" | "danger" | "primary" = "flat";
	export let type: "button" | "submit" | "reset" | null | undefined = undefined;
	export let name: Optional<string> = null;
	export let disabled = false;
	let clazz: Optional<string> = null;
	export { clazz as class };
</script>

<button
	on:click
	class={`${kind} ${icon && !text ? "icon-only" : ""} ${clazz}`}
	{type}
	{name}
	{disabled}
>
	{#if text}
		{text}
	{/if}
	{#if icon}
		<span class="material-symbols-rounded">{icon}</span>
	{/if}
</button>

<style>
	button {
		--border-colour: #1c2e58;
		--background-colour: #eee;
		--background-colour--shadow: #ccc;
		--background-colour--focus: #ddd;
		--background-colour--shadow--focus: #bbb;
		--border-width: 2px;
		--elevation: 5px;
		height: 36px;
		padding: 0 16px calc(var(--elevation) + 2px) 16px;
		font-family: "JetBrains Mono", monospace;
		font-weight: bold;
		border-radius: 999em;
		border: var(--border-width) solid var(--border-colour);
		background-color: var(--background-colour--shadow);
		display: flex;
		align-items: center;
		gap: 4px;
		position: relative;
		z-index: 0;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	button::after {
		content: "";
		position: absolute;
		width: 100%;
		height: 100%;
		top: calc(-1 * var(--elevation));
		left: -1px;
		border-radius: 999em;
		background-color: var(--background-colour);
		z-index: -1;
		border: 1px dashed var(--border-colour);
		transition: background-color 0.2s;
	}

	button:disabled {
		opacity: 0.8;
		cursor: not-allowed;
	}

	button:not(:disabled):hover,
	button:not(:disabled):focus {
		background-color: var(--background-colour--shadow--focus);
	}

	button:not(:disabled):hover::after,
	button:not(:disabled):focus::after {
		background-color: var(--background-colour--focus);
	}

	button:not(:disabled):active {
		--elevation: 0px;
	}

	button.icon-only {
		width: 36px;
		padding: 0;
		display: grid;
		align-items: center;
	}

	button.flat {
		--background-colour--shadow: #eee;
		--background-colour--shadow--focus: #ddd;
		--border-width: 1px;
		--elevation: -2px;
	}

	button.primary {
		--border-colour: #000;
		--background-colour: #03e421;
		--background-colour--shadow: #03c41d;
		--background-colour--focus: #03c91d;
		--background-colour--shadow--focus: #04a71a;
	}

	button.danger {
		--border-colour: #000;
		--background-colour: #fc5d55;
		--background-colour--shadow: #da372e;
		--background-colour--focus: #e83a31;
		--background-colour--shadow--focus: #cf342c;
	}

	button.inline {
		border: 1px dashed var(--border-colour);
		--background-colour--shadow: transparent;
		--background-colour--shadow--focus: rgba(0, 0, 0, 0.2);
		text-decoration: underline;
		height: 24px;
		--elevation: -2px;
	}

	button.inline::after,
	button.flat::after {
		display: none;
	}
</style>
