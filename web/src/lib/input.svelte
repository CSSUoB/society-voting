<script lang="ts">
	import type { Optional } from "$lib/optional";

	export let name: Optional<string> = null;
	export let type: Optional<string> = null;
	export let value = "";
</script>

<input {name} {type} {value} />

<style>
	input {
		height: 30px;
		padding: 2px 12px;
		border: 2px solid;
		border-radius: 4px;
		margin: 0 0 12px 0;
	}
</style>
