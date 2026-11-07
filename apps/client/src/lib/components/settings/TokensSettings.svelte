<script lang="ts">
	import { Alert, SettingsSection, TokensSection, toast, type TokenDraft, type TokenItem } from '@facile/muse';
	import { backend, ApiError, type TokenInfo } from '$lib/backend';
	import { t } from '$lib/i18n/index.svelte';

	let tokens: TokenInfo[] = $state([]);
	let forbidden = $state(false);

	const items = $derived<TokenItem[]>(
		tokens.map((token) => ({
			id: token.name,
			name: token.name,
			scope: token.scope || t('api.machine'),
			lastUsed: token.last_seen?.slice(0, 10) || undefined
		}))
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

	async function create(draft: TokenDraft): Promise<string> {
		const result = await backend.tokensCreate(draft.name);
		await refresh();
		return result.token ?? '';
	}

	async function revoke(token: TokenItem) {
		try {
			await backend.tokensDelete(token.name);
			toast.success(t('api.revoked', { name: token.name }));
			await refresh();
		} catch (e) {
			toast.danger(e instanceof Error ? e.message : t('api.revokeFailed'));
			throw e;
		}
	}
</script>

{#if forbidden}
	<SettingsSection title={t('api.tokens')} description={t('api.tokensHelp')}>
		<Alert tone="info">{t('api.adminOnly')}</Alert>
	</SettingsSection>
{:else}
	<TokensSection tokens={items} onCreate={create} onRevoke={revoke} expiryDays={[]} />
{/if}
