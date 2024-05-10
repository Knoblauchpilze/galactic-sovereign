
import ResponseEnvelope from './responseEnvelope';
import {buildUrl, safeFetch} from './api';

export default class User {
	readonly id: string = "00000000-0000-0000-0000-000000000000";
	readonly email: string = "";
	readonly password: string = "";
	readonly createdAt: Date = new Date();

	constructor(response: ResponseEnvelope) {
		// https://stackoverflow.com/questions/43894565/cast-object-to-interface-in-typescript
		if (response.error()) {
			return
		}

		const maybeId = (response.details as any).id;
		if (typeof maybeId === "string") {
			this.id = maybeId;
		}

		const maybeEmail = (response.details as any).email;
		if (typeof maybeId === "string") {
			this.email = maybeEmail;
		}

		const maybePassword = (response.details as any).password;
		if (typeof maybeId === "string") {
			this.password = maybePassword;
		}

		// https://stackoverflow.com/questions/643782/how-to-check-whether-an-object-is-a-date
		const maybeCreatedAt = (response.details as any).createdAt;
		if (typeof maybeCreatedAt === "string") {
			this.createdAt = new Date(maybeCreatedAt);
		}
	}
}

export async function getUser(apiKey: string, id: string): Promise<ResponseEnvelope> {
	const url = buildUrl(id);

	const params = {
		method: 'GET',
		headers: {
			'X-Api-Key' : apiKey,
		}
	};

	const response = await safeFetch(url, params);
	const jsonContent = await response.json();

	return new ResponseEnvelope(jsonContent);
}
