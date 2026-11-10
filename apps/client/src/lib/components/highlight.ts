import type { Highlight } from '@facile/encre';
import DOMPurify from 'dompurify';
import hljs from 'highlight.js';

function escape(code: string): string {
	return code.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
}

/**
 * Colours a code block for encre's `Document`, which inserts the HTML as it is: a language
 * highlight.js knows is highlighted, any other is escaped, and both go through DOMPurify.
 */
export const highlight: Highlight = (code, language) =>
	DOMPurify.sanitize(hljs.getLanguage(language) ? hljs.highlight(code, { language }).value : escape(code));
