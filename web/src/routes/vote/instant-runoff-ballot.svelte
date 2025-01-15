<script lang="ts">
	import BallotEntry from "./ballot-entry.svelte";
	import { type BallotEntry as BallotEntryT, user } from "../../store";
	import List from "../../lib/list.svelte";
	import Panel from "../../lib/panel.svelte";
	import { createEventDispatcher } from "svelte";

	export let candidates: Array<BallotEntryT>;

	let ballot: Array<BallotEntryT | undefined> = Array.from(Array(candidates.length));
	let errors = Array.from(Array(ballot.length));
	let validBallot: boolean;

	const dispatch = createEventDispatcher();

	const updateAndValidate = (index: number, changeEvent: Event) => {
		const id: number = parseInt((changeEvent.target as HTMLSelectElement).value);
		const candidate = candidates.find((c) => c.id === id);
		ballot = ballot.map((b, i) => (i === index ? candidate : b));

		errors = ballot.map((b, i) => {
			if (b === undefined && ballot.slice(i).filter((x) => x).length > 0)
				return "You cannot have gaps in your ranking";
			if (b === undefined) return undefined;
			if (ballot.filter((bb) => bb?.id === b.id).length !== 1)
				return `You cannot rank ${b.isRON ? b.name : b.name.split(" ")[0]} more than once`;
			return undefined;
		});
		
		validBallot = errors.filter((x) => x).length == 0;

		dispatch("update", {
			valid: validBallot,
			choices: validBallot ? ballot.filter((x) => x).map((b) => b?.id) : []
		});
	};
</script>

<Panel title="Your ballot">
	<p>There are {(ballot.length ?? 1) - 1} candidates on the ballot.</p>
	<p>Rank candidates in order of your choice. You do not need to rank all candidates.</p>
	<List items={ballot} let:prop={candidate}>
		<BallotEntry
			ballot={candidates ?? []}
			{candidate}
			error={errors[candidate.index]}
			on:change={(e) => updateAndValidate(candidate.index, e)}
		/>
	</List>
</Panel>
