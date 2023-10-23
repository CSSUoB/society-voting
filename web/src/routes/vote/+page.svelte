<script lang="ts">
	import Avatar from "$lib/avatar.svelte";
	import ButtonGroup from "$lib/button-group.svelte";
	import Button from "$lib/button.svelte";
	import List from "$lib/list.svelte";
	import Panel from "$lib/panel.svelte";

	let candidates = [{ name: "" }];
</script>

<svelte:head>
	<title>Vote for: X</title>
</svelte:head>

<Panel title="Electing: President">
	<p>
		Lorem ipsum dolor sit amet consectetur adipisicing elit. Quas recusandae asperiores amet iure
		vitae magni quibusdam nulla minima, soluta ad pariatur autem dolor veritatis facere ullam culpa
		nostrum velit! Nesciunt.
	</p>
</Panel>

<Panel title="Your ballot">
	<p>Rank candidates in order of your choice. You do not need to rank all candidates.</p>
	<List items={candidates} let:prop={candidate}>
		<li class="candidate">
			<p class="candidate-position">{candidate.index + 1}</p>
			<Avatar name={candidate.name} />
			<p>{candidate.name}</p>
			<ButtonGroup>
				<Button icon="arrow_upward" />
				<Button icon="arrow_downward" />
			</ButtonGroup>
		</li>
	</List>
</Panel>

<Panel title="Submit" kind="emphasis">
	<p />
</Panel>

<style>
	li.candidate {
		padding: 8px 4px;
		display: grid;
		grid-template-columns: 48px auto 1fr auto;
		align-items: center;
		gap: 8px;
	}

	li.candidate:not(:last-child) {
		border-bottom: 2px solid #eee;
	}

	li.candidate p {
		text-overflow: ellipsis;
		overflow: hidden;
	}

	p.candidate-position {
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.5rem;
		font-family: "JetBrains Mono", monospace;
		font-weight: bolder;
		background-color: #1c2e58;
		color: #fff;
		height: 48px;
		border-radius: 999em;
		position: relative;
	}
</style>
