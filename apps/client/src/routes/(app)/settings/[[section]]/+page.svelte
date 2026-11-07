<script lang="ts">
	import { getContext } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import {
		PreferencesSection,
		ProfileSection,
		SettingsLayout,
		icons,
		settingsRoute
	} from '@facile/muse';
	import AdvancedSettings from '$lib/components/settings/AdvancedSettings.svelte';
	import BusSettings from '$lib/components/settings/BusSettings.svelte';
	import ConnectSettings from '$lib/components/settings/ConnectSettings.svelte';
	import TokensSettings from '$lib/components/settings/TokensSettings.svelte';
	import { rememberLocale, t } from '$lib/i18n/index.svelte';
	import { setTheme, theme, type ThemeMode } from '$lib/theme.svelte';
	import type { AuthUser } from '$lib/backend';

	const me = getContext<() => AuthUser | null>('me');

	const sections = [
		{ id: 'profile', label: () => t('settings.tab.profile'), icon: icons.userCircle },
		{ id: 'preferences', label: () => t('settings.tab.preferences'), icon: icons.palette },
		{ id: 'connect', label: () => t('settings.tab.connect'), icon: icons.plug },
		{ id: 'bus', label: () => t('settings.tab.bus'), icon: icons.bolt },
		{ id: 'tokens', label: () => t('settings.tab.api'), icon: icons.key },
		{ id: 'advanced', label: () => t('settings.tab.advanced'), icon: icons.shield }
	];
	const items = $derived(
		sections.map((s) => ({ id: s.id, label: s.label(), icon: s.icon, href: `/settings/${s.id}` }))
	);

	const ALIASES: Record<string, string> = { appearance: 'preferences', api: 'tokens' };

	$effect(() => {
		const requested = page.params.section;
		if (!requested) return;
		const alias = ALIASES[requested];
		if (alias) void goto(`/settings/${alias}`, { replaceState: true });
		else if (!settingsRoute(requested, items).known) void goto('/settings', { replaceState: true });
	});

	const user = $derived(me?.());
</script>

<svelte:head>
	<title>{t('settings.title')} - Mycelium</title>
</svelte:head>

<SettingsLayout
	sections={items}
	active={page.params.section}
	backHref="/settings"
	title={t('settings.title')}
	description={t('settings.subtitle')}
	navLabel={t('settings.sections')}
>
	{#snippet children(section)}
		{#if section === 'profile'}
			<ProfileSection
				name={user?.name || user?.email || 'admin'}
				email={user?.email}
				avatar={user?.avatar_url || undefined}
				role={user?.admin ? 'admin' : undefined}
			/>
		{:else if section === 'preferences'}
			<PreferencesSection
				theme={theme.mode}
				onThemeChange={(next) => setTheme(next as ThemeMode)}
				onLocaleChange={rememberLocale}
			/>
		{:else if section === 'connect'}
			<ConnectSettings />
		{:else if section === 'bus'}
			<BusSettings />
		{:else if section === 'tokens'}
			<TokensSettings />
		{:else if section === 'advanced'}
			<AdvancedSettings />
		{/if}
	{/snippet}
</SettingsLayout>
