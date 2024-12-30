<script lang="ts">
	import type { Optional } from "$lib/optional";

	export let name: Optional<string> = null;
	export let placeholder: Optional<string> = null;
	export let rows: Optional<number> = 5;
	export let value = "";
	let clazz: Optional<string> = null;
	export { clazz as class };
</script>

<textarea class={clazz} {name} {placeholder} {rows} bind:value={value} />

<style>
	textarea {
		border: 2px solid;
		border-radius: 4px;
		padding: 8px 12px;
		box-sizing: border-box;
		width: 100%;
		min-height: 5rem;
		resize: vertical;
		font-family: "Inter", sans-serif;
		font-size: medium;
		outline: 2px solid #1c2e5800;
		transition: 0.2s;
	}
	
	textarea:focus {
		border: 2px solid #1c2e58;
		outline: 1px solid #1c2e58;
		background-color: #eee;
	}
</style>
