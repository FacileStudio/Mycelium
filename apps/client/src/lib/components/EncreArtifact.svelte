<script lang="ts">
	import { Alert, Spinner } from '@facile/muse';
	import { Document, loadEncre, parseMarkdown, validateDocument, validateMarkdown } from '@facile/encre';
	import { locale, t } from '$lib/i18n/index.svelte';
	import { stripFrontmatter } from './artifactSource';
	import { renderDiagram } from './diagram';
	import { toEncreMarkdown } from './encreCompat';
	import { readArtifact, type Reading } from './encreRefusal';
	import { highlight } from './highlight';

	let { content }: { content: string } = $props();

	const stored = $derived(stripFrontmatter(content));
	const body = $derived(
		toEncreMarkdown(stored, {
			note: t('md.note'),
			tip: t('md.tip'),
			important: t('md.important'),
			warning: t('md.warning'),
			caution: t('md.caution'),
			advantages: t('md.advantages'),
			drawbacks: t('md.drawbacks')
		})
	);
	let wasm = $state<'loading' | 'ready' | 'failed'>('loading');

	$effect(() => {
		loadEncre().then(
			() => (wasm = 'ready'),
			() => (wasm = 'failed')
		);
	});

	const parsed = $derived.by((): Reading | null =>
		wasm === 'ready' ? readArtifact(stored, body, { parseMarkdown, validateDocument, validateMarkdown }) : null
	);
</script>

{#snippet source()}
	<pre class="overflow-x-auto rounded-fc-lg bg-fc-surface p-4 font-fc-mono text-fc-xs leading-relaxed whitespace-pre-wrap text-fc-fg">{body}</pre>
{/snippet}

{#if wasm === 'loading'}
	<div class="flex items-center gap-3 text-fc-sm text-fc-fg-muted">
		<Spinner size="sm" label={t('common.loading')} /> {t('artifacts.loadingOne')}
	</div>
{:else if wasm === 'failed' || !parsed || !parsed.doc}
	{#if parsed?.reason}
		<Alert tone="warning" title={t('artifacts.degradedTitle')} class="mb-4">{t('artifacts.degraded', { reason: parsed.reason })}</Alert>
	{/if}
	{@render source()}
{:else}
	{#if parsed.reason}
		<Alert tone="warning" title={t('artifacts.degradedTitle')} class="mb-4">{t('artifacts.degraded', { reason: parsed.reason })}</Alert>
	{/if}
	<svelte:boundary>
		<Document doc={parsed.doc} locale={locale()} {renderDiagram} {highlight} />
		{#snippet failed()}
			{@render source()}
		{/snippet}
	</svelte:boundary>
{/if}
