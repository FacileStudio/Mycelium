import { t } from '$lib/i18n/index.svelte';
export interface AlertData {
	isAlert: boolean;
	tone: 'note' | 'tip' | 'important' | 'warning' | 'caution' | 'pros' | 'cons';
	title: string;
	body: string;
}

export { parseBarChart, parseMetrics, stripFrontmatter } from './artifactSource';

export function parseAlert(raw: string): AlertData {
	const trimmed = raw.trim();
	const match = trimmed.match(/^\[!(NOTE|TIP|IMPORTANT|WARNING|CAUTION|PROS|CONS|ADVANTAGE|DRAWBACK)\]\s*(.*)$/im);
	if (!match) {
		return { isAlert: false, tone: 'note', title: '', body: raw };
	}
	const rawType = match[1].toUpperCase();
	let tone: AlertData['tone'] = 'note';
	let defaultTitle = t('md.note');

	if (rawType === 'TIP' || rawType === 'PROS' || rawType === 'ADVANTAGE') {
		tone = rawType === 'PROS' || rawType === 'ADVANTAGE' ? 'pros' : 'tip';
		defaultTitle = rawType === 'PROS' || rawType === 'ADVANTAGE' ? t('md.advantages') : t('md.tip');
	} else if (rawType === 'IMPORTANT') {
		tone = 'important';
		defaultTitle = t('md.important');
	} else if (rawType === 'WARNING') {
		tone = 'warning';
		defaultTitle = t('md.warning');
	} else if (rawType === 'CAUTION' || rawType === 'CONS' || rawType === 'DRAWBACK') {
		tone = rawType === 'CONS' || rawType === 'DRAWBACK' ? 'cons' : 'caution';
		defaultTitle = rawType === 'CONS' || rawType === 'DRAWBACK' ? t('md.drawbacks') : t('md.caution');
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
