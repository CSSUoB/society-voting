<script lang="ts">
	import InputGroup from "$lib/input-group.svelte";
	import Input from "$lib/input.svelte";
	import Panel from "$lib/panel.svelte";
	import Textarea from "$lib/textarea.svelte";
</script>

<Panel title="Election" headerIcon="people">
	<div class="create-election">
		<InputGroup label="Role name" helpText="Maximum 64 characters.">
			<Input placeholder="Publicity Officer" name="roleName" required={true} />
		</InputGroup>
		<InputGroup label="Description" helpText="Maximum 1024 characters.">
			<Textarea
				placeholder="The Publicity Officer is responsible for..."
				name="description"
			/>
		</InputGroup>
	</div>
</Panel>

<style>
	div.create-election {
		display: flex;
		gap: 0.5rem;
		flex-direction: column;
	}	
</style>
