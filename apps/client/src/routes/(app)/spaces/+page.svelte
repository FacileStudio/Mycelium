<script lang="ts">
	import { goto } from '$app/navigation';
	import { Badge, Button, EmptyState, Field, Input, Modal, Textarea, icons, toast } from '@facile/muse';
	import { backend, type Space, type SpaceRole } from '$lib/backend';
	import { setSpaces } from '$lib/space.svelte';
	import EntityCard from '$lib/components/EntityCard.svelte';
	import { t } from '$lib/i18n/index.svelte';

	const roleTone = { owner: 'owner', admin: 'admin', member: 'neutral' } as const;

	let spaces: Space[] = $state([]);
	let registreURL = $state('');
	let createOpen = $state(false);
	let creating = $state(false);
	let draftName = $state('');
	let draftDescription = $state('');
	let createError = $state('');

	$effect(() => {
		backend
			.spacesDirectory()
			.then((d) => {
				spaces = d.spaces;
				registreURL = d.registre_url;
				setSpaces(d.spaces);
			})
			.catch((e) => toast.danger(e instanceof Error ? e.message : t('spaces.loadFailed')));
	});

	function openCreate() {
		draftName = '';
		draftDescription = '';
		createError = '';
		createOpen = true;
	}

	async function createSpace(event: Event) {
		event.preventDefault();
		const name = draftName.trim();
		if (!name) return;
		creating = true;
		createError = '';
		try {
			const space = await backend.spaceCreate(name, draftDescription.trim());
			createOpen = false;
			toast.success(t('spaces.created', { name: space.name }));
			goto(`/spaces/${space.id}`);
		} catch (e) {
			createError = e instanceof Error ? e.message : t('spaces.createFailed');
		} finally {
			creating = false;
		}
	}

	function tone(role: SpaceRole) {
		return roleTone[role] ?? 'neutral';
	}
</script>

<div class="flex flex-col gap-10">
	<div class="flex flex-wrap items-start justify-between gap-4">
		<div class="flex min-w-0 flex-col gap-2">
			<h1 class="text-fc-2xl font-semibold text-fc-fg">{t('spaces.title')}</h1>
			<p class="text-fc-sm text-fc-fg-muted">
				{registreURL
					? t('spaces.subtitleRegistre')
					: t('spaces.subtitleLocal')}
			</p>
		</div>
		{#if registreURL}
			<Button icon={icons.arrow} href={registreURL} target="_blank" rel="noopener">
				{t('spaces.manage')}
			</Button>
		{:else}
			<Button icon={icons.plus} onclick={openCreate}>{t('spaces.new')}</Button>
		{/if}
	</div>

	{#if spaces.length === 0}
		<EmptyState
			icon={icons.space}
			title={t('spaces.emptyTitle')}
			description={registreURL
				? t('spaces.emptyRegistre')
				: t('spaces.emptyLocal')}
		>
			{#if registreURL}
				<Button variant="outline" icon={icons.arrow} href={registreURL} target="_blank" rel="noopener">
					{t('spaces.manage')}
				</Button>
			{:else}
				<Button variant="outline" icon={icons.plus} onclick={openCreate}>{t('spaces.new')}</Button>
			{/if}
		</EmptyState>
	{:else}
		<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 [&>*]:min-w-0">
			{#each spaces as space (space.id)}
				<EntityCard
					href="/spaces/{space.id}"
					icon={icons.space}
					title={space.name}
					meta={space.description || t('spaces.noDescription')}
					mono={false}
				>
					{#snippet trailing()}
						<Badge tone={tone(space.role)}>{t(`role.${space.role}`)}</Badge>
					{/snippet}
				</EntityCard>
			{/each}
		</div>
	{/if}
</div>

<Modal bind:open={createOpen} title={t('spaces.new')} showClose>
	<form class="flex flex-col gap-4" onsubmit={createSpace}>
		<Field label={t('common.name')} error={createError || undefined}>
			<Input bind:value={draftName} placeholder="Facile" disabled={creating} required />
		</Field>
		<Field label={t('spaces.description')} helper={t('spaces.descriptionHelp')}>
			<Textarea
				bind:value={draftDescription}
				rows={3}
				placeholder={t('spaces.descriptionPlaceholder')}
				disabled={creating}
			/>
		</Field>
		<div class="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
			<Button
				type="button"
				variant="ghost"
				class="w-full sm:w-auto"
				onclick={() => (createOpen = false)}
			>
				{t('common.cancel')}
			</Button>
			<Button
				type="submit"
				icon={icons.plus}
				class="w-full sm:w-auto"
				disabled={creating || draftName.trim().length === 0}
			>
				{creating ? t('common.creating') : t('spaces.create')}
			</Button>
		</div>
	</form>
</Modal>
