export type BarItem = { label: string; value: number; formatted: string };

export type Metric = { label: string; val: string; sub: string };

/**
 * Drops the YAML frontmatter Mycelium stamps above every artifact body and the blank lines under it, as the
 * server does before validation, so a line number in the viewer is the one the API reports.
 */
export function stripFrontmatter(raw: string): string {
	if (raw.startsWith('---\n') || raw.startsWith('---\r\n')) {
		const end = raw.indexOf('\n---', 4);
		if (end >= 0) {
			const cut = raw.indexOf('\n', end + 4);
			return cut >= 0 ? raw.slice(cut + 1).replace(/^[\r\n]+/, '') : '';
		}
	}
	return raw;
}

/** Reads a `chart:bar` fence: `title: X` plus one `label: value` line per bar. */
export function parseBarChart(raw: string): { title: string; items: BarItem[]; max: number } {
	const lines = raw.trim().split('\n');
	let title = '';
	const items: BarItem[] = [];

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
		items.push({ label: key, value: Number.isFinite(num) ? num : 0, formatted: rawVal });
	}

	const max = Math.max(...items.map((i) => i.value), 1);
	return { title, items, max };
}

/** Reads a `chart:metrics` fence: cards split by `---`, each with label, value and sub lines. */
export function parseMetrics(raw: string): Metric[] {
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
		.filter((m): m is Metric => m !== null);
}
