<script lang="ts">
	type T = $$Generic;
	export let items: Array<T> = [];
</script>

<ul>
	{#each items as item, i}
		<slot prop={{ ...item, index: i }} />
	{/each}
</ul>

<style>
	ul {
		list-style: none;
		padding: 0;
		margin: 0;
	}
</style>
