<script lang="ts">
	import type { Optional } from "$lib/optional";

	export let name: Optional<string> = null;
	export let type: Optional<string> = null;
	export let value = "";
</script>

<input {name} {type} {value} />
