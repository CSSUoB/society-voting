<script lang="ts">
	import Avatar from "$lib/avatar.svelte";
	import Banner from "$lib/banner.svelte";
	import Button from "$lib/button.svelte";
	import List from "$lib/list.svelte";
	import Panel from "$lib/panel.svelte";

	import run from "$lib/assets/run_for_election.svg";
	import { elections, user } from "../../../store";
	import { goto } from "$app/navigation";
	import { API } from "$lib/endpoints";
	import { _getElections } from "../../+layout";

	export let data: { id: number };
	$: election = $elections?.find((e) => e.id === data.id)!;
	$: if (!election) {
		goto("/", { replaceState: true });
	}

	const standOrWithdraw = async (id: number, stand: boolean) => {
		const response = await fetch(API.ELECTION_STAND, {
			method: stand ? "POST" : "DELETE",
			body: JSON.stringify({ id }),
		});

		if (response.ok) {
			elections.set(await _getElections());
		}
	};
</script>

<svelte:head>
	<title>Electing: {election.roleName}</title>
</svelte:head>
<Panel title={`Electing: ${election.roleName}`}>
	<p>{election.description}</p>
</Panel>

{#if !$user.admin && !election.candidates?.some((c) => c.isMe)}
	<Banner title="Interested in running?" kind="emphasis">
		<img slot="image" src={run} alt="" class="banner-image" />
		<svelte:fragment slot="content">
			<p>
				If you're thinking of running for election, go for it! Being on committee is super
				rewarding, super good fun, and a great way to bolster your CV and give back to the
				community.
			</p>
			<br />
			<Button
				text="Stand for election"
				kind="primary"
				on:click={() => standOrWithdraw(election.id, true)}
			/>
		</svelte:fragment>
	</Banner>
{/if}

<Panel title="Candidates" headerIcon="groups">
	<List items={election.candidates ?? []} let:prop={candidate}>
		<li class="candidate">
			<Avatar name={candidate.name} />
			<p>
				{candidate.name}{#if candidate.isMe}
					<span><small>You</small></span>
				{/if}
			</p>
			{#if candidate.isMe}
				<Button text="Withdraw" on:click={() => standOrWithdraw(election.id, false)} />
			{/if}
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

	li.candidate p > span {
		margin-left: 8px;
		background: #000;
		color: #fff;
		padding: 0 8px;
		border-radius: 4px;
		text-transform: uppercase;
		font-family: "JetBrains Mono";
		font-weight: bold;
	}
</style>
