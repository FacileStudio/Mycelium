export function formatAge(iso: string): string {
	const then = new Date(iso).getTime();
	if (!iso || isNaN(then)) return '—';
	const s = Math.max(0, Math.floor((Date.now() - then) / 1000));
	if (s < 60) return 'just now';
	if (s < 3600) return `${Math.floor(s / 60)}m ago`;
	if (s < 86400) return `${Math.floor(s / 3600)}h ago`;
	return `${Math.floor(s / 86400)}d ago`;
}

export function formatDuration(seconds: number): string {
	const h = Math.floor(seconds / 3600);
	const m = Math.round((seconds % 3600) / 60);
	if (h > 0) return `${h}h${String(m).padStart(2, '0')}m`;
	return `${m}m`;
}

export function formatSpan(seconds: number): string {
	if (!Number.isFinite(seconds) || seconds <= 0) return 'now';
	const h = Math.floor(seconds / 3600);
	const m = Math.floor((seconds % 3600) / 60);
	if (h >= 24) return `${Math.floor(h / 24)}d ${h % 24}h`;
	if (h > 0) return `${h}h ${m}m`;
	return `${Math.max(1, m)}m`;
}

export function formatCountdown(iso: string): string {
	const target = new Date(iso).getTime();
	if (!iso || isNaN(target)) return 'unknown';
	return formatSpan(Math.floor((target - Date.now()) / 1000));
}

export function formatTokens(n: number): string {
	if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`;
	if (n >= 1_000) return `${(n / 1_000).toFixed(1)}k`;
	return `${n}`;
}

export function formatBytes(n: number): string {
	if (n >= 1_048_576) return `${(n / 1_048_576).toFixed(1)} MB`;
	if (n >= 1024) return `${Math.round(n / 1024)} kB`;
	return `${n} B`;
}

export function formatCost(n: number): string {
	if (n >= 100) return `$${Math.round(n)}`;
	if (n >= 0.01) return `$${n.toFixed(2)}`;
	return `$${n.toFixed(4)}`;
}