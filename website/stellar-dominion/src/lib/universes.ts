import ResponseEnvelope from '$lib/responseEnvelope';
import { buildUrl, safeFetch } from '$lib/api';

export interface ApiUniverse {
	readonly id: string;
	readonly name: string;
}

export default class Universe {
	readonly id: string = '00000000-0000-0000-0000-000000000000';
	readonly name: string = '';

	constructor(response: object) {
		if ('id' in response && typeof response.id === 'string') {
			this.id = response.id;
		}

		if ('name' in response && typeof response.name === 'string') {
			this.name = response.name;
		}
	}

	public toJson(): ApiUniverse {
		return {
			id: this.id,
			name: this.name
		};
	}
}

export async function getUniverses(): Promise<ResponseEnvelope> {
	const url = buildUrl('universes');

	const params = {
		method: 'GET'
	};

	const response = await safeFetch(url, params);
	const jsonContent = await response.json();

	return new ResponseEnvelope(jsonContent);
}
