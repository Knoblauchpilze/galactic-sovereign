import { writable } from 'svelte/store';

export type HeroContainerProps = {
	width: string;
	height: string;
	color: string;
};

const HOMEPAGE_HERO_CONTAINER_PROPS: HeroContainerProps = {
	width: 'w-3/5',
	height: 'h-4/5',
	color: 'bg-overlay'
};

const GAME_HERO_CONTAINER_PROPS: HeroContainerProps = {
	width: 'w-full',
	height: 'h-full',
	color: 'bg-transparent'
};

export { HOMEPAGE_HERO_CONTAINER_PROPS, GAME_HERO_CONTAINER_PROPS };

export default writable(HOMEPAGE_HERO_CONTAINER_PROPS);
