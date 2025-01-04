import { type RequestEvent, redirect } from '@sveltejs/kit';
import { resetGameCookies } from '$lib/cookies';
import { HttpStatus } from '@totocorpsoftwareinc/frontend-toolkit';

export const backToLobby = async ({ cookies }: RequestEvent) => {
	resetGameCookies(cookies);
	redirect(HttpStatus.SEE_OTHER, '/lobby');
};
