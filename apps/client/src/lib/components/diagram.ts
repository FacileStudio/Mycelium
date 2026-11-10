import type { RenderDiagram } from '@facile/encre';
import { theme } from '$lib/theme.svelte';

let rendered = 0;

function dark(): boolean {
	if (theme.mode !== 'system') return theme.mode === 'dark';
	return typeof window !== 'undefined' && window.matchMedia('(prefers-color-scheme: dark)').matches;
}

/**
 * Draws a `mermaid` diagram for encre's `Document` with the app's mermaid, loaded on first use so
 * an artifact without a diagram never pays for it. `d2` and `dot` are refused, so the document
 * shows their source under a language label.
 */
export const renderDiagram: RenderDiagram = async (language, source) => {
	if (language !== 'mermaid') throw new Error(`${language} diagrams are not rendered here`);
	const { default: mermaid } = await import('mermaid');
	mermaid.initialize({
		startOnLoad: false,
		securityLevel: 'strict',
		fontFamily: 'inherit',
		theme: dark() ? 'dark' : 'neutral'
	});
	const { svg } = await mermaid.render(`encre-diagram-${++rendered}`, source);
	return svg;
};
