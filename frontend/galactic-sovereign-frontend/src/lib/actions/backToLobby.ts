import { type RequestEvent, redirect } from '@sveltejs/kit';
import { resetGameCookies } from '$lib/cookies';

export const backToLobby = async ({ cookies }: RequestEvent) => {
	resetGameCookies(cookies);
	redirect(303, '/lobby');
};
