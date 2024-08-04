import ResponseEnvelope from '$lib/responseEnvelope';
import { buildUrl, safeFetch } from '$lib/api';

export interface ApiResource {
	readonly id: string;
	readonly name: string;
}

export class Resource {
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

	// https://stackoverflow.com/questions/65512526/cannot-stringify-arbitrary-non-pojos-and-invalid-prop-type-check-failed-for
	public toJson(): ApiResource {
		return {
			id: this.id,
			name: this.name
		};
	}
}

export async function getResources(apiKey: string): Promise<ResponseEnvelope> {
	const url = buildUrl('resources');

	const params = {
		method: 'GET',
		headers: {
			'X-Api-Key': apiKey
		}
	};

	const response = await safeFetch(url, params);
	const jsonContent = await response.json();

	return new ResponseEnvelope(jsonContent);
}

export function responseToResourcesArray(response: ResponseEnvelope): Resource[] {
	if (response.error()) {
		return [];
	}

	const details = response.getDetails();
	if (!Array.isArray(details)) {
		return [];
	}

	return details.map((maybeResource) => new Resource(maybeResource));
}

export interface PlanetResource {
	readonly id: string;
	readonly amount: number;
}

export function parseResources(data: object[]): PlanetResource[] {
	const out: PlanetResource[] = [];

	for (const maybeResource of data) {
		const hasResource = 'resource' in maybeResource && typeof maybeResource.resource === 'string';
		const hasAmount = 'amount' in maybeResource && typeof maybeResource.amount === 'number';

		if (hasResource && hasAmount) {
			const res: PlanetResource = {
				id: maybeResource.resource as string,
				amount: maybeResource.amount as number
			};

			out.push(res);
		}
	}

	return out;
}

export interface UiResource {
	readonly name: string;
	readonly amount: number;
}

export function mapPlanetResourcesToApiResources(
	planetResources: PlanetResource[],
	apiResources: ApiResource[]
): UiResource[] {
	return apiResources.map((apiResource) => {
		const maybeResource = planetResources.find((r) => r.id === apiResource.id);
		if (maybeResource === undefined) {
			return {
				name: apiResource.name,
				amount: 0
			};
		} else {
			return {
				name: apiResource.name,
				amount: maybeResource.amount
			};
		}
	});
}
