<script lang="ts">
	import Button from "$lib/button.svelte";
	import Panel from "$lib/panel.svelte";

    const openInNewTab = () => {
        window.open("/presenter", "_blank").focus()
    }
</script>

<Panel title="Presenter Mode" headerIcon="open_in_new">
    <p>Display the vote code and a QR code to attendees</p>
    <br />
    <Button kind="emphasis" text="Open" on:click={openInNewTab} />
</Panel>
