<script lang="ts">
	import { getContext } from 'svelte';
	import { SecretField, SettingsRow, SettingsSection, StatusDot } from '@facile/muse';
	import { backend, type AuthConfig, type MyceliumStatus } from '$lib/backend';
	import { t } from '$lib/i18n/index.svelte';

	const status = getContext<() => MyceliumStatus | null>('status');

	let config: AuthConfig | null = $state(null);

	$effect(() => {
		backend
			.getAuthConfig()
			.then((c) => (config = c))
			.catch(() => (config = null));
	});

	const authMode = $derived.by(() => {
		if (!config) return { tone: 'neutral' as const, label: t('advanced.auth.unknown') };
		if (config.oidc_enabled) return { tone: 'success' as const, label: t('advanced.auth.sso') };
		return { tone: 'warning' as const, label: t('advanced.auth.password') };
	});
</script>

<div class="flex flex-col gap-10">
	<SettingsSection
		title={t('advanced.instance')}
		description={t('advanced.instanceHelp')}
	>
		<SettingsRow label={t('advanced.signIn')} description={t('advanced.signInHelp')}>
			<StatusDot tone={authMode.tone} label={authMode.label} />
		</SettingsRow>

		<SettingsRow label={t('profile.thisMachine')} description={t('advanced.machineHelp')} stacked>
			<SecretField value={status?.()?.machine ?? '—'} sensitive={false} class="w-full" />
		</SettingsRow>

		<SettingsRow label={t('advanced.rulesLoaded')} description={t('advanced.rulesHelp')}>
			<span class="text-fc-sm tabular-nums text-fc-fg">{status?.()?.rules?.length ?? 0}</span>
		</SettingsRow>

		<SettingsRow label={t('advanced.skillsLoaded')} description={t('advanced.skillsHelp')}>
			<span class="text-fc-sm tabular-nums text-fc-fg">{status?.()?.skills?.length ?? 0}</span>
		</SettingsRow>
	</SettingsSection>

	<SettingsSection
		title={t('advanced.cli')}
		description={t('advanced.cliHelp')}
	>
		<SettingsRow label={t('advanced.install')} stacked>
			<SecretField
				value="curl -fsSL https://raw.githubusercontent.com/FacileStudio/Mycelium/main/install.sh | bash"
				sensitive={false}
				class="w-full"
			/>
		</SettingsRow>

		<SettingsRow label={t('advanced.pair')} description={t('advanced.pairHelp')} stacked>
			<SecretField
				value="mycelium login https://mycelium.facile.studio"
				sensitive={false}
				class="w-full"
			/>
		</SettingsRow>
	</SettingsSection>
</div>
