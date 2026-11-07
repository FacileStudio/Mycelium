<script lang="ts">
	import { Button, SecretField, SettingsSection, icons, toast } from '@facile/muse';
	import { AGENT_PROMPT } from '$lib/agentPrompt';
	import { t } from '$lib/i18n/index.svelte';

	const endpoint = $derived(
		typeof window === 'undefined' ? '/api' : `${window.location.origin}/api`
	);

	async function copyPrompt() {
		try {
			await navigator.clipboard.writeText(AGENT_PROMPT);
			toast.success(t('api.promptCopied'));
		} catch {
			toast.danger(t('api.clipboardDenied'));
		}
	}
</script>

<div class="flex flex-col gap-10">
	<SettingsSection title={t('api.endpoint')} description={t('api.endpointHelp')}>
		<SecretField value={endpoint} sensitive={false} label={t('api.baseUrl')} />
	</SettingsSection>

	<SettingsSection title={t('api.connect')} description={t('api.connectHelp')} bare>
		{#snippet actions()}
			<Button variant="outline" icon={icons.copy} onclick={copyPrompt}>{t('api.copyPrompt')}</Button>
		{/snippet}

		<pre
			class="max-h-96 overflow-auto rounded-fc-lg bg-fc-surface p-4 font-fc-mono text-fc-xs leading-relaxed whitespace-pre-wrap text-fc-fg">{AGENT_PROMPT}</pre>
		<p class="text-fc-xs text-fc-fg-muted">
			{t('api.promptNote')}
		</p>
	</SettingsSection>
</div>
