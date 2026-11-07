<script lang="ts">
	import { OptionCards, SettingsRow, SettingsSection, icons } from '@facile/muse';
	import { setTheme, theme, type ThemeMode } from '$lib/theme.svelte';
	import { t } from '$lib/i18n/index.svelte';

	const modes = $derived([
		{ value: 'system', label: t('appearance.system'), icon: icons.monitor },
		{ value: 'light', label: t('appearance.light'), icon: icons.sun },
		{ value: 'dark', label: t('appearance.dark'), icon: icons.moon }
	]);

	let mode = $state<string>(theme.mode);

	$effect(() => {
		setTheme(mode as ThemeMode);
	});
</script>

<div class="flex flex-col gap-10">
	<SettingsSection
		title={t('appearance.theme')}
		description={t('appearance.themeHelp')}
	>
		<SettingsRow
			label={t('appearance.scheme')}
			description={t('appearance.schemeHelp')}
			stacked
		>
			<OptionCards options={modes} bind:value={mode} name="theme-mode" label={t('appearance.scheme')} />
		</SettingsRow>
	</SettingsSection>
</div>
