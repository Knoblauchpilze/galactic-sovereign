/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: {
				// Palette
				primary: '#263037',
				'primary-hover': '#36454f',
				secondary: '#b87333',
				'secondary-hover': '#fff',

				// State
				enabled: "#45d90f",
				disabled: "#d92d0f",
				error: '#d92d0f',

				// Miscellaneous
				overlay: '#0005',
			},
			backgroundImage: {
				homepage: "url('$lib/assets/background.webp')",
				overview: "url('$lib/assets/overview.webp')"
			}
		}
	},
	plugins: []
};
