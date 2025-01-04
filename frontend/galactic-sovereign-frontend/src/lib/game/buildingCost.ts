export function computeCostForLevel(cost: number, progress: number, level: number): number {
	return Math.floor(cost * Math.pow(progress, level - 1));
}
