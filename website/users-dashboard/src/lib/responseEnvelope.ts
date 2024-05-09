
export default class ResponseEnvelope {
	readonly RequestId: string;
	readonly Status: string;
	readonly Details: object;

	public constructor(response: {RequestId: string, Status: string, Details: object}) {
		this.RequestId = response.RequestId;
		this.Status = response.Status;
		this.Details = response.Details;
	}

	public success(): boolean {
		return this.Status === "SUCCESS";
	}

	public error(): boolean {
		return !this.success();
	}
}


