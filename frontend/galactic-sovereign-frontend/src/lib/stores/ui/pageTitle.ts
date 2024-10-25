import { writable } from 'svelte/store';

const HOMEPAGE_TITLE: string = 'Stellar Dominion';

export { HOMEPAGE_TITLE };

export default writable(HOMEPAGE_TITLE);
