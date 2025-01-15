<script lang="ts">
	import type { Optional } from "$lib/optional";

	export let label: Optional<string> = null;
	export let helpText: Optional<string> = null;
	let clazz: Optional<string> = null;
	export { clazz as class };
</script>

<div class={`input-group ${clazz}`}>
	{#if label}
		<p class="label"><b>{label}</b></p>
	{/if}
	<slot />
	{#if helpText}
		<p class="help-text">{helpText}</p>
	{/if}
</div>

<style>
	p.label {
		font-family: "JetBrains Mono", monospace;
	}

	p.help-text {
		font-size: smaller;
		color: #555;
	}
</style>
