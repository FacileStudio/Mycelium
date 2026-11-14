<script lang="ts">
	import { untrack } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { Alert, Button, EmptyState, Spinner, icons } from '@facile/muse';
	import ArtifactCard from '$lib/components/ArtifactCard.svelte';
	import ArtifactSearchBar from '$lib/components/ArtifactSearchBar.svelte';
	import { ArtifactFeed } from '$lib/artifactFeed.svelte';
	import { normaliseQuery } from '$lib/artifactSearch';
	import { t } from '$lib/i18n/index.svelte';

	const DEBOUNCE_MS = 250;

	const feed = new ArtifactFeed();
	let query = $state(page.url.searchParams.get('q') ?? '');
	let first = true;

	function syncUrl(q: string) {
		const url = new URL(page.url);
		if (q) url.searchParams.set('q', q);
		else url.searchParams.delete('q');
		if (url.search === page.url.search) return;
		goto(url, { replaceState: true, keepFocus: true, noScroll: true });
	}

	function run(q: string) {
		syncUrl(q);
		untrack(() => feed.load(q));
	}

	$effect(() => {
		const q = normaliseQuery(query);
		if (first) {
			first = false;
			untrack(() => feed.load(q));
			return;
		}
		const timer = setTimeout(() => run(q), q ? DEBOUNCE_MS : 0);
		return () => clearTimeout(timer);
	});

	const libraryEmpty = $derived(feed.ready && !feed.error && !feed.query && !query && feed.items.length === 0);
</script>

<div class="flex flex-col gap-6">
	<div class="flex flex-col gap-2">
		<h1 class="text-fc-2xl font-bold text-fc-fg">{t('artifacts.heading')}</h1>
		<p class="text-fc-sm text-fc-fg-muted">{t('artifacts.subtitle')}</p>
	</div>

	{#if !libraryEmpty}
		<div class="flex flex-col gap-2">
			<ArtifactSearchBar bind:value={query} busy={feed.loading && feed.ready} />
			<p class="min-h-4 text-fc-xs text-fc-fg-muted" aria-live="polite">
				{#if feed.ready && feed.query && !feed.error && feed.items.length > 0}
					{t('artifacts.results', { count: feed.total, query: feed.query })}
				{/if}
			</p>
		</div>
	{/if}

	{#if !feed.ready}
		<div class="flex items-center gap-3 text-fc-sm text-fc-fg-muted">
			<Spinner size="sm" label={t('common.loading')} />
			{t('artifacts.loading')}
		</div>
	{:else if feed.error}
		<Alert tone="danger" title={feed.query || query ? t('artifacts.searchFailed') : t('artifacts.loadFailed')}>
			{feed.error}
		</Alert>
	{:else if libraryEmpty}
		<EmptyState icon="solar:file-smile-linear" title={t('artifacts.emptyTitle')} description={t('artifacts.emptyText')} />
	{:else if feed.items.length === 0}
		<EmptyState icon={icons.search} title={t('artifacts.noMatch', { query: feed.query })} description={t('artifacts.noMatchText')}>
			<Button variant="secondary" onclick={() => (query = '')}>{t('artifacts.clearSearch')}</Button>
		</EmptyState>
	{:else}
		<div class="flex flex-col gap-6">
			<ul
				class="grid grid-cols-[repeat(auto-fill,minmax(min(18rem,100%),1fr))] gap-4 transition-opacity [&>*]:min-w-0 {feed.loading
					? 'opacity-60'
					: ''}"
				aria-busy={feed.loading}
			>
				{#each feed.items as art (`${art.space_id ?? ''}/${art.id}`)}
					<li class="flex [&>*]:w-full"><ArtifactCard artifact={art} query={feed.query} /></li>
				{/each}
			</ul>
			{#if feed.nextCursor}
				<div class="flex justify-center">
					<Button variant="secondary" disabled={feed.loadingMore} onclick={() => feed.more()}>
						{feed.loadingMore ? t('artifacts.loadingMore') : t('artifacts.loadMore')}
					</Button>
				</div>
			{/if}
		</div>
	{/if}
</div>
