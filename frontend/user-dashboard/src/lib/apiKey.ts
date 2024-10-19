import ResponseEnvelope from '$lib/responseEnvelope';

export default class ApiKey {
	readonly user: string = '00000000-0000-0000-0000-000000000000';
	readonly key: string = '00000000-0000-0000-0000-000000000000';
	readonly validUntil: Date = new Date();

	constructor(response: ResponseEnvelope) {
		if (response.error()) {
			return;
		}

		// https://stackoverflow.com/questions/455338/how-do-i-check-if-an-object-has-a-key-in-javascript
		if ('user' in response.details && typeof response.details.user === 'string') {
			this.user = response.details.user;
		}

		if ('key' in response.details && typeof response.details.key === 'string') {
			this.key = response.details.key;
		}

		if ('validUntil' in response.details && typeof response.details.validUntil === 'string') {
			this.validUntil = new Date(response.details.validUntil);
		}
	}
}
