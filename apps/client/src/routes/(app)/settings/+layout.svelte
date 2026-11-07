<script lang="ts">
	import { page } from '$app/state';
	import { Divider, PageTransition, Tabs, icons } from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';

	let { children } = $props();

	/*
	 * The section is a real route, not local state: /settings/bus opens on Bus, reload keeps
	 * you there and Back walks the sections. Profile lives at the bare /settings because that
	 * is where the sidebar's user card points.
	 */
	const sections = $derived([
		{ id: 'profile', label: t('settings.tab.profile'), icon: icons.userCircle, href: '/settings' },
		{ id: 'appearance', label: t('settings.tab.appearance'), icon: icons.palette, href: '/settings/appearance' },
		{ id: 'api', label: t('settings.tab.api'), icon: icons.key, href: '/settings/api' },
		{ id: 'bus', label: t('settings.tab.bus'), icon: icons.plug, href: '/settings/bus' },
		{ id: 'advanced', label: t('settings.tab.advanced'), icon: icons.shield, href: '/settings/advanced' }
	]);

	const active = $derived(
		sections.find((s) => s.href !== '/settings' && page.url.pathname.startsWith(s.href))?.id ??
			'profile'
	);
</script>

<div class="flex flex-col gap-8">
	<div class="flex flex-col gap-2">
		<h1 class="text-fc-2xl font-semibold text-fc-fg">{t('settings.title')}</h1>
		<p class="text-fc-sm text-fc-fg-muted">
			{t('settings.subtitle')}
		</p>
	</div>

	<!-- gap-4, not tighter: pulled under the strip the rule reads as an underline welded to
	     the active pill. -->
	<div class="flex flex-col gap-4">
		<Tabs items={sections} value={active} label={t('settings.sections')} />
		<Divider class="my-0" />
	</div>

	<PageTransition key={active} distance={8} duration={0.25}>
		{@render children()}
	</PageTransition>
</div>
