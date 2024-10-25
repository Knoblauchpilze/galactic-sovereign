import { writable } from 'svelte/store';

const HOMEPAGE_TITLE: string = 'Galactic Sovereign';

export { HOMEPAGE_TITLE };

export default writable(HOMEPAGE_TITLE);
