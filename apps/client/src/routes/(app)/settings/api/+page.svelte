<script lang="ts">
	import {
		Alert,
		Badge,
		Button,
		Card,
		ConfirmModal,
		Drawer,
		Field,
		Input,
		SecretField,
		SettingsSection,
		Table,
		icons,
		toast
	} from '@facile/muse';
	import { backend, ApiError, type TokenInfo } from '$lib/backend';
	import { AGENT_PROMPT } from '$lib/agentPrompt';
	import { t } from '$lib/i18n/index.svelte';

	let tokens: TokenInfo[] = $state([]);
	let forbidden = $state(false);

	let createOpen = $state(false);
	let creating = $state(false);
	let newName = $state('');
	let createdToken = $state('');

	let revokeTarget = $state<TokenInfo | null>(null);
	let revokeOpen = $state(false);

	const endpoint = $derived(
		typeof window === 'undefined' ? '/api' : `${window.location.origin}/api`
	);

	$effect(() => {
		refresh();
	});

	async function refresh() {
		try {
			forbidden = false;
			tokens = await backend.tokensList();
		} catch (e) {
			if (e instanceof ApiError && e.status === 403) {
				forbidden = true;
				tokens = [];
				return;
			}
			toast.danger(e instanceof Error ? e.message : t('api.loadFailed'));
		}
	}

	function openCreate() {
		/* Reset first: reopening the drawer must never re-show a token from a previous run. */
		createdToken = '';
		newName = '';
		createOpen = true;
	}

	async function create(event: Event) {
		event.preventDefault();
		if (!newName.trim()) return;
		creating = true;
		try {
			const result = await backend.tokensCreate(newName.trim());
			createdToken = result.token ?? '';
			newName = '';
			await refresh();
		} catch (e) {
			toast.danger(e instanceof Error ? e.message : t('api.createFailed'), {
				title: t('api.notCreated')
			});
		} finally {
			creating = false;
		}
	}

	async function revoke() {
		const target = revokeTarget;
		if (!target) return;
		try {
			await backend.tokensDelete(target.name);
			toast.success(t('api.revoked', { name: target.name }));
			revokeTarget = null;
			await refresh();
		} catch (e) {
			toast.danger(e instanceof Error ? e.message : t('api.revokeFailed'));
			throw e;
		}
	}

	async function copyPrompt() {
		try {
			await navigator.clipboard.writeText(AGENT_PROMPT);
			toast.success(t('api.promptCopied'));
		} catch {
			toast.danger(t('api.clipboardDenied'));
		}
	}
</script>

<div class="flex flex-col gap-10">
	<SettingsSection
		title={t('api.endpoint')}
		description={t('api.endpointHelp')}
	>
		<SecretField value={endpoint} sensitive={false} label={t('api.baseUrl')} />
	</SettingsSection>

	<SettingsSection
		title={t('api.connect')}
		description={t('api.connectHelp')}
		bare
	>
		{#snippet actions()}
			<Button variant="outline" icon={icons.copy} onclick={copyPrompt}>{t('api.copyPrompt')}</Button>
		{/snippet}

		<Card class="overflow-hidden">
			<pre
				class="max-h-96 overflow-auto whitespace-pre-wrap font-fc-mono text-fc-xs leading-relaxed text-fc-fg">{AGENT_PROMPT}</pre>
		</Card>
		<p class="text-fc-xs text-fc-fg-muted">
			{t('api.promptNote')}
		</p>
	</SettingsSection>

	<SettingsSection
		title={t('api.tokens')}
		description={t('api.tokensHelp')}
		bare
	>
		{#snippet actions()}
			{#if !forbidden}
				<Button icon={icons.plus} onclick={openCreate}>{t('api.newToken')}</Button>
			{/if}
		{/snippet}

		{#if forbidden}
			<Alert tone="info">
				{t('api.adminOnly')}
			</Alert>
		{:else if tokens.length === 0}
			<Alert tone="info">
				{t('api.noTokens')}
			</Alert>
		{:else}
			<Table>
				<thead>
					<tr>
						<th scope="col">{t('common.name')}</th>
						<th scope="col">{t('api.scope')}</th>
						<th scope="col">{t('api.created')}</th>
						<th scope="col">{t('api.lastSeen')}</th>
						<th scope="col" class="text-right" aria-label={t('members.actions')}></th>
					</tr>
				</thead>
				<tbody>
					{#each tokens as token (token.name)}
						<tr>
							<td class="font-fc-mono font-medium text-fc-fg">{token.name}</td>
							<td>
								<Badge tone={token.scope === 'admin' ? 'accent' : 'neutral'}>
									{token.scope || t('api.machine')}
								</Badge>
							</td>
							<td class="whitespace-nowrap text-fc-fg-muted">
								{token.created_at?.slice(0, 10) || '—'}
							</td>
							<td class="whitespace-nowrap text-fc-fg-muted">
								{token.last_seen?.slice(0, 10) || t('api.never')}
							</td>
							<td class="text-right">
								<Button
									variant="ghost-danger"
									size="sm"
									icon={icons.revoke}
									aria-label={t('api.revokeLabel', { name: token.name })}
									onclick={() => {
										revokeTarget = token;
										revokeOpen = true;
									}}
								>
									{t('common.revoke')}
								</Button>
							</td>
						</tr>
					{/each}
				</tbody>
			</Table>
		{/if}
	</SettingsSection>
</div>

<Drawer bind:open={createOpen} title={t('api.drawerTitle')} showClose>
	{#if createdToken}
		<div class="flex flex-col gap-4">
			<Alert tone="warning" title={t('api.copyNow')}>
				{t('api.copyNowText')}
			</Alert>

			<!-- The one-time token starts revealed and stays that way: hiding a value the user has
			     not copied yet is theatre. -->
			<SecretField
				value={createdToken}
				visible
				autoHideMs={0}
				label={t('api.tokenLabel')}
				helper={t('api.tokenHelp')}
			/>

			<div class="flex justify-end">
				<Button onclick={() => (createOpen = false)}>{t('api.done')}</Button>
			</div>
		</div>
	{:else}
		<form class="flex flex-col gap-4" onsubmit={create}>
			<Field label={t('common.name')} helper={t('api.nameHelp')}>
				<Input bind:value={newName} placeholder="lucy" required disabled={creating} />
			</Field>
			<div class="flex justify-end gap-2 pt-1">
				<Button
					type="button"
					variant="ghost"
					disabled={creating}
					onclick={() => (createOpen = false)}
				>
					{t('common.cancel')}
				</Button>
				<Button type="submit" icon={icons.key} disabled={creating || newName.trim().length === 0}>
					{creating ? t('common.creating') : t('api.createToken')}
				</Button>
			</div>
		</form>
	{/if}
</Drawer>

<ConfirmModal
	bind:open={revokeOpen}
	tone="danger"
	title={t('api.revokeTitle', { name: revokeTarget?.name ?? '' })}
	description={t('api.revokeText')}
	confirmLabel={t('api.revokeConfirm')}
	cancelLabel={t('common.keep')}
	onConfirm={revoke}
	onCancel={() => (revokeTarget = null)}
/>
