export interface AlertData {
	isAlert: boolean;
	tone: 'note' | 'tip' | 'important' | 'warning' | 'caution' | 'pros' | 'cons';
	title: string;
	body: string;
}

export function stripFrontmatter(raw: string): string {
	if (raw.startsWith('---\n') || raw.startsWith('---\r\n')) {
		const end = raw.indexOf('\n---', 4);
		if (end >= 0) {
			const cut = raw.indexOf('\n', end + 4);
			return cut >= 0 ? raw.slice(cut + 1) : '';
		}
	}
	return raw;
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

export function parseAlert(raw: string): AlertData {
	const trimmed = raw.trim();
	const match = trimmed.match(/^\[!(NOTE|TIP|IMPORTANT|WARNING|CAUTION|PROS|CONS|ADVANTAGE|DRAWBACK)\]\s*(.*)$/im);
	if (!match) {
		return { isAlert: false, tone: 'note', title: '', body: raw };
	}
	const rawType = match[1].toUpperCase();
	let tone: AlertData['tone'] = 'note';
	let defaultTitle = 'Note';

	if (rawType === 'TIP' || rawType === 'PROS' || rawType === 'ADVANTAGE') {
		tone = rawType === 'PROS' || rawType === 'ADVANTAGE' ? 'pros' : 'tip';
		defaultTitle = rawType === 'PROS' || rawType === 'ADVANTAGE' ? 'Advantages' : 'Tip';
	} else if (rawType === 'IMPORTANT') {
		tone = 'important';
		defaultTitle = 'Important';
	} else if (rawType === 'WARNING') {
		tone = 'warning';
		defaultTitle = 'Warning';
	} else if (rawType === 'CAUTION' || rawType === 'CONS' || rawType === 'DRAWBACK') {
		tone = rawType === 'CONS' || rawType === 'DRAWBACK' ? 'cons' : 'caution';
		defaultTitle = rawType === 'CONS' || rawType === 'DRAWBACK' ? 'Drawbacks' : 'Caution';
	}

	const firstLineExtra = match[2].trim();
	const afterMatch = trimmed.slice(match[0].length).trim();
	const body = firstLineExtra ? `${firstLineExtra}\n${afterMatch}` : afterMatch;

	return {
		isAlert: true,
		tone,
		title: defaultTitle,
		body: body || raw
	};
}

export function parseListItem(text: string) {
	const trimmed = text.trim();
	if (
		trimmed.startsWith('[+]') ||
		trimmed.startsWith('+ ') ||
		trimmed.startsWith('✅') ||
		trimmed.toLowerCase().startsWith('[pro]') ||
		trimmed.toLowerCase().startsWith('[pros]') ||
		trimmed.toLowerCase().startsWith('[advantage]')
	) {
		const clean = trimmed.replace(/^(\[\+\]|\+\s+|✅|\[pro[s]?\]|\[advantage\])\s*/i, '');
		return { kind: 'pro' as const, clean };
	}
	if (
		trimmed.startsWith('[-]') ||
		trimmed.startsWith('- ') ||
		trimmed.startsWith('❌') ||
		trimmed.toLowerCase().startsWith('[con]') ||
		trimmed.toLowerCase().startsWith('[cons]') ||
		trimmed.toLowerCase().startsWith('[drawback]')
	) {
		const clean = trimmed.replace(/^(\[-\]|-\s+|❌|\[con[s]?\]|\[drawback\])\s*/i, '');
		return { kind: 'con' as const, clean };
	}
	if (trimmed.startsWith('[!]') || trimmed.startsWith('⚠️') || trimmed.toLowerCase().startsWith('[warn]')) {
		const clean = trimmed.replace(/^(\[!\]|⚠️|\[warn\])\s*/i, '');
		return { kind: 'warn' as const, clean };
	}
	if (trimmed.startsWith('[?]') || trimmed.startsWith('ℹ️') || trimmed.toLowerCase().startsWith('[info]')) {
		const clean = trimmed.replace(/^(\[\?\]|ℹ️|\[info\])\s*/i, '');
		return { kind: 'info' as const, clean };
	}
	return { kind: 'normal' as const, clean: text };
}

export function parseDiff(text: string) {
	return text.split('\n').map((line) => {
		if (line.startsWith('+') && !line.startsWith('+++')) return { type: 'add' as const, line };
		if (line.startsWith('-') && !line.startsWith('---')) return { type: 'del' as const, line };
		if (line.startsWith('@@')) return { type: 'hunk' as const, line };
		return { type: 'ctx' as const, line };
	});
}

function buildComparisonSection(title: string, items: string[]) {
	const isPro = /advantage|pro|positive|plus|benefit/i.test(title);
	const isCon = /drawback|con|negative|minus|risk|issue/i.test(title);
	return {
		title,
		tone: isPro ? 'pro' as const : isCon ? 'con' as const : 'neutral' as const,
		items
	};
}

export function parseCompare(raw: string) {
	const sections = raw.split(/^###\s+/m).filter((s) => s.trim().length > 0);
	if (sections.length >= 2) {
		return sections.map((section) => {
			const lines = section.trim().split('\n');
			const title = lines[0].trim();
			const items = lines.slice(1).map((l) => l.trim()).filter((l) => l.length > 0);
			return buildComparisonSection(title, items);
		});
	}
	const blocks = raw.split(/---+/).filter((s) => s.trim().length > 0);
	if (blocks.length >= 2) {
		return blocks.map((block) => {
			const lines = block.trim().split('\n');
			const title = lines[0].replace(/^#+\s*/, '').trim();
			const items = lines.slice(1).map((l) => l.trim()).filter((l) => l.length > 0);
			return buildComparisonSection(title || /advantage|pro|positive|plus|benefit/i.test(title) ? 'Advantages' : 'Drawbacks', items);
		});
	}
	return [];
}
