<script lang="ts">
	import Avatar from "$lib/avatar.svelte";
	import Banner from "$lib/banner.svelte";
	import Button from "$lib/button.svelte";
	import List from "$lib/list.svelte";
	import Panel from "$lib/panel.svelte";

	import run from "$lib/assets/run_for_election.svg";

	let candidates = [{ name: "" }];
</script>

<Panel title="Electing: President">
	<p>
		Lorem ipsum dolor sit amet consectetur adipisicing elit. Quas recusandae asperiores amet iure
		vitae magni quibusdam nulla minima, soluta ad pariatur autem dolor veritatis facere ullam culpa
		nostrum velit! Nesciunt.
	</p>
</Panel>

<Banner title="Interested in running?" kind="emphasis">
	<img slot="image" src={run} alt="" class="banner-image" />
	<svelte:fragment slot="content">
		<p>
			If you're thinking of running for election, go for it! Being on committee is super rewarding,
			super good fun, and a great way to bolster your CV and give back to the community.
		</p>
		<br />
		<Button text="Stand for election" kind="primary" />
	</svelte:fragment>
</Banner>

<Panel title="Candidates" headerIcon="groups">
	<List items={candidates} let:prop={candidate}>
		<li class="candidate">
			<Avatar name={candidate.name} />
			<p>{candidate.name}</p>
			<Button text="Withdraw" />
		</li>
	</List>
</Panel>

<Panel title="Admin stuff" headerIcon="admin_panel_settings">
	<li class="candidate">
		<Button text="Start election" />
	</li>
</Panel>

<style>
	img.banner-image {
		height: 120px;
		transform: translateY(16px);
	}

	li.candidate {
		padding: 8px 4px;
		display: grid;
		grid-template-columns: auto 1fr auto;
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
</style>
