/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: {
				primary: '#263037',
				'primary-hover': '#36454f',
				secondary: '#b87333',
				'secondary-hover': '#fff',
				error: '#d92d0f',
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