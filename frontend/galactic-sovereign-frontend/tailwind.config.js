/** @type {import('tailwindcss').Config} */
export default {
	content: [
		'./src/**/*.{html,js,svelte,ts}',
		'./node_modules/@totocorpsoftwareinc/frontend-toolkit/**/*.svelte'
	],
	theme: {
		extend: {
			colors: {
				// Palette
				primary: '#263037',
				'primary-hover': '#36454f',
				'primary-selected': '#535a5e',
				secondary: '#b87333',
				'secondary-hover': '#fff',

				// State
				enabled: '#2a7a0c',
				'enabled-hover': '#45d90f',
				disabled: '#751c0d',
				'disabled-hover': '#d92d0f',
				error: '#d92d0f',

				// Miscellaneous
				overlay: '#0005'
			},
			backgroundImage: {
				homepage: "url('$lib/assets/background.webp')",
				overview: "url('$lib/assets/overview.webp')"
			}
		}
	},
	plugins: []
};
