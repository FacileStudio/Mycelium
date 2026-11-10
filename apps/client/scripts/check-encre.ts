import { stripFrontmatter } from '../src/lib/components/artifactSource';
import { toEncreMarkdown, type AlertTitles } from '../src/lib/components/encreCompat';
import { loadEncre, parseMarkdown, serializeMarkdown, validateDocument } from '../node_modules/@facile/encre/src/lib/wasm';
import type { EncreNode } from '../node_modules/@facile/encre/src/lib/schema';

const wasm = await Bun.file(new URL('../node_modules/@facile/encre/src/lib/encre-wasm/encre_wasm_bg.wasm', import.meta.url)).arrayBuffer();
await loadEncre(wasm);

const titles: AlertTitles = { note: 'Note', tip: 'Tip', important: 'Important', warning: 'Warning', caution: 'Caution', advantages: 'Advantages', drawbacks: 'Drawbacks' };

function render(markdown: string): { mapped: string; doc: EncreNode; reason: string | null } {
	const mapped = toEncreMarkdown(stripFrontmatter(markdown), titles);
	const doc = parseMarkdown(mapped);
	return { mapped, doc, reason: validateDocument(doc) };
}

const blocks = (doc: EncreNode) => doc.content ?? [];
const attrs = (node: EncreNode) => node.attrs ?? {};

const legacy = [
	'---',
	'title: "Weekly"',
	'machine: "ruche"',
	'---',
	'',
	'# Weekly',
	'',
	'```chart:bar',
	'title: Sales',
	'Jan: 10',
	'Feb: 25 k',
	'```',
	'',
	'```chart:metrics',
	'label: Repos',
	'val: 42',
	'---',
	'label: Active',
	'val: 16',
	'sub: In production',
	'```',
	'',
	'> [!WARNING]',
	'> Mind the **gap**',
	'',
	'```mermaid',
	'flowchart LR',
	'  A-->B',
	'```',
	'',
	'A [badge:success kept] line.',
	''
].join('\n');

const cases: [string, () => boolean][] = [
	[
		'legacy fences and alerts become encre blocks, the rest stays',
		() => {
			const { doc, reason } = render(legacy);
			const types = blocks(doc).map((node) => node.type);
			const chart = attrs(blocks(doc)[1]);
			const active = attrs(blocks(doc)[3]);
			const callout = blocks(doc)[4];
			const text = blocks(doc)[6].content?.map((node) => node.text).join('') ?? '';
			return (
				reason === null &&
				types.join() === 'heading,chart,stat,stat,callout,codeBlock,paragraph' &&
				JSON.stringify(chart.labels) === '["Jan","Feb"]' &&
				JSON.stringify(chart.series) === '[{"data":[10,25],"name":"Sales"}]' &&
				active.delta === 'In production' &&
				attrs(callout).tone === 'warning' &&
				attrs(callout).title === 'Warning' &&
				attrs(blocks(doc)[5]).language === 'mermaid' &&
				text === 'A [badge:success kept] line.'
			);
		}
	],
	[
		'an encre chart fence is left alone',
		() => {
			const source = '```chart\n{"kind": "line", "labels": ["a"], "series": [{"name": "s", "data": [1]}]}\n```\n';
			const { mapped, doc } = render(source);
			return mapped === source && blocks(doc)[0].type === 'chart' && attrs(blocks(doc)[0]).kind === 'line';
		}
	],
	[
		'a callout fence outgrows the backtick runs inside it',
		() => {
			const { mapped, doc } = render('> [!NOTE]\n> four ```` ticks\n');
			return mapped.startsWith('`````callout tone="info" title="Note"\n') && blocks(doc)[0].type === 'callout';
		}
	],
	[
		'the mapped document round-trips through the canonical writer',
		() => {
			const { doc } = render(legacy);
			const again = parseMarkdown(serializeMarkdown(doc));
			return JSON.stringify(again) === JSON.stringify(doc) && validateDocument(again) === null;
		}
	],
	[
		'an invalid body still parses and names its reason',
		() => {
			const { doc, reason } = render('# Old\n\na\u0001b\n');
			return blocks(doc).length === 2 && typeof reason === 'string' && reason.includes('forbidden');
		}
	]
];

const failures = cases.filter(([, ok]) => !ok());
for (const [name] of failures) console.error(`encre: ${name}`);
if (failures.length > 0) process.exit(1);
console.log(`encre: ${cases.length} cases map legacy artifacts onto the document format`);
