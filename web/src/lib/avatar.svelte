<script lang="ts">
	export let name: string = "";
	export let size: number = 48;

	const colours = [
		"ff9c00,9d68dc", // orange
		"58cad6,4a9efe", // blue
		"03e421,3b37b6", // green
		"fdff8a,3b8b94", // yellow
		"1234ff,cc22ff", // magenta
		"cc5500,112233", // brown
	];

	$: colour = colours[name.charCodeAt(0) % colours.length ?? 0];

	const hash = async (data: string) => {
		const encoder = new TextEncoder();
		const encoded = encoder.encode(data);
		const hashBuffer = await crypto.subtle.digest("SHA-256", encoded);
		const hashArray = Array.from(new Uint8Array(hashBuffer)); // convert buffer to byte array
		const hashHex = hashArray.map((b) => b.toString(16).padStart(2, "0")).join(""); // convert bytes to hex string
		return hashHex;
	};
	$: hashedName = hash(name);
</script>

{#await hashedName}
	<img src="" alt={`Minimalistic avatar representing ${name}`} height={size} width={size} />
{:then n}
	<img
		src={`https://source.boringavatars.com/beam/${size}/${n}?colors=${colour}`}
		alt={`Minimalistic avatar representing ${name}`}
		height={size}
		width={size}
	/>
{/await}

<style>
	img {
		border: 2px solid #1c2e58;
		border-radius: 999em;
	}
</style>
