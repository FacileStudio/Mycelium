<script lang="ts">
	import {
		Avatar,
		Badge,
		Button,
		Card,
		ConfirmModal,
		Field,
		Select,
		Table,
		icons,
		toast
	} from '@facile/muse';
	import { backend, type SpaceMember, type SpaceRole, type UserInfo } from '$lib/backend';
	import { t } from '$lib/i18n/index.svelte';

	let {
		id,
		spaceName,
		members,
		candidates,
		isOwner,
		canManage,
		onChanged
	}: {
		id: string;
		spaceName: string | undefined;
		members: SpaceMember[];
		candidates: UserInfo[];
		isOwner: boolean;
		canManage: boolean;
		onChanged: () => Promise<void> | void;
	} = $props();

	const roleTone = { owner: 'owner', admin: 'admin', member: 'neutral' } as const;

	let newMemberEmail = $state('');
	let newMemberRole = $state<SpaceRole>('member');
	let adding = $state(false);

	let removeTarget = $state<SpaceMember | null>(null);
	let removeOpen = $state(false);

	async function addMember(event: Event) {
		event.preventDefault();
		if (!newMemberEmail) return;
		adding = true;
		try {
			await backend.spaceMemberAdd(id, newMemberEmail, newMemberRole);
			toast.success(t('members.added', { email: newMemberEmail, space: spaceName ?? t('members.theSpace') }));
			newMemberEmail = '';
			newMemberRole = 'member';
			await onChanged();
		} catch (e) {
			toast.danger(e instanceof Error ? e.message : t('members.addFailed'));
		} finally {
			adding = false;
		}
	}

	async function changeRole(email: string, role: SpaceRole) {
		try {
			await backend.spaceMemberUpdate(id, email, role);
			toast.success(t('members.roleChanged', { email, role: t(`role.${role}`) }));
		} catch (e) {
			toast.danger(e instanceof Error ? e.message : t('members.roleFailed'));
		} finally {
			await onChanged();
		}
	}

	async function removeMember() {
		const target = removeTarget;
		if (!target) return;
		try {
			await backend.spaceMemberRemove(id, target.email);
			toast.success(t('members.removed', { email: target.email }));
			removeTarget = null;
			await onChanged();
		} catch (e) {
			toast.danger(e instanceof Error ? e.message : t('members.removeFailed'));
			throw e;
		}
	}
</script>

<section class="flex flex-col gap-4">
	<div class="flex flex-col gap-1">
		<h2 class="text-fc-lg font-semibold text-fc-fg">{t('members.title')}</h2>
		<p class="text-fc-sm text-fc-fg-muted">
			{t('members.subtitle')}
		</p>
	</div>

	<Table>
		<thead>
			<tr>
				<th scope="col">{t('members.email')}</th>
				<th scope="col">{t('common.name')}</th>
				<th scope="col">{t('members.role')}</th>
				{#if canManage}
					<th scope="col" class="text-right" aria-label={t('members.actions')}></th>
				{/if}
			</tr>
		</thead>
		<tbody>
			{#each members as member (member.email)}
				<tr>
					<td class="font-medium text-fc-fg">
						<span class="flex items-center gap-3">
							<Avatar
								src={member.avatar_url || undefined}
								name={member.name || member.email}
								alt=""
								size="sm"
							/>
							{member.email}
						</span>
					</td>
					<td class="text-fc-fg-muted">{member.name}</td>
					<td>
						{#if isOwner}
							<Select
								value={member.role}
								aria-label={t('members.roleFor', { email: member.email })}
								class="h-9 min-w-28"
								onchange={(e) =>
									changeRole(member.email, e.currentTarget.value as SpaceRole)}
							>
								<option value="member">{t('role.member')}</option>
								<option value="admin">{t('role.admin')}</option>
								<option value="owner">{t('role.owner')}</option>
							</Select>
						{:else}
							<Badge tone={roleTone[member.role]}>{t(`role.${member.role}`)}</Badge>
						{/if}
					</td>
					{#if canManage}
						<td class="text-right">
							<Button
								variant="ghost-danger"
								size="sm"
								icon={icons.remove}
								aria-label={t('members.removeLabel', { email: member.email })}
								onclick={() => {
									removeTarget = member;
									removeOpen = true;
								}}
							>
								{t('members.remove')}
							</Button>
						</td>
					{/if}
				</tr>
			{/each}
		</tbody>
	</Table>

	{#if canManage}
		<Card class="flex flex-col gap-4">
			<p class="text-fc-sm font-medium text-fc-fg">{t('members.add')}</p>
			<form class="flex flex-col gap-3 sm:flex-row sm:items-end" onsubmit={addMember}>
				<div class="min-w-0 flex-1">
					<Field label={t('members.user')} helper={t('members.userHelp')}>
						<Select bind:value={newMemberEmail} disabled={adding || candidates.length === 0}>
							<option value="" disabled>{t('members.select')}</option>
							{#each candidates as user (user.email)}
								<option value={user.email}>
									{user.name ? `${user.name} — ${user.email}` : user.email}
								</option>
							{/each}
						</Select>
					</Field>
				</div>
				<Field label={t('members.role')}>
					<Select bind:value={newMemberRole} class="min-w-32" disabled={adding}>
						<option value="member">{t('role.member')}</option>
						<option value="admin">{t('role.admin')}</option>
						{#if isOwner}
							<option value="owner">{t('role.owner')}</option>
						{/if}
					</Select>
				</Field>
				<Button type="submit" icon={icons.plus} disabled={adding || !newMemberEmail}>
					{adding ? t('members.adding') : t('members.addButton')}
				</Button>
			</form>
		</Card>
	{/if}
</section>

<ConfirmModal
	bind:open={removeOpen}
	tone="danger"
	title={t('members.removeTitle', { email: removeTarget?.email ?? t('members.thisMember') })}
	description={t('members.removeText')}
	confirmLabel={t('members.removeConfirm')}
	cancelLabel={t('members.keep')}
	onConfirm={removeMember}
	onCancel={() => (removeTarget = null)}
/>
