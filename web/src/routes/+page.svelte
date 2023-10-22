<script lang="ts">
	import Panel from "$lib/panel.svelte";
	const images = [
		"original",
		"christmas",
		"ball",
		"pride",
		"bbq",
		"graduation",
		"old-joe",
		"halloween",
	];
	const image = images[Math.floor(Math.random() * images.length)];
</script>

<Panel title="There are no upcoming elections">
	<p>Check this space later for updates. Here's a random TeX for now.</p>
	<img
		src={`https://cssuob.github.io/resources/dinosaur/tex_${image}.svg`}
		alt="A variant of TeX, the mascot of CSS"
		height="100px"
	/>
</Panel>

<style>
	img {
		margin-top: 16px;
	}
</style>
