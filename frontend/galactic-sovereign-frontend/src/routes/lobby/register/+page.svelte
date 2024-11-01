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


	interface Props {
		form: HTMLFormElement;
		data: any;
	}

	let { form = $bindable(), data }: Props = $props();

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
	<div class="fixed right-4 top-4">
		<p class="text-secondary">
			Back to the <StyledLink text="lobby" link="/lobby" />
		</p>
	</div>

	<FlexContainer extensible={false} styling="h-1/5">
		<StyledTitle text="Galactic Sovereign" />
		<StyledText text="Join new universe" />
	</FlexContainer>

	<FlexContainer extensible={false} styling="h-3/5">
		{#if data.universes.length > 0}
			<form method="POST" action="?/register" class="flex flex-col flex-1 justify-evenly">
				<FormField label="universe:" labelId="universe">
					<select id="universe" name="universe">
						{#each data.universes as universe}
							<option value={universe.id}>{universe.name}</option>
						{/each}
					</select>
				</FormField>
				<FormField label="player:" labelId="player">
					<input
						id="player"
						type="text"
						name="player"
						placeholder="Choose a name"
						required
						oninput={resetFormError}
					/></FormField
				>
				<StyledButton text="Start" />
			</form>
		{:else}
			<StyledError text="You already have an account in all universes, please login!" />
			<StyledLink text="Login" link="/lobby/login" showAsButton={true} />
		{/if}

		{#if form?.message}
			<div class="fixed bottom-4">
				<StyledError text="Failed to register in universe: {form.message}" />
			</div>
		{/if}
	</FlexContainer>
</FlexContainer>
