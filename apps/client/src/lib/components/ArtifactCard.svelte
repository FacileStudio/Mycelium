<script lang="ts">
	import { icons } from '@facile/muse';
	import EntityCard from './EntityCard.svelte';
	import type { ArtifactHit } from '$lib/backend';
	import { highlightParts } from '$lib/artifactSearch';
	import { formatDate } from '$lib/i18n/format';
	import { t } from '$lib/i18n/index.svelte';

	let { artifact, query = '' }: { artifact: ArtifactHit; query?: string } = $props();

	const href = $derived(
		`/artifacts/${encodeURIComponent(artifact.id)}?space=${encodeURIComponent(artifact.space_id ?? '')}`
	);
	const parts = $derived(query && artifact.snippet ? highlightParts(artifact.snippet, query) : []);
	const date = $derived(artifact.updated || artifact.created);

	function shortDate(iso: string): string {
		try {
			return formatDate(iso, { year: 'numeric', month: 'short', day: 'numeric' });
		} catch {
			return iso;
		}
	}
</script>

<EntityCard {href} icon="solar:file-smile-linear" title={artifact.title}>
	{#if parts.length}
		<p class="line-clamp-3 text-fc-xs leading-relaxed text-fc-fg-muted">
			{#each parts as part, i (i)}
				{#if part.hit}<mark class="bg-transparent font-medium text-fc-fg">{part.text}</mark>{:else}{part.text}{/if}
			{/each}
		</p>
	{/if}
	{#snippet footer()}
		<span class="inline-flex min-w-0 items-center gap-1">
			<iconify-icon icon={icons.folder} width="14" height="14" class="block shrink-0"></iconify-icon>
			<span class="truncate">{artifact.space_name || t('nav.personal')}</span>
		</span>
		<time class="shrink-0" datetime={date}>{shortDate(date)}</time>
	{/snippet}
</EntityCard>
