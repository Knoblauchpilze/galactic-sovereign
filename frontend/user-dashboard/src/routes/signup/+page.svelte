<script lang="ts">
	import { FlexContainer, StyledLink, StyledTitle } from '@totocorpsoftwareinc/frontend-toolkit';
	import { FormField, StyledButton, StyledError, StyledText } from '$lib/components';

	interface Props {
		form: HTMLFormElement;
	}

	let { form = $bindable() }: Props = $props();

	function resetFormError() {
		if (!form) {
			return;
		}
		form.message = '';
	}
</script>

<FlexContainer>
	<div class="fixed right-4 top-4">
		<StyledLink text="Login" link="/login" />
	</div>

	<FlexContainer extensible={false} styling="h-1/5">
		<StyledTitle text="Admin dashboard" />
		<StyledText text="Sign up" />
	</FlexContainer>

	<FlexContainer extensible={false} styling="h-3/5">
		<form method="POST" action="?/signup" class="flex flex-col flex-1 justify-evenly">
			<FormField label="email:" labelId="email">
				<input
					id="email"
					type="text"
					name="email"
					placeholder="Enter your email address"
					required
					value={form?.email ?? ''}
					oninput={resetFormError}
				/>
			</FormField>
			<FormField label="password:" labelId="password">
				<input
					id="password"
					type="text"
					name="password"
					placeholder="Enter your password"
					required
					oninput={resetFormError}
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
