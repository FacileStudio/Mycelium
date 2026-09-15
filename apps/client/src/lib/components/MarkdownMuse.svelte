import { Badge, Button, icons, toast } from '@facile/muse';
import { marked, type Tokens, type Token } from 'marked';
import hljs from 'highlight.js';
import 'highlight.js/styles/github-dark.css';
import AlertBlock from './AlertBlock.svelte';
import MermaidBlock from './MermaidBlock.svelte';
import ListItem from './ListItem.svelte';
import { parseAlert, parseBarChart, parseCompare, parseDiff, parseListItem, stripFrontmatter } from './MarkdownMuseUtils';
import { richInline, richBlock } from './MarkdownMuseRich';

function highlightCode(code: string, lang: string): string {
	if (!hljs) return code;
	try {
		return hljs.highlight(code, { language: lang || 'plaintext', ignoreIllegals: true }).value;
	} catch {
		return code;
	}
}

function copyToClipboard(text: string) {
	navigator.clipboard.writeText(text);
	toast.success('Code copied to clipboard');
}

let { content = '' }: { content: string } = $props();

const cleanBody = $derived(stripFrontmatter(content));
const tokens = $derived(marked.lexer(cleanBody) as Token[]);
