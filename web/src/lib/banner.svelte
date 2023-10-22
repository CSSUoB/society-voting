<script lang="ts">
	export let title: string = "Title";
	export let kind: "emphasis" | "success" | "caution" | "info" = "info";
</script>

<section class={kind}>
	<slot name="image" />
	<div>
		<h2>{title}</h2>
		<slot name="content" />
	</div>
</section>

<style>
	section {
		padding: 16px;
		border-radius: 12px;
		border: 3px solid;
		display: grid;
		grid-template-columns: auto 1fr;
		gap: 16px;
		align-items: flex-end;
	}

	section.emphasis {
		background-color: #111;
		border-color: #1c2e58;
		color: #eee;
	}

	section.caution {
		background-color: yellow;
		border-color: #1c2e58;
		color: #000;
	}

	h2 {
		margin-bottom: 12px;
	}
</style>
