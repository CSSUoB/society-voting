<script lang="ts">
	import { AvatarBeam } from "svelte-boring-avatars";

	export let name: string = "";
	export let size: number = 48;

	const colours = [
		["#ff9c00", "#9d68dc"], // orange
		["#58cad6", "#4a9efe"], // blue
		["#03e421", "#3b37b6"], // green
		["#fdff8a", "#3b8b94"], // yellow
		["#1234ff", "#cc22ff"], // magenta
		["#cc5500", "#112233"], // brown
	];

	$: colour = colours[name.charCodeAt(0) % colours.length || 0];
</script>

{#key name}
	<div class="avatar" style="--size: {size}px">
		<AvatarBeam {size} {name} colors={colour} />
	</div>
{/key}

<style>
	.avatar {
		height: var(--size);
		width: var(--size);
		border: 2px solid #1c2e58;
		border-radius: 999em;
	}
</style>
