<script lang="ts">
	import type { Optional } from "$lib/optional";

	export let name: Optional<string> = null;
	export let placeholder: Optional<string> = null;
	export let value = "";
	export let required = false
	let clazz: Optional<string> = null;
	export { clazz as class };
</script>

<input class={clazz} {name} {placeholder} {required} bind:value={value} />

<style>
	input {
		height: 40px;
		padding: 2px 12px;
		border: 2px solid;
		border-radius: 4px;
		width: 100%;
		box-sizing: border-box;
		font-family: "Inter", sans-serif;
		font-size: medium;
		outline: 2px solid #1c2e5800;
		transition: 0.2s;
	}
	
	input:focus {
		border: 2px solid #1c2e58;
		outline: 1px solid #1c2e58;
		background-color: #eee;
	}
</style>
