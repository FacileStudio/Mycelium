<script lang="ts">
	import {
		Alert,
		Badge,
		Button,
		Input,
		SecretField,
		SettingsRow,
		SettingsSection,
		Spinner,
		StatusDot,
		Switch,
		icons,
		toast
	} from '@facile/muse';
	import { ApiError, backend, type AntenneSettings, type EmitterStatus } from '$lib/backend';
	import BusMachineOverrides from '$lib/components/BusMachineOverrides.svelte';
	import BusUsageAlerts from '$lib/components/BusUsageAlerts.svelte';
	import { t } from '$lib/i18n/index.svelte';

	let loaded = $state(false);
	let denied = $state(false);
	let saving = $state(false);
	let error = $state('');
	let loadError = $state('');

	let enabled = $state(false);
	let instance = $state('');
	let secret = $state('');
	let userEmail = $state('');
	let emitSince = $state('');
	let usageAlerts = $state(false);
	let usageThreshold: number | string = $state(80);
	let machineEmails: { machine: string; email: string }[] = $state([]);
	let status: EmitterStatus | null = $state(null);
	let envManaged = $state<Record<string, boolean>>({});

	$effect(() => {
		backend
			.settingsGet()
			.then((s) => {
				apply(s.antenne, s.status);
				envManaged = s.env_managed ?? {};
				loaded = true;
			})
			.catch((e) => {
				// Only a 403 means "not an admin". Everything else — the API
				// restarting mid-deploy, a dropped connection, a 500 — used to
				// land here too and get reported as a permission problem,
				// which sends you looking at your account instead of at the
				// server. Say what actually happened.
				if (e instanceof ApiError && e.status === 403) {
					denied = true;
					return;
				}
				// A failed load is not a failed save. They were one variable, so
				// opening the page on a broken response greeted you with "Not
				// saved" for something you never submitted.
				loadError = e instanceof Error ? e.message : t('bus.loadFailed');
			});
	});

	function apply(antenne: AntenneSettings, next: EmitterStatus) {
		enabled = antenne.enabled;
		instance = antenne.instance;
		secret = antenne.secret;
		userEmail = antenne.user_email;
		emitSince = antenne.emit_since ?? '';
		usageAlerts = antenne.usage_alerts ?? false;
		usageThreshold = normalizeThreshold(antenne.usage_threshold);
		machineEmails = Object.entries(antenne.machine_emails ?? {}).map(([machine, email]) => ({
			machine,
			email
		}));
		status = next;
	}

	const DEFAULT_THRESHOLD = 80;

	function normalizeThreshold(raw: number | string | undefined): number {
		const n = Math.round(Number(raw));
		if (!Number.isFinite(n) || n <= 0) return DEFAULT_THRESHOLD;
		return Math.min(100, n);
	}

	function settleThreshold() {
		usageThreshold = normalizeThreshold(usageThreshold);
	}

	/*
	 * Four states, not a boolean. "Not connected" hides a switch that is off, a socket still
	 * opening and a handshake that failed — three situations with three different fixes.
	 */
	const connection = $derived.by(
		(): { tone: 'success' | 'warning' | 'danger' | 'neutral'; label: string; pulse: boolean } => {
			if (!enabled) return { tone: 'neutral', label: t('bus.disabled'), pulse: false };
			if (status?.connected) return { tone: 'success', label: t('bus.connected'), pulse: false };
			if (status?.last_error)
				return { tone: 'danger', label: t('bus.disconnected', { error: status.last_error }), pulse: false };
			return { tone: 'warning', label: t('bus.connecting'), pulse: true };
		}
	);

	async function save(event?: Event) {
		event?.preventDefault();
		saving = true;
		error = '';
		try {
			const emails: Record<string, string> = {};
			for (const row of machineEmails) {
				if (row.machine.trim()) emails[row.machine.trim()] = row.email.trim();
			}
			const s = await backend.settingsSave({
				enabled,
				instance,
				secret,
				user_email: userEmail,
				machine_emails: emails,
				emit_since: emitSince || undefined,
				usage_alerts: usageAlerts,
				usage_threshold: normalizeThreshold(usageThreshold)
			});
			apply(s.antenne, s.status);
			envManaged = s.env_managed ?? {};
			toast.success(t('bus.saved'));
		} catch (e) {
			error = e instanceof Error ? e.message : t('bus.saveFailed');
		} finally {
			saving = false;
		}
	}
