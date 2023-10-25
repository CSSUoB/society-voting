<script lang="ts">
	import logo from "$lib/assets/logo.svg";
	import { createEventDispatcher } from "svelte";
	import Button from "./button.svelte";

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
</nav>

<style>
	nav {
		display: grid;
		grid-template-columns: auto 1fr;
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

	@media only screen and (max-width: 850px) {
		nav {
			grid-template-columns: auto auto 1fr;
		}

		nav > :global(button.menu) {
			display: inherit;
		}
	}
</style>
