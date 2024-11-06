import { writable } from 'svelte/store';

export type HeroContainerProps = {
	width: string;
	height: string;
	color: string;
};

const HOMEPAGE_HERO_CONTAINER_PROPS: HeroContainerProps = {
	width: 'w-3/12',
	height: 'h-3/6',
	color: 'bg-overlay'
};

export { HOMEPAGE_HERO_CONTAINER_PROPS };

export default writable(HOMEPAGE_HERO_CONTAINER_PROPS);
