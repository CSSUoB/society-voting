<script lang="ts">
	import { createEventDispatcher, onDestroy, onMount } from "svelte";

	export let dialog: HTMLDialogElement;
	export let title = "Title";
	let ref: HTMLDivElement;
	let portal: HTMLDivElement;

	onMount(() => {
		portal = document.createElement("div");
		portal.className = "portal";
		document.body.appendChild(portal);
		portal.appendChild(ref);
	});

	onDestroy(() => {
		document.body.removeChild(portal);
	});

	const dispatch = createEventDispatcher();
	const submit = (e: SubmitEvent) => {
		if ((e.submitter as HTMLButtonElement).name !== "submit") return;
		e.preventDefault();
		const formData = new FormData(e.target as HTMLFormElement);

		dispatch("submit", Object.fromEntries(formData));

		dialog.close();
	};
</script>

<div bind:this={ref}>
	<dialog bind:this={dialog} on:close>
		<h2>{title}</h2>
		<form method="dialog" on:submit={submit}>
			<slot />
			<div class="actions">
				<slot name="actions" />
			</div>
		</form>
	</dialog>
</div>

<style>
	dialog {
		border-radius: 20px;
		border: 3px solid #1c2e58;
		padding: 16px;
		min-width: min(calc(100vw - 32px - 32px), 400px);
	}

	dialog::backdrop {
		background-image: radial-gradient(rgba(0, 0, 0, 0.1) 20%, transparent 20%),
			radial-gradient(rgba(0, 0, 0, 0.1) 20%, transparent 20%);
		background-position: 0 0, 100px 100px;
		background-size: 40px 40px;
		backdrop-filter: blur(5px);
	}

	div.actions {
		display: flex;
		margin-top: 12px;
		gap: 8px;
		justify-content: flex-end;
	}
</style>
