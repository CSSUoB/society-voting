<script>
	import logo from "$lib/assets/logo.svg";
</script>

<nav>
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
</style>
