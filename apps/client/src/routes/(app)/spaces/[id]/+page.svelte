<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import {
		Badge,
		Button,
		Card,
		ConfirmModal,
		Field,
		Input,
		Modal,
		Spinner,
		Textarea,
		icons,
		toast
	} from '@facile/muse';
	import { backend, type Space, type SpaceMember, type UserInfo } from '$lib/backend';
	import { setSpaces } from '$lib/space.svelte';
	import SpaceMembers from '$lib/components/SpaceMembers.svelte';
	import { t } from '$lib/i18n/index.svelte';

	const roleTone = { owner: 'owner', admin: 'admin', member: 'neutral' } as const;

	const id = $derived(page.params.id ?? '');
	let space = $state<Space | null>(null);
	let members: SpaceMember[] = $state([]);
	let users: UserInfo[] = $state([]);
	let loading = $state(true);
	let registreURL = $state('');

	let editOpen = $state(false);
	let draftName = $state('');
	let draftDescription = $state('');
	let savingEdit = $state(false);
	let editError = $state('');

	let leaveOpen = $state(false);
	let deleteOpen = $state(false);

	const managed = $derived(registreURL !== '');
	const linkedID = $derived(space?.registre_space_id ?? '');
	const manageHref = $derived(linkedID ? `${registreURL}/${encodeURIComponent(linkedID)}` : registreURL);
	const isOwner = $derived(space?.role === 'owner');
	const canManage = $derived(space?.role === 'owner' || space?.role === 'admin');
	const candidates = $derived(users.filter((u) => !members.some((m) => m.email === u.email)));

	$effect(() => {
		const spaceId = id;
		loading = true;
		Promise.all([
			backend.spacesDirectory().then((d) => {
				setSpaces(d.spaces);
				registreURL = d.registre_url;
				space = d.spaces.find((s) => s.id === spaceId) ?? null;
			}),
			backend.spaceMembers(spaceId).then((m) => (members = m)),
			backend
				.usersList()
				.then((u) => (users = u ?? []))
				.catch(() => (users = []))
		])
			.catch((e) => toast.danger(e instanceof Error ? e.message : t('spaces.loadOne')))
			.finally(() => (loading = false));
	});

	async function refreshMembers() {
		members = await backend.spaceMembers(id);
	}

	function openEdit() {
		if (!space) return;
		draftName = space.name;
		draftDescription = space.description;
		editError = '';
		editOpen = true;
	}

	async function saveEdit(event: Event) {
		event.preventDefault();
		if (!draftName.trim()) return;
		savingEdit = true;
		editError = '';
		try {
			space = await backend.spaceUpdate(id, draftName.trim(), draftDescription.trim());
			editOpen = false;
			toast.success(t('spaces.updated'));
		} catch (e) {
			editError = e instanceof Error ? e.message : t('spaces.updateFailed');
		} finally {
			savingEdit = false;
		}
	}

	async function leaveSpace() {
		try {
			await backend.spaceLeave(id);
			toast.success(t('spaces.left'));
			goto('/spaces');
		} catch (e) {
			toast.danger(e instanceof Error ? e.message : t('spaces.leaveFailed'));
			throw e;
		}
	}

	async function deleteSpace() {
		try {
			await backend.spaceDelete(id);
			toast.success(t('spaces.deleted'));
			goto('/spaces');
		} catch (e) {
			toast.danger(e instanceof Error ? e.message : t('spaces.deleteFailed'));
			throw e;
		}
	}
</script>

