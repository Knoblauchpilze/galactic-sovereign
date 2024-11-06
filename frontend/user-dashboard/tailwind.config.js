/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: {
				// Palette
				primary: '#000',
				'primary-hover': '#222',
				'primary-selected': '#444',
				secondary: '#fff',
				'secondary-hover': '#fff',

				// State
				error: '#d92d0f',

				// Miscellaneous
				overlay: '#0b0b45'
			},
			backgroundImage: {
				homepage: "url('$lib/assets/background.webp')",
				overview: "url('$lib/assets/overview.webp')"
			}
		}
	},
	plugins: []
};
