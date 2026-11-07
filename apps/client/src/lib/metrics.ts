import { formatDate, formatNumber } from '$lib/i18n/format';
import { locale, t } from '$lib/i18n/index.svelte';

export function sum(values: number[]): number {
	return values.reduce((total, v) => total + (Number.isFinite(v) ? v : 0), 0);
}

export function periodDelta(values: number[], unit: 'hour' | 'day' | 'week' | 'month'): string | undefined {
	const half = Math.floor(values.length / 2);
	if (half < 1) return undefined;
	const previous = sum(values.slice(values.length - 2 * half, values.length - half));
	const current = sum(values.slice(values.length - half));
	if (previous <= 0) return undefined;
	const pct = Math.round(((current - previous) / previous) * 100);
	const span = formatNumber(half, { style: 'unit', unit, unitDisplay: 'long' });
	return t('metrics.delta', { pct: `${pct >= 0 ? '+' : ''}${formatNumber(pct)}`, span });
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
	const year = Number(parts[0]);
	const month = Number(parts[1]) - 1;
	if (parts.length >= 3) {
		return formatDate(Date.UTC(year, month, Number(parts[2])), { month: 'short', day: 'numeric', timeZone: 'UTC' });
	}
	return formatDate(Date.UTC(year, month, 1), { month: 'short', year: '2-digit', timeZone: 'UTC' });
}

function unitText(value: number, unit: string, digits = 0): string {
	return formatNumber(value, { style: 'unit', unit, unitDisplay: 'narrow', maximumFractionDigits: digits });
}

export function formatDuration(seconds: number): string {
	const h = Math.floor(seconds / 3600);
	const m = Math.floor((seconds % 3600) / 60);
	const s = seconds % 60;
	const parts = [
		h > 0 ? unitText(h, 'hour') : '',
		h > 0 || m > 0 ? (m > 0 ? unitText(m, 'minute') : '') : '',
		h > 0 || m > 0 ? (s > 0 ? unitText(s, 'second') : '') : unitText(s, 'second')
	];
	return parts.filter(Boolean).join(' ');
}

export function formatTokens(tokens: number): string {
	return formatNumber(tokens, { notation: 'compact', maximumFractionDigits: 1 });
}

export function formatCost(cost: number): string {
	const digits = cost >= 100 ? 0 : cost >= 0.01 ? 2 : 4;
	return formatNumber(cost, {
		style: 'currency',
		currency: 'USD',
		minimumFractionDigits: digits,
		maximumFractionDigits: digits
	});
}

export function formatBytes(bytes: number): string {
	const scales: [number, string][] = [[1_000_000_000, 'gigabyte'], [1_000_000, 'megabyte'], [1_000, 'kilobyte']];
	for (const [size, unit] of scales) {
		if (bytes >= size) return unitText(bytes / size, unit, 1);
	}
	return unitText(bytes, 'byte');
}

export function formatAge(timestamp: number): string {
	const days = Math.floor((Date.now() - timestamp) / 86_400_000);
	const relative = new Intl.RelativeTimeFormat(locale(), { numeric: 'auto' });
	if (days < 7) return relative.format(-days, 'day');
	if (days < 30) return relative.format(-Math.floor(days / 7), 'week');
	if (days < 365) return relative.format(-Math.floor(days / 30), 'month');
	return relative.format(-Math.floor(days / 365), 'year');
}
