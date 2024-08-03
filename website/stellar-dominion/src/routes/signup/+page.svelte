<script lang="ts">
	import '$styles/app.css';
	import {
		CenteredWrapper,
		FormField,
		StyledButton,
		StyledError,
		StyledLink,
		StyledText,
		StyledTitle
	} from '$lib/components';

	import heroImage, { HOMEPAGE_HERO_IMAGE } from '$lib/stores/ui/heroImage';
	import heroContainer, { HOMEPAGE_HERO_CONTAINER_PROPS } from '$lib/stores/ui/heroContainer';

	/** @type {import('./$types').ActionData} */
	export let form: HTMLFormElement;

	/** @type {import('./$types').PageData} */
	export let data;

	function resetFormError() {
		if (!form) {
			return;
		}
		form.message = '';
	}

	heroImage.set(HOMEPAGE_HERO_IMAGE);
	heroContainer.set(HOMEPAGE_HERO_CONTAINER_PROPS);
</script>

<CenteredWrapper>
	<div class="fixed left-4 top-4">
		<p class="text-secondary">
			Already have an account yet? Click <StyledLink text="here" link="/login" /> to login!"
		</p>
	</div>

	<CenteredWrapper width="w-full" height="h-1/5">
		<StyledTitle text="Stellar Dominion" />
		<StyledText text="Sign up" />
	</CenteredWrapper>

	<CenteredWrapper width="w-full" height="h-3/5">
		<form method="POST" action="?/signup" class="flex flex-col grow justify-evenly">
			<FormField label="universe:" labelId="universe">
				<select id="universe" name="universe">
					{#each data.universes as universe}
						<option value={universe.id}>{universe.name}</option>
					{/each}
				</select>
			</FormField>
			<FormField label="email:" labelId="email">
				<input
					id="email"
					type="text"
					name="email"
					placeholder="Enter your email address"
					required
					value={form?.email ?? ''}
					on:input={resetFormError}
				/>
			</FormField>
			<FormField label="password:" labelId="password">
				<input
					id="password"
					type="text"
					name="password"
					placeholder="Enter your password"
					required
					on:input={resetFormError}
				/></FormField
			>
			<FormField label="player:" labelId="player">
				<input
					id="player"
					type="text"
					name="player"
					placeholder="Choose a name"
					required
					on:input={resetFormError}
				/></FormField
			>
			<StyledButton text="Sign up" />
		</form>

		{#if form?.message}
			<div class="fixed bottom-4">
				<StyledError text="Failed to sign up: {form.message}" />
			</div>
		{/if}
	</CenteredWrapper>
</CenteredWrapper>
