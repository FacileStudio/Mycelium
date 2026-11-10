import { renderRich } from '../src/lib/components/markdownRich';

const cases: [string, boolean, (html: string) => boolean][] = [
	['a <script>alert(1)</script> b', true, (html) => !html.includes('<script') && html.includes('&lt;script&gt;')],
	['<div onclick="x()">raw</div>', false, (html) => !html.includes('<div') && html.includes('&lt;div')],
	['[badge:info <b>hi</b>] ==<i>m</i>==', true, (html) => html.includes('<span class=') && !html.includes('<b>') && !html.includes('<i>')],
	['`<i>` [status:done]', true, (html) => html.includes('<code>&lt;i&gt;</code>') && html.includes('done</span>')]
];

const failures = cases.filter(([source, inline, ok]) => !ok(renderRich(source, inline)));
for (const [source, inline] of failures) console.error(`markdown: ${JSON.stringify(source)} -> ${renderRich(source, inline)}`);
if (failures.length > 0) process.exit(1);
console.log(`markdown: ${cases.length} cases render raw HTML as text`);
