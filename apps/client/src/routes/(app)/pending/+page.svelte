<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { Badge, Button, Card, EmptyState, Skeleton, icons, toast } from '@facile/muse';
	import { ApiError, backend, type PendingAction } from '$lib/backend';
	import { formatDate } from '$lib/i18n/format';
	import { t } from '$lib/i18n/index.svelte';

	let actions = $state<PendingAction[]>([]);
	let loading = $state(true);
	let busy = $state('');

	const focused = $derived(page.url.searchParams.get('action') ?? '');
	const ordered = $derived([
		...actions.filter((action) => action.id === focused),
		...actions.filter((action) => action.id !== focused)
	]);

	onMount(load);

	/** Loads the actions waiting for the signed-in person. */
	async function load() {
		loading = true;
		try {
			actions = await backend.pendingList();
		} catch (err) {
			toast.danger(err instanceof ApiError ? err.message : t('pending.loadFailed'));
		}
		loading = false;
	}

	/** Runs one action once, then drops it from the list. */
	async function confirm(action: PendingAction) {
		busy = action.id;
		try {
			await backend.pendingConfirm(action.id);
			toast.success(t(isArtifact(action) ? 'pending.artifactDeleted' : 'pending.deleted'));
			actions = actions.filter((item) => item.id !== action.id);
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t('pending.failed'));
			await load();
		}
		busy = '';
	}

	/** Closes one action without running it. */
	async function dismiss(action: PendingAction) {
		busy = action.id;
		try {
			await backend.pendingDismiss(action.id);
			toast.success(t('pending.dismissed'));
			actions = actions.filter((item) => item.id !== action.id);
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t('pending.failed'));
		}
		busy = '';
	}

	/** Tells an artifact delete from a page delete. */
	function isArtifact(action: PendingAction): boolean {
		return action.kind === 'delete_artifact';
	}

	/** Names the client that asked: the host of a client URL, else its id. */
	function clientName(clientId: string): string {
		if (clientId === 'bulle') return 'Bulle';
		try {
			return new URL(clientId).host;
		} catch {
			return clientId || '?';
		}
	}
</script>

<svelte:head>
	<title>{t('pending.title')}</title>
</svelte:head>

<div class="w-full px-4 py-6 md:px-8">
	<div class="mx-auto flex max-w-2xl flex-col gap-8">
		<div class="flex flex-col gap-1">
			<h1 class="text-fc-2xl font-semibold text-fc-fg">{t('pending.title')}</h1>
			<p class="text-fc-sm text-fc-fg-muted">{t('pending.intro')}</p>
		</div>

		{#if loading}
			<Skeleton class="h-48 w-full" />
		{:else if actions.length === 0}
			<EmptyState icon={icons.check} title={t('pending.none')} description={t('pending.noneBody')} />
		{:else}
			<section class="flex flex-col gap-4">
				{#each ordered as action (action.id)}
					<Card class="flex flex-col gap-4">
						<div class="flex flex-wrap items-center justify-between gap-2">
							<Badge tone="danger">{t(isArtifact(action) ? 'pending.deleteArtifact' : 'pending.deletePage')}</Badge>
							<span class="text-fc-xs text-fc-fg-muted">
								{t('pending.askedBy', { client: clientName(action.client_id) })} · {formatDate(action.created_at, {
									day: 'numeric',
									month: 'short',
									hour: '2-digit',
									minute: '2-digit'
								})}
							</span>
						</div>

						{#if action.path}
							<div class="flex flex-col gap-1 rounded-fc-md bg-fc-surface p-4">
								{#if action.title}
									<p class="break-words text-fc-md font-medium text-fc-fg">{action.title}</p>
								{/if}
								<p class="break-all font-mono text-fc-sm text-fc-fg">{action.path}</p>
								<p class="text-fc-sm text-fc-fg-muted">{action.space_name}</p>
								{#if action.log}
									<p class="break-words text-fc-sm text-fc-fg-muted">{t('pending.reason', { log: action.log })}</p>
								{/if}
							</div>
						{:else}
							<p class="break-words text-fc-sm text-fc-fg-muted">{action.summary}</p>
						{/if}

						{#if action.state === 'changed'}
							<p class="text-fc-sm text-fc-warning">{t('pending.changed')}</p>
						{:else if action.state === 'missing'}
							<p class="text-fc-sm text-fc-warning">{t('pending.missing')}</p>
						{/if}

						<div class="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
							<Button variant="ghost" class="w-full sm:w-auto" disabled={busy === action.id} onclick={() => dismiss(action)}>
								{t('pending.dismiss')}
							</Button>
							{#if action.state === 'ready'}
								<Button
									variant="danger"
									class="w-full sm:w-auto"
									disabled={busy === action.id}
									icon={icons.remove}
									onclick={() => confirm(action)}
								>
									{t(isArtifact(action) ? 'pending.confirmArtifact' : 'pending.confirmDelete')}
								</Button>
							{/if}
						</div>
					</Card>
				{/each}
			</section>
		{/if}
	</div>
</div>
