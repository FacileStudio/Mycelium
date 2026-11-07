<script lang="ts">
	import { goto } from '$app/navigation';
	import { Alert, Button, EmptyState, Field, Input, Modal, icons, toast } from '@facile/muse';
	import { backend } from '$lib/backend';
	import EntityCard from '$lib/components/EntityCard.svelte';
	import { t } from '$lib/i18n/index.svelte';

	let rules: string[] = $state([]);
	let createOpen = $state(false);
	let draftName = $state('');
	let creating = $state(false);
	/* Separate from loadError: this one belongs to the create form, and showing a failed
	   list fetch inside a modal nobody opened is how a 403 came to look like an empty tree. */
	let error = $state('');
	let loadError = $state('');

	$effect(() => {
		backend
			.rulesList()
			.then((r) => (rules = r ?? []))
			.catch((e) => (loadError = e instanceof Error ? e.message : t('rules.loadFailed')));
	});

	function openCreate() {
		draftName = '';
		error = '';
		createOpen = true;
	}

	async function createRule(event: Event) {
		event.preventDefault();
		const name = draftName.trim();
		if (!name) return;
		creating = true;
		error = '';
		try {
			await backend.ruleSave(name, `# ${name}\n`);
			createOpen = false;
			toast.success(t('rules.created', { name }));
			goto(`/rules/${name}`);
		} catch (e) {
			error = e instanceof Error ? e.message : t('rules.createFailed');
		} finally {
			creating = false;
		}
	}
</script>

<div class="flex flex-col gap-10">
	<div class="flex flex-wrap items-start justify-between gap-4">
		<p class="min-w-0 text-fc-sm text-fc-fg-muted">
			{t('rules.subtitle')}
		</p>
		<Button icon={icons.plus} onclick={openCreate}>{t('rules.new')}</Button>
	</div>

	{#if loadError}
		<Alert tone="danger" title={t('rules.loadFailed')}>{loadError}</Alert>
	{:else if rules.length === 0}
		<EmptyState
			icon={icons.shield}
			title={t('rules.emptyTitle')}
			description={t('rules.emptyText')}
		>
			<Button variant="secondary" icon={icons.plus} onclick={openCreate}>{t('rules.new')}</Button>
		</EmptyState>
	{:else}
		<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 [&>*]:min-w-0">
			{#each rules as rule (rule)}
				<EntityCard
					href="/rules/{rule}"
					icon={icons.shield}
					title={rule}
					meta="rules/{rule}.md"
				/>
			{/each}
		</div>
	{/if}
</div>

<Modal bind:open={createOpen} title={t('rules.new')} showClose>
	<form class="flex flex-col gap-4" onsubmit={createRule}>
		<Field
			label={t('common.name')}
			helper={t('rules.nameHelp')}
			error={error || undefined}
		>
			<Input bind:value={draftName} placeholder="30-testing" disabled={creating} required />
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
				{creating ? t('common.creating') : t('rules.create')}
			</Button>
		</div>
	</form>
</Modal>
