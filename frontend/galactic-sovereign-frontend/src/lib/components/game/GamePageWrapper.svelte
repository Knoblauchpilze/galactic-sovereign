<script lang="ts">
	import { type Snippet } from 'svelte';

	import { FlexContainer, Header, StyledText } from '@totocorpsoftwareinc/frontend-toolkit';

	import { floorToInteger, toFlooredShortString } from '$lib/displayUtils';
	import type { PlanetResourceUiDto } from '$lib/communication/ui/planetResourceUiDto';

	interface Props {
		universeName: string;
		planetName: string;
		playerName: string;
		resources: PlanetResourceUiDto[];
		children?: Snippet;
	}

	let { universeName, planetName, playerName, resources, children }: Props = $props();

	function resourceTextColor(resource: PlanetResourceUiDto): string {
		if (resource.amount < resource.storage) {
			return 'text-white';
		}
		return 'text-disabled';
	}
	function productionTextColor(resource: PlanetResourceUiDto): string {
		if (resource.amount < resource.storage) {
			return 'text-enabled';
		}
		return 'text-disabled';
	}
</script>

<FlexContainer>
	<!-- https://stackoverflow.com/questions/67852559/pass-svelte-component-as-props -->
	<Header>
		<StyledText text={universeName} textColor="text-white" />
		<StyledText text={playerName} textColor="text-white" />
		<StyledText text={planetName} textColor="text-white" />
		<form method="POST" action="?/backToLobby">
			<button class="hover:underline">Back to the lobby</button>
		</form>
		<form method="POST" action="?/logout">
			<button class="hover:underline">Logout</button>
		</form>
	</Header>

	<FlexContainer align={'stretch'} styling={'w-full'}>
		<FlexContainer vertical={false} extensible={false} bgColor={'bg-black'} styling={'w-full'}>
			{#each resources as resource}
				<FlexContainer vertical={false} extensible={false} styling={'space-between'}>
					<StyledText text="{resource.name}:" textColor="text-white" />
					<StyledText
						text={floorToInteger(resource.amount).toString()}
						textColor={resourceTextColor(resource)}
						styling="px-1"
					/>
					<StyledText
						text="(+{floorToInteger(resource.production)}/h)"
						textColor={productionTextColor(resource)}
						styling="pr-1"
					/>
					<StyledText
						text="(storage: {toFlooredShortString(resource.storage)})"
						textColor="text-white"
					/>
				</FlexContainer>
			{/each}
		</FlexContainer>

		{@render children?.()}
	</FlexContainer>
</FlexContainer>
