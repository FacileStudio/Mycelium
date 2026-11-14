/** The longest query `GET /artifacts/search` accepts, in characters. */
export const MAX_QUERY = 200;

/** The query the server receives: trimmed and cut to the length it accepts. */
export function normaliseQuery(raw: string): string {
	return raw.trim().slice(0, MAX_QUERY);
}

/** The search route over every tree the caller reads, with the cursor of the next page if any. */
export function artifactSearchPath(q: string, cursor = ''): string {
	const params = new URLSearchParams({ all: 'true', q: normaliseQuery(q) });
	if (cursor) params.set('cursor', cursor);
	return `/artifacts/search?${params}`;
}

export interface Part {
	text: string;
	hit: boolean;
}

function escapeRegExp(word: string): string {
	return word.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

/**
 * Splits `text` around every case-insensitive occurrence of a query word, so a snippet can
 * show what matched without rendering HTML from the server.
 */
export function highlightParts(text: string, query: string): Part[] {
	const words = normaliseQuery(query).toLowerCase().split(/\s+/).filter(Boolean);
	if (!text || words.length === 0) return text ? [{ text, hit: false }] : [];
	const pattern = new RegExp(`(${words.sort((a, b) => b.length - a.length).map(escapeRegExp).join('|')})`, 'gi');
	return text
		.split(pattern)
		.filter(Boolean)
		.map((chunk) => ({ text: chunk, hit: words.includes(chunk.toLowerCase()) }));
}

/** True when a key press lands in a place that takes text, where `/` must type a slash. */
export function isEditableTarget(target: EventTarget | null): boolean {
	if (!target || typeof (target as HTMLElement).tagName !== 'string') return false;
	const el = target as HTMLElement;
	return el.isContentEditable || ['INPUT', 'TEXTAREA', 'SELECT'].includes(el.tagName);
}

/**
 * Hands out tickets so only the newest request may write its answer: a slow search for "enc"
 * must not overwrite the results of "encre" that came back first.
 */
export function latestOnly() {
	let current = 0;
	return {
		next: () => ++current,
		isCurrent: (ticket: number) => ticket === current
	};
}
