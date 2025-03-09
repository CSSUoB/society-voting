<script lang="ts">
	import type { Optional } from "$lib/optional";

	export let label: Optional<string> = null;
	export let value = false;
</script>

<label>
	<input type="checkbox" bind:checked={value} />
	<span class="check material-symbols-rounded"></span>
	{label}
</label>

<style>
	label {
		display: grid;
		grid-template-columns: 0px 30px auto;
		align-items: center;
		cursor: pointer;
		user-select: none;
	}

	input {
		opacity: 0;
		margin: 0;
	}
	
	.check::after {
		opacity: 1;
		transition: 0.2s;
		color: #aaa;
		content: "check_box_outline_blank";
	}

	input:checked ~ .check::after {
		opacity: 1;
		color: #1c2e58;
		content: "check_box";
	}

</style>
