
export default class ResponseEnvelope {
	readonly requestId: string;
	readonly status: string;
	readonly details: object;

	constructor(response: {requestId: string, status: string, details: object}) {
		this.requestId = response.requestId;
		this.status = response.status;
		this.details = response.details;
	}

	public success(): boolean {
		return this.status === "SUCCESS";
	}

	public error(): boolean {
		return !this.success();
	}

	public failureReason(): string {
		if (typeof this.details === "string") {
			return this.details;
		}

		return "Unexpected error";
	}
}

export function createFailedResponseEnvelope(details: object): ResponseEnvelope {
	return new ResponseEnvelope({
		requestId: "00000000-0000-0000-0000-000000000000",
		status: "ERROR",
		details: details,
	});
}


