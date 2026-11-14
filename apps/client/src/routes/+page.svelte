<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Button, DownloadSection, SignInWithFacile, icons } from '@facile/muse';
	import FlowTile from '$lib/components/landing/FlowTile.svelte';
	import HeroDemo from '$lib/components/landing/HeroDemo.svelte';
	import LandingOffers from '$lib/components/LandingOffers.svelte';
	import LandingShell from '$lib/components/landing/LandingShell.svelte';
	import { TOKEN_KEY } from '$lib/backend';
	import { CLI_COMMANDS, CLI_RELEASES } from '$lib/downloads';
	import { t } from '$lib/i18n/index.svelte';

	const SIGN_IN = '/api/auth/oidc';

	const FEATURES = [
		{ id: 'memory', icon: icons.folder },
		{ id: 'rules', icon: icons.shield },
		{ id: 'flows', icon: icons.bolt },
		{ id: 'sync', icon: icons.refresh }
	] as const;

	let visible = $state(false);

	onMount(() => {
		if (localStorage.getItem(TOKEN_KEY)) {
			goto('/memory');
			return;
		}
		visible = true;
	});
</script>

<svelte:head>
	<title>{t('landing.title')}</title>
	<meta name="description" content={t('landing.description')} />
</svelte:head>

{#if visible}
	<LandingShell>
		<section class="mx-auto flex w-full max-w-fc-xl flex-col gap-10 px-4 pt-12 sm:px-6 md:pt-20">
			<div class="flex max-w-fc-md flex-col gap-6">
				<h1 class="font-fc-title text-fc-3xl font-bold text-balance md:text-5xl">{t('landing.hero.title')}</h1>
				<p class="text-fc-md text-balance text-fc-fg-muted">{t('landing.hero.line')}</p>
				<div class="flex flex-wrap items-center gap-3">
					<SignInWithFacile href={SIGN_IN} size="lg" data-sveltekit-reload />
					<Button href="#download" variant="secondary" size="lg" icon={icons.download}>{t('landing.hero.cli')}</Button>
				</div>
			</div>
			<div class="w-full max-w-[960px]"><HeroDemo /></div>
		</section>

		<section id="features" aria-labelledby="features-title" class="mx-auto flex w-full max-w-fc-xl scroll-mt-24 flex-col gap-8 px-4 sm:px-6">
			<h2 id="features-title" class="font-fc-title text-fc-2xl font-bold">{t('landing.features.title')}</h2>
			<ul class="grid gap-4 sm:grid-cols-2">
				{#each FEATURES as feature (feature.id)}
					<li class="flex flex-col gap-3 rounded-fc-lg bg-fc-component p-6">
						<span class="flex size-10 items-center justify-center rounded-fc-md bg-fc-surface text-fc-fg" aria-hidden="true">
							<iconify-icon icon={feature.icon} width="20" height="20" class="block"></iconify-icon>
						</span>
						<h3 class="text-fc-md font-bold">{t(`landing.feature.${feature.id}.title`)}</h3>
						<p class="text-fc-sm text-fc-fg-muted">{t(`landing.feature.${feature.id}.body`)}</p>
					</li>
				{/each}
			</ul>
		</section>

		<section id="download" aria-labelledby="cli-title" class="mx-auto flex w-full max-w-fc-xl scroll-mt-24 flex-col gap-8 px-4 sm:px-6">
			<h2 id="cli-title" class="font-fc-title text-fc-2xl font-bold">{t('landing.cli.title')}</h2>
			<div class="grid items-start gap-4 lg:grid-cols-2">
				<FlowTile />
				<DownloadSection mode="device" allHref="/downloads" commands={CLI_COMMANDS} releasesHref={CLI_RELEASES} headingLevel={3} />
			</div>
		</section>

		<LandingOffers />

		<section class="mx-auto flex w-full max-w-fc-md flex-col items-center gap-6 px-4 text-center sm:px-6">
			<h2 class="font-fc-title text-fc-2xl font-bold text-balance">{t('landing.cta.title')}</h2>
			<div class="flex flex-wrap items-center justify-center gap-3">
				<SignInWithFacile href={SIGN_IN} size="lg" data-sveltekit-reload />
				<Button
					href="https://github.com/FacileStudio/Mycelium"
					target="_blank"
					rel="noopener noreferrer"
					variant="secondary"
					size="lg"
					iconRight={icons.code}
				>
					{t('landing.code')}
				</Button>
			</div>
		</section>
	</LandingShell>
{/if}
