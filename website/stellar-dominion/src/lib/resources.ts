
export interface Resource {
	readonly id: string;
	readonly amount: number;
};

export function parseResources(data: object[]): Resource[] {
	const out: Resource[] = [];

	for (const maybeResource of data) {
		const hasResource = 'resource' in maybeResource && typeof maybeResource.resource === 'string';
		const hasAmount = 'amount' in maybeResource && typeof maybeResource.amount === 'number';

		if (hasResource && hasAmount) {
			const res: Resource = {
				id: maybeResource.resource as string,
				amount: maybeResource.amount as number
			};

			out.push(res);
		}
	}

	return out;
}