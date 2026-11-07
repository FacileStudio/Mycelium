<script lang="ts">
	import { page } from '$app/state';
	import { Divider, PageTransition, Tabs, icons } from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';

	let { children } = $props();

	const sections = $derived([
		{ id: 'flows', label: t('flows.title'), icon: icons.plug, href: '/flows' },
		{ id: 'models', label: t('models.title'), icon: icons.code, href: '/models' }
	]);

	const active = $derived(
		sections.find((s) => page.url.pathname.startsWith(s.href))?.id ?? 'flows'
	);
</script>

<div class="flex flex-col gap-8">
	<div class="flex flex-col gap-2">
		<h1 class="text-fc-2xl font-semibold text-fc-fg">{t('automation.title')}</h1>
		<p class="text-fc-sm text-fc-fg-muted">
			{t('automation.subtitle')}
		</p>
	</div>

	<div class="flex flex-col gap-4">
		<Tabs items={sections} value={active} label={t('automation.sections')} />
		<Divider class="my-0" />
	</div>

	<PageTransition key={active} distance={8} duration={0.25}>
		{@render children()}
	</PageTransition>
</div>
