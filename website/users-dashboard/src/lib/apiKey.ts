
import ResponseEnvelope from './responseEnvelope';

export default class ApiKey {
	readonly user: string = "00000000-0000-0000-0000-000000000000";
	readonly key: string = "00000000-0000-0000-0000-000000000000";
	readonly validUntil: Date = new Date();

	constructor(response: ResponseEnvelope) {
		if (response.error()) {
			return
		}

		const maybeUser = (response.details as any).user;
		if (typeof maybeUser === "string") {
			this.user = maybeUser;
		}

		const maybeKey = (response.details as any).key;
		if (typeof maybeKey === "string") {
			this.key = maybeKey;
		}

		const maybeValidUntil = (response.details as any).validUntil;
		if (typeof maybeValidUntil === "string") {
			this.validUntil = new Date(maybeValidUntil);
		}
	}
}
