<script lang="ts">
	import { Alert, Input, SettingsRow, SettingsSection, Switch } from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';

	let {
		usageAlerts = $bindable(),
		usageThreshold = $bindable(),
		enabled,
		saving,
		onSettle
	}: {
		usageAlerts: boolean;
		usageThreshold: number | string;
		enabled: boolean;
		saving: boolean;
		onSettle: () => void;
	} = $props();
</script>

<SettingsSection
	title={t('alerts.title')}
	description={t('alerts.help')}
>
	{#if !enabled}
		<Alert tone="info" title={t('alerts.emittingOff')}>
			{t('alerts.emittingOffText')}
		</Alert>
	{/if}

	<SettingsRow
		label={t('alerts.on')}
		description={t('alerts.onHelp')}
	>
		<Switch bind:checked={usageAlerts} disabled={!enabled || saving} aria-label={t('alerts.on')} />
	</SettingsRow>

	<SettingsRow
		label={t('alerts.threshold')}
		description={t('alerts.thresholdHelp')}
		for="bus-usage-threshold"
	>
		<Input
			bind:value={usageThreshold}
			id="bus-usage-threshold"
			type="number"
			inputmode="numeric"
			min="1"
			max="100"
			step="1"
			onblur={onSettle}
			disabled={!enabled || !usageAlerts || saving}
			class="w-24 tabular-nums"
		/>
	</SettingsRow>
</SettingsSection>
