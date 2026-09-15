const MONTHS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

export function sum(values: number[]): number {
	return values.reduce((total, v) => total + (Number.isFinite(v) ? v : 0), 0);
}

export function periodDelta(values: number[], unit: string): string | undefined {
	const half = Math.floor(values.length / 2);
	if (half < 1) return undefined;
	const previous = sum(values.slice(values.length - 2 * half, values.length - half));
	const current = sum(values.slice(values.length - half));
	if (previous <= 0) return undefined;
	const pct = Math.round(((current - previous) / previous) * 100);
	return `${pct >= 0 ? '+' : ''}${pct}% vs previous ${half} ${unit}${half === 1 ? '' : 's'}`;
}

export function columnTotals(rows: number[][]): number[] {
	const width = rows.reduce((max, r) => Math.max(max, r.length), 0);
	const out = new Array(width).fill(0);
	for (const row of rows) {
		for (let i = 0; i < row.length; i++) out[i] += Number.isFinite(row[i]) ? row[i] : 0;
	}
	return out;
}

export function columnActive(rows: number[][]): number[] {
	const width = rows.reduce((max, r) => Math.max(max, r.length), 0);
	const out = new Array(width).fill(0);
	for (const row of rows) {
		for (let i = 0; i < row.length; i++) if (row[i] > 0) out[i] += 1;
	}
	return out;
}

export function hours(seconds: number): number {
	return Math.round((seconds / 3600) * 10) / 10;
}

export function dayKey(d: Date): string {
	const month = String(d.getUTCMonth() + 1).padStart(2, '0');
	const day = String(d.getUTCDate()).padStart(2, '0');
	return `${d.getUTCFullYear()}-${month}-${day}`;
}

export function dayWindow(days: number, end = new Date()): string[] {
	const labels: string[] = [];
	const cursor = Date.UTC(end.getUTCFullYear(), end.getUTCMonth(), end.getUTCDate());
	for (let i = days - 1; i >= 0; i--) labels.push(dayKey(new Date(cursor - i * 86_400_000)));
	return labels;
}

export function bucketByDay(
	entries: { iso: string; weight?: number }[],
	labels: string[]
): number[] {
	const index = new Map(labels.map((l, i) => [l, i]));
	const out = new Array(labels.length).fill(0);
	for (const entry of entries) {
		const d = new Date(entry.iso);
		if (isNaN(d.getTime())) continue;
		const i = index.get(dayKey(d));
		if (i === undefined) continue;
		out[i] += entry.weight ?? 1;
	}
	return out;
}

export function bucketLabel(label: string): string {
	const parts = label.split('-');
	const month = MONTHS[Number(parts[1]) - 1] ?? label;
	if (parts.length >= 3) return `${month} ${Number(parts[2])}`;
	return `${month} ${parts[0]?.slice(2) ?? ''}`;
}