<script lang="ts">
	import { Button, Input, SettingsRow, SettingsSection, icons } from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';

	let {
		machineEmails = $bindable()
	}: {
		machineEmails: { machine: string; email: string }[];
	} = $props();
</script>

<SettingsSection
	title={t('overrides.title')}
	description={t('overrides.help')}
>
	{#each machineEmails as row, i (i)}
		<SettingsRow label={t('overrides.row', { n: i + 1 })} stacked>
			<div class="flex w-full flex-col gap-2 sm:flex-row">
				<Input bind:value={row.machine} placeholder={t('overrides.machine')} aria-label={t('overrides.machineLabel')} />
				<Input
					bind:value={row.email}
					type="email"
					placeholder={t('overrides.email')}
					aria-label={t('overrides.emailLabel')}
					class="flex-1"
				/>
				<Button
					variant="ghost-danger"
					icon={icons.remove}
					aria-label={t('overrides.remove', { n: i + 1 })}
					onclick={() => (machineEmails = machineEmails.filter((_, j) => j !== i))}
				>
					{t('members.remove')}
				</Button>
			</div>
		</SettingsRow>
	{/each}

	<SettingsRow
		label={t('overrides.add')}
		description={t('overrides.addHelp')}
	>
		<Button
			variant="secondary"
			icon={icons.plus}
			onclick={() => (machineEmails = [...machineEmails, { machine: '', email: '' }])}
		>
			{t('overrides.addButton')}
		</Button>
	</SettingsRow>
</SettingsSection>
