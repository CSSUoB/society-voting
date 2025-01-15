<script lang="ts">
	import InputGroup from "$lib/input-group.svelte";
	import Input from "$lib/input.svelte";
	import Panel from "$lib/panel.svelte";
	import Textarea from "$lib/textarea.svelte";
</script>

<Panel title="Referendum" headerIcon="rule">
	<div class="create-referendum">
		<InputGroup label="Title" helpText="Maximum 64 characters.">
			<Input placeholder="TeX Dictatorship" name="title" required={true} />
		</InputGroup>
		<InputGroup label="Question" helpText="Maximum 256 characters.">
			<Input placeholder="Should the Committee be abolished and replaced with a dictatorship with TeX as the sole leader?" name="question" required={true} />
		</InputGroup>
		<InputGroup label="Description" helpText="Maximum 1024 characters.">
			<Textarea
				placeholder="This motion if passed will..."
				name="description"
			/>
		</InputGroup>
	</div>
</Panel>

<style>
	.create-referendum {
		display: flex;
		gap: 0.5rem;
		flex-direction: column;
	}	
</style>
