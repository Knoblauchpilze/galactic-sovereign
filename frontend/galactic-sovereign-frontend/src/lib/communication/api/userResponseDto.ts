export class UserResponseDto {
	readonly id: string = '00000000-0000-0000-0000-000000000000';
	readonly email: string = '';
	readonly password: string = '';
	readonly createdAt: Date = new Date();

	constructor(data: object) {
		// https://stackoverflow.com/questions/43894565/cast-object-to-interface-in-typescript
		if ('id' in data && typeof data.id === 'string') {
			this.id = data.id;
		}

		if ('email' in data && typeof data.email === 'string') {
			this.email = data.email;
		}

		if ('password' in data && typeof data.password === 'string') {
			this.password = data.password;
		}

		// https://stackoverflow.com/questions/643782/how-to-check-whether-an-object-is-a-date
		if ('createdAt' in data && typeof data.createdAt === 'string') {
			this.createdAt = new Date(data.createdAt);
		}
	}
}
