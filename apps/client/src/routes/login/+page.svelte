<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { Alert, AuthScreen } from '@facile/muse';
	import { backend, TOKEN_KEY } from '$lib/backend';
	import { rememberLocale, t } from '$lib/i18n/index.svelte';
	import { rememberRedirect, safeRedirect } from '$lib/redirect';

	let failed = $state(false);

	onMount(async () => {
		const target = safeRedirect(page.url.searchParams.get('redirect'));
		if (page.url.searchParams.get('error') === 'sso') {
			localStorage.removeItem(TOKEN_KEY);
			failed = true;
		}
		if (target !== '/memory') rememberRedirect(target);
		if (!localStorage.getItem(TOKEN_KEY)) return;
		try {
			await backend.authMe();
			goto(target);
		} catch {
			localStorage.removeItem(TOKEN_KEY);
		}
	});
</script>

<svelte:head>
	<title>{t('login.title')}</title>
</svelte:head>

<AuthScreen
	appName="Mycelium"
	href="/api/auth/oidc"
	data-sveltekit-reload
	onLocaleChange={rememberLocale}
>
	{#snippet logo()}
		<iconify-icon icon="solar:structure-bold-duotone" width="40" height="40" class="block text-fc-fg"
		></iconify-icon>
	{/snippet}
	{#snippet tagline()}
		{t('login.tagline')}
	{/snippet}
	{#snippet notice()}
		{#if failed}
			<Alert tone="danger">{t('login.failed')}</Alert>
		{/if}
	{/snippet}
</AuthScreen>
