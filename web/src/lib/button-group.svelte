<div>
	<slot />
</div>

<style>
	div {
		display: flex;
		gap: 4px;
	}

	div > :global(button) {
		border-radius: 8px;
	}

	div > :global(button:last-child) {
		border-radius: 16px 36px 36px 16px;
		padding-right: 4px;
	}

	div > :global(button:first-child) {
		border-radius: 36px 16px 16px 36px;
		padding-left: 4px;
	}
</style>
