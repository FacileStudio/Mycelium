<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import LandingOffers from '$lib/components/LandingOffers.svelte';
	import {
		Button,
		Card,
		LandingFooter,
		LandingHeader,
		SignInWithFacile,
		icons
	} from '@facile/muse';
	import WordReveal from '$lib/components/motion/WordReveal.svelte';
	import { TOKEN_KEY } from '$lib/backend';
	import { t } from '$lib/i18n/index.svelte';

	let visible = $state(false);

	const signInHref = '/api/auth/oidc';

	onMount(() => {
		if (localStorage.getItem(TOKEN_KEY)) {
			goto('/memory');
			return;
		}
		visible = true;
	});

	/* Inside an inverted band nothing may read `fc-fg`: the ink there is `fc-accent-fg`. */
	const invertedDim = 'color-mix(in oklab, var(--color-fc-accent-fg) 25%, transparent)';
	const invertedInk = 'var(--color-fc-accent-fg)';

	const adapters = [
		{ agent: 'pi', file: 'SOUL.md', icon: 'solar:widget-6-bold-duotone' },
		{ agent: 'Claude Code', file: 'CLAUDE.md', icon: 'simple-icons:anthropic' },
		{ agent: 'Gemini CLI', file: 'GEMINI.md', icon: 'simple-icons:googlegemini' },
		{ agent: 'Codex', file: 'AGENTS.md', icon: 'simple-icons:openai' },
		{ agent: 'Cursor', file: '.cursor/rules/', icon: 'simple-icons:cursor' },
		{ agent: 'Copilot', file: 'copilot-instructions.md', icon: 'simple-icons:githubcopilot' }
	];

	const layers = $derived([
		{ icon: icons.folder, title: t('landing.layer.memory.title'), body: t('landing.layer.memory.body') },
		{ icon: icons.shield, title: t('landing.layer.rules.title'), body: t('landing.layer.rules.body') },
		{ icon: icons.bolt, title: t('landing.layer.skills.title'), body: t('landing.layer.skills.body') }
	]);

	const cells = $derived([
		{ name: 'personal', desc: t('landing.cell.personal'), active: true },
		{ name: 'facile', desc: t('landing.cell.team'), active: false },
		{ name: 'client-x', desc: t('landing.cell.client'), active: false }
	]);

	const footerColumns = [
		{
			title: 'Mycelium',
			links: [
				{ href: 'https://github.com/FacileStudio/Mycelium', label: 'GitHub', external: true },
				{ href: 'https://facile.studio', label: 'Facile', external: true }
			]
		}
	];

	function lines(key: Parameters<typeof t>[0]): string[] {
		return t(key).split('|');
	}
</script>

<svelte:head>
	<title>{t('landing.title')}</title>
	<meta name="description" content={t('landing.description')} />
</svelte:head>

