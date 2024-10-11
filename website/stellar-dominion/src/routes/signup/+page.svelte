<script lang="ts">
	import {
		FlexContainer,
		FormField,
		StyledButton,
		StyledError,
		StyledLink,
		StyledText,
		StyledTitle
	} from '$lib/components';

	import heroImage, { HOMEPAGE_HERO_IMAGE } from '$lib/stores/ui/heroImage';
	import heroContainer, { HOMEPAGE_HERO_CONTAINER_PROPS } from '$lib/stores/ui/heroContainer';

	export let form: HTMLFormElement;

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

<FlexContainer>
	<div class="fixed left-4 top-4">
		<p class="text-secondary">
			Already have an account yet? Click <StyledLink text="here" link="/login" /> to login!"
		</p>
	</div>

	<FlexContainer extensible={false} styling="h-1/5">
		<StyledTitle text="Stellar Dominion" />
		<StyledText text="Sign up" />
	</FlexContainer>

	<FlexContainer extensible={false} styling="h-3/5">
		<form method="POST" action="?/signup" class="flex flex-col flex-1 justify-evenly">
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
	</FlexContainer>
</FlexContainer>
