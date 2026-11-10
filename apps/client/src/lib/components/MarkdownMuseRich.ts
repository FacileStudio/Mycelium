import DOMPurify from 'dompurify';
import { renderRich } from './markdownRich';

/**
 * Strips scripts, event handlers and unsafe URLs from rendered markdown, so an
 * artifact written by an agent or a synced machine cannot run code in the page.
 */
export function safeHtml(html: string): string {
	return DOMPurify.sanitize(html, { FORBID_TAGS: ['style', 'form', 'input', 'button', 'textarea', 'select'] });
}

export function richInline(rawText: string): string {
	return safeHtml(renderRich(rawText, true));
}

export function richBlock(rawText: string): string {
	return safeHtml(renderRich(rawText, false));
}
