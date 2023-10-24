<script lang="ts">
	import { fetching } from "../store";
</script>

<div class={$fetching ? "visible" : ""}>
	<p>Working...</p>
</div>

<style>
	div {
		position: absolute;
		top: 0;
		left: 0;
		z-index: 99999;
		width: 100vw;
		max-width: 100%;
		height: 100%;
		pointer-events: none;
		display: flex;
		justify-content: center;
		background-color: rgba(0, 0, 0, var(--alpha, 0));
		backdrop-filter: blur(var(--blur, 0));
		opacity: var(--opacity, 0);
		--opacity: 1;
		transition: backdrop-filter 0.2s, background-color 0.2s;
	}

	div.visible {
		--blur: 5px;
		--opacity: 1;
		--alpha: 0.2;
		--top: 12px;
		pointer-events: all;
	}

	p {
		height: 36px;
		padding: 2px 24px;
		border-radius: 999em;
		background: #fff;
		position: relative;
		transform: translateY(var(--top, -40px));
		display: flex;
		align-items: center;
		font-family: "JetBrains Mono", monospace;
		transition: transform 0.2s;
	}
</style>
