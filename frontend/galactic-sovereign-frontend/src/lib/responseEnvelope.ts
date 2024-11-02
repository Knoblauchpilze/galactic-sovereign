import { error, redirect } from '@sveltejs/kit';

export class ResponseEnvelope {
	readonly requestId: string;
	readonly status: string;
	readonly details: object;

	constructor(response: { requestId: string; status: string; details: object }) {
		this.requestId = response.requestId;
		this.status = response.status;
		this.details = response.details;
	}

	public success(): boolean {
		return this.status === 'SUCCESS';
	}

	public error(): boolean {
		return !this.success();
	}

	public failureMessage(): string {
		if (typeof this.details === 'string') {
			return this.details;
		}

		return 'Unexpected error';
	}

	public failureReason(): ApiFailureReason {
		if (!this.error()) {
			return ApiFailureReason.NONE;
		}

		switch (this.failureMessage()) {
			case 'API key expired':
				return ApiFailureReason.API_KEY_EXPIRED;
			default:
				return ApiFailureReason.UNKNOWN_ERROR;
		}
	}

	public getDetails(): object {
		return this.details;
	}
}

export function createFailedResponseEnvelope(details: object): ResponseEnvelope {
	return new ResponseEnvelope({
		requestId: '00000000-0000-0000-0000-000000000000',
		status: 'ERROR',
		details: details
	});
}

export function createEmptySuccessResponseEnvelope(): ResponseEnvelope {
	return new ResponseEnvelope({
		requestId: '00000000-0000-0000-0000-000000000000',
		status: 'SUCCESS',
		details: 'No content' as unknown as object
	});
}

export enum ApiFailureReason {
	NONE = 0,
	UNKNOWN_ERROR = 1,
	API_KEY_EXPIRED = 2
}

export function analayzeResponseEnvelopAndRedirectIfNeeded(response: ResponseEnvelope) {
	if (!response.error()) {
		return;
	}

	const reason = response.failureReason();

	switch (reason) {
		case ApiFailureReason.API_KEY_EXPIRED:
			redirect(303, '/login');
	}

	error(404, { message: response.failureMessage() });
}
