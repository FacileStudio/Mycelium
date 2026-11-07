import { formatNumber } from '$lib/i18n/format';
import { t } from '$lib/i18n/index.svelte';

function unitText(value: number, unit: string, digits = 0): string {
	return formatNumber(value, { style: 'unit', unit, unitDisplay: 'narrow', maximumFractionDigits: digits });
}

export function formatAge(iso: string): string {
	const then = new Date(iso).getTime();
	if (!iso || isNaN(then)) return '—';
	const s = Math.max(0, Math.floor((Date.now() - then) / 1000));
	if (s < 60) return t('time.justNow');
	if (s < 3600) return t('time.minutesAgo', { count: Math.floor(s / 60) });
	if (s < 86400) return t('time.hoursAgo', { count: Math.floor(s / 3600) });
	return t('time.daysAgo', { count: Math.floor(s / 86400) });
}

export function formatDuration(seconds: number): string {
	const h = Math.floor(seconds / 3600);
	const m = Math.round((seconds % 3600) / 60);
	if (h > 0) return `${unitText(h, 'hour')} ${unitText(m, 'minute')}`;
	return unitText(m, 'minute');
}

export function formatSpan(seconds: number): string {
	if (!Number.isFinite(seconds) || seconds <= 0) return t('time.now');
	const h = Math.floor(seconds / 3600);
	const m = Math.floor((seconds % 3600) / 60);
	if (h >= 24) return `${unitText(Math.floor(h / 24), 'day')} ${unitText(h % 24, 'hour')}`;
	if (h > 0) return `${unitText(h, 'hour')} ${unitText(m, 'minute')}`;
	return unitText(Math.max(1, m), 'minute');
}

export function formatCountdown(iso: string): string {
	const target = new Date(iso).getTime();
	if (!iso || isNaN(target)) return t('time.unknown');
	return formatSpan(Math.floor((target - Date.now()) / 1000));
}

export function formatTokens(n: number): string {
	return formatNumber(n, { notation: 'compact', maximumFractionDigits: 1 });
}

export function formatBytes(n: number): string {
	if (n >= 1_048_576) return unitText(n / 1_048_576, 'megabyte', 1);
	if (n >= 1024) return unitText(Math.round(n / 1024), 'kilobyte');
	return unitText(n, 'byte');
}

export function formatCost(n: number): string {
	const digits = n >= 100 ? 0 : n >= 0.01 ? 2 : 4;
	return formatNumber(n, { style: 'currency', currency: 'USD', minimumFractionDigits: digits, maximumFractionDigits: digits });
}

export function parseBarChart(raw: string) {
	const lines = raw.trim().split('\n');
	let title = '';
	const items: { label: string; value: number; formatted: string }[] = [];

	for (const line of lines) {
		const trimmed = line.trim();
		if (!trimmed) continue;
		const colon = trimmed.indexOf(':');
		if (colon < 0) continue;
		const key = trimmed.slice(0, colon).trim();
		const rawVal = trimmed.slice(colon + 1).trim();
		if (key.toLowerCase() === 'title') {
			title = rawVal;
			continue;
		}
		const numMatch = rawVal.match(/[\d.,]+/);
		const num = numMatch ? parseFloat(numMatch[0].replace(',', '.')) : 0;
		items.push({ label: key, value: num, formatted: rawVal });
	}

	const max = Math.max(...items.map((i) => i.value), 1);
	return { title, items, max };
}

export function parseMetrics(raw: string) {
	const blocks = raw.trim().split(/---+/);
	return blocks
		.map((block) => {
			const lines = block.trim().split('\n');
			let label = '';
			let val = '';
			let sub = '';
			for (const line of lines) {
				const colon = line.indexOf(':');
				if (colon < 0) continue;
				const k = line.slice(0, colon).trim().toLowerCase();
				const v = line.slice(colon + 1).trim();
				if (k === 'label' || k === 'title') label = v;
				else if (k === 'val' || k === 'value') val = v;
				else if (k === 'sub' || k === 'desc' || k === 'subtitle') sub = v;
			}
			if (!val && !label) return null;
			return { label, val, sub };
		})
		.filter((m): m is { label: string; val: string; sub: string } => m !== null);
}