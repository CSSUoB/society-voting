<script lang="ts">
	import Panel from "$lib/panel.svelte";
	import { elections } from "../store";
	import { goto } from "$app/navigation";
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

	$: currentElections = $elections?.filter((e) => e.isActive) ?? [];
	$: upcomingElections = $elections?.filter((e) => !e.isActive) ?? [];

	$: if (currentElections.length > 0) {
		goto(`/vote/${currentElections[0].id}`);
	}

	$: if (upcomingElections.length > 0) {
		goto(`/election/${upcomingElections[0].id}`);
	}
</script>

<svelte:head>
	<title>CSS Elects</title>
</svelte:head>

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
