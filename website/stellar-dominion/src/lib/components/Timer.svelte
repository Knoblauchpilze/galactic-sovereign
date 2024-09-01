<script lang="ts">
	// https://svelte.dev/repl/aa6c438ff6c04851b63328e3eb54466c?version=3.32.3
	import { onMount, onDestroy } from 'svelte';
	import { StyledText } from '$lib/components';
	import { msToTimeStringOrFinished } from '$lib/time';

	export let durationMs;
	// https://svelte.dev/repl/5f4a327999cd49e5a79e91f6fbe994c8?version=3.59.2
	// https://stackoverflow.com/questions/29689966/how-to-define-type-for-a-function-callback-as-any-function-type-not-universal
	export let onFinished: () => void;

	let msElapsed = 0;
	let interval: number;
	// https://geoffrich.net/posts/svelte-$-meanings/
	$: remainingMs = durationMs - msElapsed;
	$: remaining = msToTimeStringOrFinished(remainingMs);

	$: textColor = remainingMs <= 0 ? 'text-enabled' : 'text-white';

	const UPDATE_INTERVAL_MS = 1000;

	onMount(() => {
		interval = setInterval(() => {
			msElapsed += UPDATE_INTERVAL_MS;

			if (remainingMs <= 0) {
				clearInterval(interval);
				onFinished();
			}
		}, UPDATE_INTERVAL_MS);
	});

	onDestroy(() => {
		clearInterval(interval);
	});
</script>

<StyledText text={remaining} {textColor} />
