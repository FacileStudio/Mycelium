import { Marked } from 'marked';

const toneColors: Record<string, string> = {
	success: 'bg-fc-success/15 text-fc-success ',
	danger: 'bg-fc-danger/15 text-fc-danger ',
	warning: 'bg-fc-warning/15 text-fc-warning ',
	info: 'bg-fc-info/15 text-fc-info ',
	primary: 'bg-fc-primary/15 text-fc-primary ',
	neutral: 'bg-fc-surface text-fc-fg-muted '
};

const pill = 'inline-flex items-center rounded-full px-2 py-0.5 text-[0.7rem] font-semibold';
const statusPill = 'inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-[0.68rem] font-semibold';

function escapeHtml(text: string): string {
	return text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}

const md = new Marked({ renderer: { html: ({ text }) => escapeHtml(text) } });

const rules: [RegExp, (...groups: string[]) => string][] = [
	[
		/\[(badge|tag|chip):(success|danger|warning|info|primary|neutral)\s+([^\]]+)\]/gi,
		(_type, tone, label) =>
			`<span class="${pill} tracking-wide ${toneColors[tone.toLowerCase()] || toneColors.neutral}">${escapeHtml(label)}</span>`
	],
	[
		/\[\+\s+([^\]]+)\]/g,
		(label) =>
			`<span class="${pill} gap-1 bg-fc-success/15 text-fc-success"><span class="font-bold text-fc-success">+</span> ${escapeHtml(label)}</span>`
	],
	[
		/\[-\s+([^\]]+)\]/g,
		(label) =>
			`<span class="${pill} gap-1 bg-fc-danger/15 text-fc-danger"><span class="font-bold text-fc-danger">-</span> ${escapeHtml(label)}</span>`
	],
	[
		/\[status:(active|done|ready|success)\]/gi,
		(label) =>
			`<span class="${statusPill} bg-fc-success/10 text-fc-success"><span class="size-1.5 rounded-full bg-fc-success animate-pulse"></span>${label}</span>`
	],
	[
		/\[status:(pending|wip|running|loading)\]/gi,
		(label) =>
			`<span class="${statusPill} bg-fc-warning/10 text-fc-warning"><span class="size-1.5 rounded-full bg-fc-warning"></span>${label}</span>`
	],
	[
		/\[status:(error|failed|stopped|dead)\]/gi,
		(label) =>
			`<span class="${statusPill} bg-fc-danger/10 text-fc-danger"><span class="size-1.5 rounded-full bg-fc-danger"></span>${label}</span>`
	],
	[
		/==([^=]+)==/g,
		(text) => `<mark class="rounded bg-fc-warning/20 px-1 py-0.5 font-medium text-fc-fg">${escapeHtml(text)}</mark>`
	]
];

const slotRef = /\uE000(\d+)\uE001/g;

/**
 * Renders markdown with the rich formats (badges, status pills, ==marks==).
 * Raw HTML in the source is shown as text, never rendered: artifacts are
 * markdown only. The rich formats are swapped for placeholders before parsing
 * so marked never sees them as raw HTML.
 */
export function renderRich(rawText: string, inline: boolean): string {
	const slots: string[] = [];
	let text = rawText;
	for (const [pattern, build] of rules) {
		text = text.replace(pattern, (_match, ...groups) => {
			slots.push(build(...(groups.slice(0, -2) as string[])));
			return `\uE000${slots.length - 1}\uE001`;
		});
	}
	const html = (inline ? md.parseInline(text) : md.parse(text)) as string;
	return html.replace(slotRef, (_m, i) => slots[Number(i)] ?? '');
}
