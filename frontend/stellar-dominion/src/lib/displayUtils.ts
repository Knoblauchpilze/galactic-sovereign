export function floorToInteger(value: number): number {
	return Math.floor(value);
}

interface abbrevation {
	readonly threshold: number;
	readonly divisor: number;
	readonly suffix: string;
}

export const ABBREVATIONS: abbrevation[] = [
	{
		threshold: 10_000,
		divisor: 1,
		suffix: ''
	},
	{
		threshold: 1_000_000,
		divisor: 1_000,
		suffix: 'k'
	},
	{
		threshold: 1_000_000_000,
		divisor: 1_000_000,
		suffix: 'M'
	}
];

export function toFlooredShortString(value: number): string {
	let id = 0;
	while (id < ABBREVATIONS.length && ABBREVATIONS[id].threshold <= value) {
		++id;
	}

	const abbrevation = ABBREVATIONS[Math.min(id, ABBREVATIONS.length - 1)];
	const out = value / abbrevation.divisor;

	const withTwoSignificantDigits = Math.floor(out * 10) / 10;

	return withTwoSignificantDigits.toString() + abbrevation.suffix;
}
