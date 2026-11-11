import type { EncreNode, MarkdownProblem, MarkdownValidation } from '@facile/encre';

/** The encre calls the viewer needs, passed in so the bun test can hand the wasm module loaded by path. */
export type EncreReader = {
	parseMarkdown(markdown: string): EncreNode;
	validateDocument(doc: EncreNode): string | null;
	validateMarkdown(markdown: string): MarkdownValidation;
};

/** A document ready for `Document`, or null when encre cannot read it, with the reason it is not valid. */
export type Reading = { doc: EncreNode | null; reason: string | null };

type Placed = Partial<Pick<MarkdownProblem, 'line' | 'column'>> & { message: string };

function placed(problem: Placed): string {
	return problem.line ? `line ${problem.line}, column ${problem.column}: ${problem.message}` : problem.message;
}

function refusal(source: string, encre: EncreReader): string | null {
	const verdict = encre.validateMarkdown(source);
	return verdict.valid || verdict.errors.length === 0 ? null : verdict.errors.map(placed).join('; ');
}

function messageOf(error: unknown): string {
	if (error instanceof Error) return placed(error as Error & Placed);
	return String(error);
}

/**
 * Reads the mapped body for the viewer. When encre refuses it, the reason is the one the API gives for the
 * stored body (`line L, column C: <message>`, from `validateMarkdown` on `source`), else the error's own.
 */
export function readArtifact(source: string, mapped: string, encre: EncreReader): Reading {
	let doc: EncreNode;
	try {
		doc = encre.parseMarkdown(mapped);
	} catch (error) {
		return { doc: null, reason: refusal(source, encre) ?? messageOf(error) };
	}
	const invalid = encre.validateDocument(doc);
	return { doc, reason: invalid === null ? null : (refusal(source, encre) ?? invalid) };
}
