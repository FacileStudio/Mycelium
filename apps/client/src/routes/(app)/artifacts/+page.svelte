<script lang="ts">
	import { Alert, EmptyState, Input, Spinner, icons } from '@facile/muse';
	import EntityCard from '$lib/components/EntityCard.svelte';
	import { backend, type ArtifactSummary } from '$lib/backend';
	import { formatDate } from '$lib/i18n/format';
	import { t } from '$lib/i18n/index.svelte';

	let artifacts: ArtifactSummary[] = $state([]);
	let loading = $state(true);
	let error = $state('');
	let query = $state('');

	$effect(() => {
		loading = true;
		backend
			.artifactsList()
			.then((r) => (artifacts = r ?? []))
			.catch((e) => (error = e instanceof Error ? e.message : t('artifacts.loadFailed')))
			.finally(() => (loading = false));
	});

	const filtered = $derived(
		(artifacts ?? []).filter((r) => {
			const q = query.trim().toLowerCase();
			if (!q) return true;
			return r.title.toLowerCase().includes(q) || r.machine.toLowerCase().includes(q) || r.id.toLowerCase().includes(q);
		})
	);

	function shortDate(iso: string): string {
		try {
			return formatDate(iso, { year: 'numeric', month: 'short', day: 'numeric' });
		} catch {
			return iso;
		}
	}
</script>

<div class="flex flex-col gap-8">
	<div class="flex flex-col gap-2">
		<h1 class="text-fc-2xl font-bold text-fc-fg">{t('artifacts.heading')}</h1>
		<p class="text-fc-sm text-fc-fg-muted">
			{t('artifacts.subtitle')}
		</p>
	</div>

	{#if loading}
		<div class="flex items-center gap-3 text-fc-sm text-fc-fg-muted">
			<Spinner size="sm" label={t('common.loading')} /> {t('artifacts.loading')}
		</div>
	{:else if error}
		<Alert tone="danger" title={t('artifacts.loadFailed')}>{error}</Alert>
	{:else if artifacts.length === 0}
		<EmptyState
			icon="solar:file-smile-linear"
			title={t('artifacts.emptyTitle')}
			description={t('artifacts.emptyText')}
		/>
	{:else}
		<div class="flex flex-col gap-4">
			<div class="max-w-md">
				<Input
					type="search"
					placeholder={t('artifacts.filter')}
					bind:value={query}
				/>
			</div>

			{#if filtered.length === 0}
				<p class="py-8 text-center text-fc-sm text-fc-fg-muted">
					{t('artifacts.noMatch', { query })}
				</p>
			{:else}
				<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 [&>*]:min-w-0">
					{#each filtered as art (art.id)}
						<EntityCard
							href="/artifacts/{art.id}"
							icon="solar:file-smile-linear"
							title={art.title}
							meta="artifacts/{art.id}.md"
						>
							{#snippet footer()}
								<span class="inline-flex min-w-0 items-center gap-1 truncate font-fc-mono">
									<iconify-icon
										icon={icons.server}
										width="14"
										height="14"
										class="block shrink-0"
									></iconify-icon>
									{art.machine || t('artifacts.unknownMachine')}
								</span>
								<span class="shrink-0">{t('artifacts.updated', { date: shortDate(art.updated || art.created) })}</span>
							{/snippet}
						</EntityCard>
					{/each}
				</div>
			{/if}
		</div>
	{/if}
</div>
