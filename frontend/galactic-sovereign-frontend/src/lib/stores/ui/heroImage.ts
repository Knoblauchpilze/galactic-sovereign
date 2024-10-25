import { writable } from 'svelte/store';

const HOMEPAGE_HERO_IMAGE: string = 'bg-homepage';
const GAME_HERO_IMAGE: string = 'bg-overview';

export { HOMEPAGE_HERO_IMAGE, GAME_HERO_IMAGE };

export default writable(HOMEPAGE_HERO_IMAGE);
