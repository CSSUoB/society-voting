<script lang="ts">
	import logo from "$lib/assets/logo.svg";
	import { createEventDispatcher } from "svelte";
	import Button from "./button.svelte";
	import { navigating } from "../store";

	let menuOpen = false;

	const dispatch = createEventDispatcher();
	const toggleMenu = (open: boolean) => {
		menuOpen = !open;
		dispatch("menuToggle", menuOpen);
	};
</script>

<nav>
	<Button class="menu" icon={menuOpen ? "close" : "menu"} on:click={() => toggleMenu(menuOpen)} />
	<img src={logo} alt="Logo" />
	<!-- TODO: Replace society name with value from config.yml -->
	<span>CSS Elects</span>
	{#if $navigating}
		<div class="spinner"></div>
	{/if}
</nav>

<style>
	nav {
		display: grid;
		grid-template-columns: auto 1fr auto;
		gap: 20px;
		align-items: center;
		padding: 12px;
		border-bottom: 3px solid #1c2e58;
		background-color: #000;
		grid-area: var(--location);
		position: sticky;
		top: 0;
		z-index: 1000;
	}

	nav > :global(button.menu) {
		display: none;
	}

	img {
		height: 40px;
		width: 40px;
	}

	span {
		color: #fff;
		font-family: "JetBrains Mono", monospace;
		font-weight: bold;
		font-size: 1.2rem;
	}
	
	div.spinner {
	    width: 25px;
		height: 25px;
		border: 5px solid #FFF;
		border-bottom-color: transparent;
		border-radius: 50%;
		display: inline-block;
		box-sizing: border-box;
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		from {
			transform: rotate(0deg);
		}

		to {
			transform: rotate(360deg);
		}
	}

	@media only screen and (max-width: 850px) {
		nav {
			grid-template-columns: auto auto 1fr auto;
		}

		nav > :global(button.menu) {
			display: inherit;
		}
	}
</style>
