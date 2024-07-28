<script lang="ts">
	import Panel from "$lib/panel.svelte";
	import { currentElection, elections } from "../store";
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

	$: upcomingElections = $elections?.filter((e) => !e.isActive && !e.isConcluded) ?? [];

	$: if ($currentElection && !$currentElection.hasVoted) {
		goto(`/vote`);
	} else if (upcomingElections.length > 0) {
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
