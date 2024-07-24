// https://kit.svelte.dev/docs/advanced-routing#matching

// https://stackoverflow.com/questions/7905929/how-to-test-valid-uuid-guid
const UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-5][0-9a-f]{3}-[089ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

/**
 * @param {string} param
 * @return {param is ('id')}
 * @satisfies {import('@sveltejs/kit').ParamMatcher}
 */
export function match(param) {
	return UUID_REGEX.test(param);
}