<div class="flex flex-col gap-10">
	<a
		href="/spaces"
		class="inline-flex w-fit items-center gap-1 text-fc-sm text-fc-fg-muted transition-colors hover:text-fc-fg focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-fc-ring"
	>
		<iconify-icon icon={icons.chevronLeft} width="16" height="16" class="block"></iconify-icon>
		{t('spaces.title')}
	</a>

	{#if loading}
		<div class="flex items-center gap-3 text-fc-sm text-fc-fg-muted">
			<Spinner size="sm" label={t('common.loading')} /> {t('common.loadingDots')}
		</div>
	{:else if !space}
		<Card class="text-fc-sm text-fc-fg-muted">{t('spaces.missing')}</Card>
	{:else}
		<div class="flex flex-wrap items-start justify-between gap-4">
			<div class="flex min-w-0 items-center gap-3">
				<span
					class="flex size-10 shrink-0 items-center justify-center rounded-fc-md bg-fc-surface text-fc-fg"
				>
					<iconify-icon icon={icons.space} width="20" height="20" class="block"
					></iconify-icon>
				</span>
				<div class="min-w-0">
					<div class="flex items-center gap-2">
						<h1 class="truncate text-fc-xl font-semibold text-fc-fg">{space.name}</h1>
						<Badge tone={roleTone[space.role]}>{t(`role.${space.role}`)}</Badge>
					</div>
					<p class="truncate text-fc-sm text-fc-fg-muted">
						{space.description || t('spaces.noDescription')}
					</p>
				</div>
			</div>

			<div class="flex shrink-0 flex-wrap items-center gap-2">
				{#if managed}
					<Button variant="outline" icon={icons.arrow} href={manageHref} target="_blank" rel="noopener">
						{t('spaces.manage')}
					</Button>
				{/if}
				{#if canManage && !linkedID}
					<Button variant="outline" icon={icons.edit} onclick={openEdit}>{t('common.edit')}</Button>
				{/if}
				{#if !managed}
					<Button variant="ghost" icon={icons.logout} onclick={() => (leaveOpen = true)}>
						{t('spaces.leave')}
					</Button>
				{/if}
				{#if isOwner && !linkedID}
					<Button variant="ghost-danger" icon={icons.remove} onclick={() => (deleteOpen = true)}>
						{t('common.delete')}
					</Button>
				{/if}
			</div>
		</div>

		<SpaceMembers
			{id}
			spaceName={space.name}
			{members}
			{candidates}
			isOwner={isOwner && !managed}
			canManage={canManage && !managed}
			onChanged={refreshMembers}
		/>
	{/if}
</div>

<Modal bind:open={editOpen} title={t('spaces.edit')} showClose>
	<form class="flex flex-col gap-4" onsubmit={saveEdit}>
		<Field label={t('common.name')} error={editError || undefined}>
			<Input bind:value={draftName} disabled={savingEdit} required />
		</Field>
		<Field label={t('spaces.description')}>
			<Textarea bind:value={draftDescription} rows={3} disabled={savingEdit} />
		</Field>
		<div class="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
			<Button
				type="button"
				variant="ghost"
				class="w-full sm:w-auto"
				onclick={() => (editOpen = false)}
			>
				{t('common.cancel')}
			</Button>
			<Button
				type="submit"
				icon={icons.check}
				class="w-full sm:w-auto"
				disabled={savingEdit || draftName.trim().length === 0}
			>
				{savingEdit ? t('common.saving') : t('spaces.saveChanges')}
			</Button>
		</div>
	</form>
</Modal>

<ConfirmModal
	bind:open={leaveOpen}
	tone="danger"
	title={t('spaces.leaveTitle', { name: space?.name ?? t('spaces.thisSpace') })}
	description={t('spaces.leaveText')}
	confirmLabel={t('spaces.leaveConfirm')}
	cancelLabel={t('spaces.stay')}
	onConfirm={leaveSpace}
/>

<ConfirmModal
	bind:open={deleteOpen}
	tone="danger"
	title={t('spaces.deleteTitle', { name: space?.name ?? t('spaces.thisSpace') })}
	description={t('spaces.deleteText')}
	confirmLabel={t('spaces.deleteConfirm')}
	cancelLabel={t('common.keep')}
	onConfirm={deleteSpace}
/>