</script>

{#if denied}
	<Alert tone="info" title={t('bus.adminOnly')}>
		{t('bus.adminOnlyText')}
	</Alert>
{:else if loadError}
	<Alert tone="danger" title={t('bus.loadFailed')}>{loadError}</Alert>
{:else if !loaded}
	<div class="flex items-center gap-3 text-fc-sm text-fc-fg-muted">
		<Spinner size="sm" label={t('common.loading')} /> {t('common.loadingDots')}
	</div>
{:else}
	<div class="flex flex-col gap-10">
		<SettingsSection
			title={t('bus.title')}
			description={t('bus.help')}
		>
			<SettingsRow label={t('bus.status')} description={t('bus.statusHelp')}>
				<StatusDot tone={connection.tone} label={connection.label} pulse={connection.pulse} />
			</SettingsRow>

			<SettingsRow label={t('bus.emit')} description={t('bus.emitHelp')}>
				<Switch bind:checked={enabled} aria-label={t('bus.emit')} />
			</SettingsRow>

			<SettingsRow label={t('bus.emitted')} description={t('bus.emittedHelp')}>
				<Badge tone="neutral">{t('bus.sent', { count: status?.emitted ?? 0 })}</Badge>
			</SettingsRow>

			<SettingsRow label={t('bus.outbox')} description={t('bus.outboxHelp')}>
				<Badge tone={(status?.pending ?? 0) > 0 ? 'accent' : 'neutral'}>
					{t('bus.pending', { count: status?.pending ?? 0 })}
				</Badge>
			</SettingsRow>

			{#if usageAlerts}
				<SettingsRow
					label={t('bus.alertOutbox')}
					description={t('bus.alertOutboxHelp')}
				>
					<Badge tone={(status?.usage_alerts_pending ?? 0) > 0 ? 'accent' : 'neutral'}>
						{t('bus.pending', { count: status?.usage_alerts_pending ?? 0 })}
					</Badge>
				</SettingsRow>
			{/if}

			{#if emitSince}
				<SettingsRow
					label={t('bus.since')}
					description={t('bus.sinceHelp')}
					stacked
				>
					<SecretField value={emitSince} sensitive={false} class="w-full" />
				</SettingsRow>
			{/if}
		</SettingsSection>

		<BusUsageAlerts
			bind:usageAlerts
			bind:usageThreshold
			{enabled}
			{saving}
			onSettle={settleThreshold}
		/>

		<SettingsSection
			title={t('bus.connection')}
			description={t('bus.connectionHelp')}
		>
			<SettingsRow
				label={t('bus.url')}
				description={envManaged.instance
					? t('bus.urlPinned')
					: t('bus.urlHelp')}
				for="bus-url"
				stacked
			>
				<Input
					bind:value={instance}
					id="bus-url"
					placeholder="https://antenne.facile.studio"
					disabled={saving || envManaged.instance}
				/>
			</SettingsRow>

			<SettingsRow
				label={t('bus.secret')}
				description={envManaged.secret
					? t('bus.secretPinned')
					: t('bus.secretHelp')}
				stacked
			>
				<SecretField
					bind:value={secret}
					editable={!envManaged.secret}
					disabled={saving || envManaged.secret}
					class="w-full"
				/>
			</SettingsRow>

			<SettingsRow
				label={t('bus.email')}
				description={envManaged.user_email
					? t('bus.emailPinned')
					: t('bus.emailHelp')}
				for="bus-email"
				stacked
			>
				<Input
					bind:value={userEmail}
					id="bus-email"
					type="email"
					placeholder="you@facile.studio"
					disabled={saving || envManaged.user_email}
				/>
			</SettingsRow>

			{#if error}
				<Alert tone="danger" title={t('bus.notSaved')}>{error}</Alert>
			{/if}

			<div class="flex flex-wrap gap-2 pt-2">
				<Button icon={icons.plug} disabled={saving} onclick={save}>
					{saving ? t('common.saving') : t('bus.saveConnect')}
				</Button>
			</div>
		</SettingsSection>

		<BusMachineOverrides bind:machineEmails />
	</div>
{/if}
