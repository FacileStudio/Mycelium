import DOMPurify from 'dompurify';
import { marked } from 'marked';

export const toneColors: Record<string, string> = {
	success: 'bg-fc-success/15 text-fc-success ',
	danger: 'bg-fc-danger/15 text-fc-danger ',
	warning: 'bg-fc-warning/15 text-fc-warning ',
	info: 'bg-fc-info/15 text-fc-info ',
	primary: 'bg-fc-primary/15 text-fc-primary ',
	neutral: 'bg-fc-surface text-fc-fg-muted '
};

function applyBadge(tone: string, label: string): string {
	const cls = toneColors[tone.toLowerCase()] || toneColors.neutral;
	return `<span class="inline-flex items-center rounded-full px-2 py-0.5 text-[0.7rem] font-semibold tracking-wide ${cls}">${label}</span>`;
}

function applyRichFormats(processed: string): string {
	processed = processed.replace(
		/\[(badge|tag|chip):(success|danger|warning|info|primary|neutral)\s+([^\]]+)\]/gi,
		(_type, tone, label) => applyBadge(tone, label)
	);
	processed = processed.replace(
		/\[\+\s+([^\]]+)\]/g,
		`<span class="inline-flex items-center gap-1 rounded-full bg-fc-success/15 px-2 py-0.5 text-[0.7rem] font-semibold text-fc-success"><span class="font-bold text-fc-success">+</span> $1</span>`
	);
	processed = processed.replace(
		/\[-\s+([^\]]+)\]/g,
		`<span class="inline-flex items-center gap-1 rounded-full bg-fc-danger/15 px-2 py-0.5 text-[0.7rem] font-semibold text-fc-danger"><span class="font-bold text-fc-danger">-</span> $1</span>`
	);
	processed = processed.replace(
		/\[status:(active|done|ready|success)\]/gi,
		`<span class="inline-flex items-center gap-1.5 rounded-full bg-fc-success/10 px-2 py-0.5 text-[0.68rem] font-semibold text-fc-success"><span class="size-1.5 rounded-full bg-fc-success animate-pulse"></span>$1</span>`
	);
	processed = processed.replace(
		/\[status:(pending|wip|running|loading)\]/gi,
		`<span class="inline-flex items-center gap-1.5 rounded-full bg-fc-warning/10 px-2 py-0.5 text-[0.68rem] font-semibold text-fc-warning"><span class="size-1.5 rounded-full bg-fc-warning"></span>$1</span>`
	);
	processed = processed.replace(
		/\[status:(error|failed|stopped|dead)\]/gi,
		`<span class="inline-flex items-center gap-1.5 rounded-full bg-fc-danger/10 px-2 py-0.5 text-[0.68rem] font-semibold text-fc-danger"><span class="size-1.5 rounded-full bg-fc-danger"></span>$1</span>`
	);
	processed = processed.replace(
		/==([^=]+)==/g,
		`<mark class="rounded bg-fc-warning/20 px-1 py-0.5 font-medium text-fc-fg">$1</mark>`
	);
	return processed;
}

/**
 * Strips scripts, event handlers and unsafe URLs from rendered markdown, so an
 * artifact written by an agent or a synced machine cannot run code in the page.
 */
export function safeHtml(html: string): string {
	return DOMPurify.sanitize(html, { FORBID_TAGS: ['style', 'form', 'input', 'button', 'textarea', 'select'] });
}

export function richInline(rawText: string): string {
	return safeHtml(marked.parseInline(applyRichFormats(rawText)) as string);
}

export function richBlock(rawText: string): string {
	return safeHtml(marked.parse(applyRichFormats(rawText)) as string);
}
