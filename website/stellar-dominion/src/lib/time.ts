const MS_IN_ONE_SECOND = 1000;
const SECOND_IN_A_MINUTE = 60;
const MINUTE_IN_AN_HOUR = 60;
const HOUR_IN_A_DAY = 24;

const MS_IN_A_MINUTE = SECOND_IN_A_MINUTE * MS_IN_ONE_SECOND;
const MS_IN_AN_HOUR = MINUTE_IN_AN_HOUR * MS_IN_A_MINUTE;
const MS_IN_A_DAY = HOUR_IN_A_DAY * MS_IN_AN_HOUR;

function positiveNumberToStringWithSuffix(value: number, suffix: string): string {
	let out = '';

	out += value.toString();
	out += ' ';
	out += suffix;
	if (value !== 1) {
		out += 's';
	}

	return out;
}

function addSpaceIfNotEmpty(str: string): string {
	if (str.length === 0) {
		return str;
	}
	return str + ' ';
}

export function msToTimeString(timeInMs: number): string {
	if (timeInMs < 0) {
		return '-' + msToTimeString(-timeInMs);
	} else if (timeInMs == 0) {
		return '0s';
	}

	const days = Math.floor(timeInMs / MS_IN_A_DAY);
	const hours = Math.floor((timeInMs % MS_IN_A_DAY) / MS_IN_AN_HOUR);
	const minutes = Math.floor((timeInMs % MS_IN_AN_HOUR) / MS_IN_A_MINUTE);
	const seconds = Math.floor((timeInMs % MS_IN_A_MINUTE) / MS_IN_ONE_SECOND);

	let out = '';
	if (days !== 0) {
		out += positiveNumberToStringWithSuffix(days, 'day');
	}
	if (hours !== 0) {
		out = addSpaceIfNotEmpty(out);
		out += positiveNumberToStringWithSuffix(hours, 'hour');
	}
	if (minutes !== 0) {
		out = addSpaceIfNotEmpty(out);
		out += positiveNumberToStringWithSuffix(minutes, 'minute');
	}
	if (seconds !== 0) {
		out = addSpaceIfNotEmpty(out);
		out += positiveNumberToStringWithSuffix(seconds, 'second');
	}
	if (out.length === 0) {
		out = 'Almost done';
	}

	return out;
}

export function msToTimeStringOrFinished(timeInMs: number): string {
	if (timeInMs < 0) {
		return 'Finished';
	}

	return msToTimeString(timeInMs);
}
