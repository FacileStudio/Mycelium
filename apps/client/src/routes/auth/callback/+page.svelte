<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { TOKEN_KEY } from '$lib/backend';
	import { t } from '$lib/i18n/index.svelte';
	import { takeRedirect } from '$lib/redirect';

	onMount(() => {
		const params = new URLSearchParams(location.hash.replace(/^#/, ''));
		const token = params.get('token');
		if (!token) {
			goto('/login');
			return;
		}
		localStorage.setItem(TOKEN_KEY, token);
		history.replaceState(null, '', '/auth/callback');
		goto(takeRedirect());
	});
</script>

<svelte:head>
	<title>{t('login.title')}</title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center bg-fc-page">
	<p class="text-fc-sm text-fc-fg-muted">{t('callback.signingIn')}</p>
</div>
