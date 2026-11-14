<script lang="ts">
	import type { Snippet } from 'svelte';
	import { LandingFooter, LandingHeader, SignInWithFacile } from '@facile/muse';
	import { rememberLocale, t } from '$lib/i18n/index.svelte';

	let { children }: { children: Snippet } = $props();

	const ICON = 'solar:structure-bold-duotone';

	const links = $derived([
		{ href: '/#features', label: t('landing.nav.features') },
		{ href: '/downloads', label: t('landing.nav.downloads') }
	]);

	const columns = $derived([
		{ title: 'Mycelium', links },
		{
			title: 'Facile',
			links: [
				{ href: 'https://github.com/FacileStudio/Mycelium', label: 'GitHub', external: true },
				{ href: 'https://facile.studio', label: 'facile.studio', external: true }
			]
		}
	]);
</script>

<div class="min-h-dvh overflow-x-clip bg-fc-page text-fc-fg">
	<LandingHeader appName="Mycelium" icon={ICON} {links} skipHref="#main">
		{#snippet actions()}
			<SignInWithFacile href="/api/auth/oidc" size="sm" data-sveltekit-reload />
		{/snippet}
	</LandingHeader>

	<main id="main" class="flex flex-col gap-24 pb-24">
		{@render children()}
	</main>

	<LandingFooter
		appName="Mycelium"
		icon={ICON}
		tagline={t('landing.footerTagline')}
		{columns}
		legal={{ onLocaleChange: rememberLocale }}
	/>
</div>
