<script lang="ts">
	import Button from "$lib/button.svelte";
	import Panel from "$lib/panel.svelte";
</script>

<Panel title="Hi Vishwas!" headerIcon="waving_hand">
	<Button slot="header-action" text="Use a different name" kind="emphasis" />
	<!-- TODO: Move to next line and add blurb about voting system -->
	<p>Welcome to CSS' voting system!</p>
	<span
		><p>Not you?</p>
		<Button text="Log out" kind="emphasis" />
	</span>
</Panel>
