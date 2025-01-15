<script lang="ts">
	import type { Optional } from "$lib/optional";

	export let kind: "regular" | "emphasis" | "primary" = "regular";
	export let title: string | undefined = undefined;
	export let headerIcon: Optional<string> = null;
</script>

<section class={kind}>
	{#if title}
		<span class="header">
			{#if headerIcon}<span class="material-symbols-rounded">{headerIcon}</span>{/if}
			<h1>{title}</h1>
			<slot name="header-action" />
		</span>
	{/if}
	<slot />
</section>

<style>
	section {
		background-color: #fff;
		padding: 16px;
		border-radius: 12px;
		border: 3px solid #1c2e58;
	}

	section:first-child {
		border-top-right-radius: 20px;
		border-top-left-radius: 20px;
	}

	section:last-child {
		border-bottom-right-radius: 20px;
		border-bottom-left-radius: 20px;
	}

	section.emphasis {
		background-color: #c3abdf;
	}

	section.primary {
		background-color: #93f193;
	}

	span.header {
		display: flex;
		gap: 8px;
		align-items: flex-start;
		margin-bottom: 12px;
		flex-wrap: wrap;
	}

	h1 {
		margin-bottom: -3px;
		text-overflow: ellipsis;
		overflow: hidden;
		flex-grow: 1;
	}

	span.material-symbols-rounded {
		font-size: 30px;
	}
</style>
