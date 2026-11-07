<script lang="ts">
	import { Alert, EmptyState, icons } from '@facile/muse';
	import { backend } from '$lib/backend';
	import EntityCard from '$lib/components/EntityCard.svelte';
	import { t } from '$lib/i18n/index.svelte';

	let models: { type: string; path: string }[] = $state([]);
	let error = $state('');

	$effect(() => {
		backend
			.modelsList()
			.then((m) => (models = m))
			.catch((e) => (error = e instanceof Error ? e.message : t('models.loadFailed')));
	});
</script>

<div class="flex flex-col gap-10">
	<p class="text-fc-sm text-fc-fg-muted">
		{t('models.subtitle')}
	</p>

	{#if error}
		<Alert tone="danger" title={t('models.loadFailed')}>{error}</Alert>
	{:else if models.length === 0}
		<EmptyState
			icon={icons.code}
			title={t('models.emptyTitle')}
			description={t('models.emptyText')}
		/>
	{:else}
		<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 [&>*]:min-w-0">
			{#each models as model (model.path)}
				<EntityCard
					href="/models/{model.path}"
					icon={icons.code}
					title={model.type}
					meta="extensions/models/{model.path}"
				/>
			{/each}
		</div>
	{/if}
</div>
