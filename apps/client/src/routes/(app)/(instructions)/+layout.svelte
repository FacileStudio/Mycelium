<script lang="ts">
	import { page } from '$app/state';
	import { Divider, PageTransition, Tabs, icons } from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';

	let { children } = $props();

	const sections = $derived([
		{ id: 'rules', label: t('rules.title'), icon: icons.shield, href: '/rules' },
		{ id: 'skills', label: t('skills.title'), icon: icons.bolt, href: '/skills' }
	]);

	const active = $derived(sections.find((s) => page.url.pathname.startsWith(s.href))?.id ?? 'rules');
</script>

<div class="flex flex-col gap-8">
	<div class="flex flex-col gap-2">
		<h1 class="text-fc-2xl font-semibold text-fc-fg">{t('instructions.title')}</h1>
		<p class="text-fc-sm text-fc-fg-muted">
			{t('instructions.subtitle')}
		</p>
	</div>

	<div class="flex flex-col gap-4">
		<Tabs items={sections} value={active} label={t('instructions.sections')} />
		<Divider class="my-0" />
	</div>

	<PageTransition key={active} distance={8} duration={0.25}>
		{@render children()}
	</PageTransition>
</div>