{#if visible}
	<div class="min-h-dvh bg-fc-page text-fc-fg">
		<LandingHeader appName="Mycelium" icon="solar:structure-bold-duotone" skipHref="#main">
			{#snippet actions()}
				<SignInWithFacile href={signInHref} data-sveltekit-reload />
			{/snippet}
		</LandingHeader>

		<main id="main">
			<section class="mx-auto max-w-fc-lg px-6 pb-28 pt-36 md:pb-36 md:pt-44">
				<p
					class="mb-6 inline-flex items-center gap-2 rounded-fc-pill px-3.5 py-1 text-fc-xs text-fc-fg-muted"
				>
					<iconify-icon icon={icons.server} width="14" height="14" class="block"></iconify-icon>
					{t('landing.badge')}
				</p>
				<h1 class="max-w-3xl text-fc-3xl font-semibold leading-tight md:text-6xl">
					{lines('landing.heroTitle')[0]}<br />
					<span class="text-fc-fg-muted">{lines('landing.heroTitle')[1]}</span>
				</h1>
				<p class="mt-8 max-w-lg text-fc-md leading-relaxed text-fc-fg-muted">
					{t('landing.heroText')}
				</p>
				<div class="mt-10 flex flex-wrap items-center gap-4">
					<SignInWithFacile href={signInHref} size="lg" data-sveltekit-reload />
					<Button
						href="https://github.com/FacileStudio/Mycelium"
						target="_blank"
						rel="noopener noreferrer"
						variant="outline"
						size="lg"
						iconRight={icons.code}
						class="rounded-fc-md"
					>
						GitHub
					</Button>
				</div>
			</section>

			<section class="bg-fc-accent text-fc-accent-fg">
				<div class="mx-auto max-w-fc-lg px-6 py-28 md:py-36">
					<div class="mb-10 flex items-center gap-3">
						<iconify-icon icon={icons.code} width="24" height="24" class="block"></iconify-icon>
						<span class="text-fc-xl font-semibold">{t('landing.adaptersLabel')}</span>
					</div>
					<h2 class="max-w-2xl text-fc-3xl font-semibold leading-tight md:text-5xl">
						{lines('landing.adaptersTitle')[0]}<br />
						<span class="opacity-50">{lines('landing.adaptersTitle')[1]}</span>
					</h2>
					<WordReveal
						text={t('landing.adaptersText')}
						dimColor={invertedDim}
						revealColor={invertedInk}
						class="mt-8 max-w-xl text-fc-md"
					/>

					<div class="mt-16 grid gap-4 sm:grid-cols-3">
						{#each adapters as adapter (adapter.file)}
							<div class="flex items-center gap-3 rounded-fc-md bg-fc-accent-fg/10 px-4 py-3">
								<iconify-icon icon={adapter.icon} width="20" height="20" class="block shrink-0"></iconify-icon>
								<div class="min-w-0">
									<p class="text-fc-sm font-medium">{adapter.agent}</p>
									<p class="mt-1 truncate font-fc-mono text-fc-xs opacity-60">{adapter.file}</p>
								</div>
							</div>
						{/each}
					</div>
				</div>
			</section>

			<section class="mx-auto max-w-fc-lg px-6 py-28 md:py-36">
				<div class="mb-16 max-w-lg">
					<h2 class="text-fc-3xl font-semibold leading-tight md:text-5xl">
						{lines('landing.layersTitle')[0]}<br /><span class="text-fc-fg-muted">{lines('landing.layersTitle')[1]}</span>
					</h2>
					<p class="mt-4 text-fc-sm text-fc-fg-muted">
						{t('landing.layersText')}
					</p>
				</div>

				<div class="grid gap-4 sm:grid-cols-3">
					{#each layers as layer (layer.title)}
						<Card class="flex flex-col gap-4">
							<span
								class="flex size-10 items-center justify-center rounded-fc-md bg-fc-surface text-fc-fg"
							>
								<iconify-icon icon={layer.icon} width="20" height="20" class="block"
								></iconify-icon>
							</span>
							<h3 class="text-fc-lg font-semibold text-fc-fg">{layer.title}</h3>
							<p class="text-fc-sm leading-relaxed text-fc-fg-muted">{layer.body}</p>
						</Card>
					{/each}
				</div>
			</section>

			<section class="bg-fc-surface">
				<div class="mx-auto max-w-fc-lg px-6 py-28 md:py-36">
					<div class="grid items-center gap-16 md:grid-cols-2">
						<div>
							<div class="mb-10 flex items-center gap-3">
								<iconify-icon icon={icons.usersGroup} width="24" height="24" class="block"
								></iconify-icon>
								<span class="text-fc-xl font-semibold text-fc-fg">{t('landing.cellsLabel')}</span>
							</div>
							<h2 class="text-fc-3xl font-semibold leading-tight md:text-5xl">
								{t('landing.cellsTitle')}
							</h2>
							<WordReveal
								text={t('landing.cellsText')}
								class="mt-8 max-w-lg text-fc-md"
							/>
						</div>
						<div class="grid gap-4">
							{#each cells as cell (cell.name)}
								<div
									class="rounded-fc-md px-5 py-4 {cell.active
										? 'bg-fc-accent text-fc-accent-fg'
										: 'bg-fc-component text-fc-fg'}"
								>
									<p class="font-fc-mono text-fc-sm font-medium">{cell.name}</p>
									<p
										class="mt-1 text-fc-xs {cell.active
											? 'text-fc-accent-fg/60'
											: 'text-fc-fg-muted'}"
									>
										{cell.desc}
									</p>
								</div>
							{/each}
						</div>
					</div>
				</div>
			</section>

			<section class="bg-fc-accent text-fc-accent-fg">
				<div class="mx-auto max-w-fc-lg px-6 py-28 md:py-36">
					<div class="mb-10 flex items-center gap-3">
						<iconify-icon icon={icons.refresh} width="24" height="24" class="block"
						></iconify-icon>
						<span class="text-fc-xl font-semibold">{t('landing.syncLabel')}</span>
					</div>
					<h2 class="text-fc-3xl font-semibold leading-tight md:text-5xl">
						{lines('landing.syncTitle')[0]}<br />
						<span class="opacity-50">{lines('landing.syncTitle')[1]}</span>
					</h2>
					<WordReveal
						text={t('landing.syncText')}
						dimColor={invertedDim}
						revealColor={invertedInk}
						class="mt-8 max-w-xl text-fc-md"
					/>

					<div class="mt-12 rounded-fc-md bg-fc-accent-fg/10 p-6 font-fc-mono text-fc-sm">
						<p class="opacity-50">{t('landing.syncOnServer')}</p>
						<p>$ mycelium serve --port 8420</p>
						<p class="mt-4 opacity-50">{t('landing.syncOnMachine')}</p>
						<p>$ mycelium sync</p>
						<p class="mt-1 opacity-60">&nbsp; ↓ memory/tools/dokploy.md</p>
						<p class="opacity-60">&nbsp; ↑ rules/engineering-ladder.md</p>
						<p class="opacity-60">&nbsp; {t('landing.syncDone')}</p>
					</div>
				</div>
			</section>

			<section class="">
				<div class="mx-auto max-w-fc-lg px-6 py-28 text-center md:py-36">
					<h2 class="text-fc-2xl font-semibold text-fc-fg">
						{t('landing.ctaTitle')}
					</h2>
					<p class="mt-4 text-fc-sm text-fc-fg-muted">
						{t('landing.ctaText')}
					</p>
					<div class="mt-10 flex flex-wrap justify-center gap-3">
						<SignInWithFacile href={signInHref} size="lg" data-sveltekit-reload />
						<Button
							href="https://github.com/FacileStudio/Mycelium"
							target="_blank"
							rel="noopener noreferrer"
							variant="outline"
							size="lg"
							iconRight={icons.code}
							class="rounded-fc-md"
						>
							{t('landing.code')}
						</Button>
					</div>
				</div>
			</section>
		</main>

		<LandingOffers />

		<LandingFooter
			appName="Mycelium"
			icon="solar:structure-bold-duotone"
			tagline={t('landing.footerTagline')}
			columns={footerColumns}
		/>
	</div>
{/if}
