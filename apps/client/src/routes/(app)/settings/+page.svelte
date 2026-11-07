<script lang="ts">
	import { getContext } from 'svelte';
	import { goto } from '$app/navigation';
	import { Button, ProfileCard, SettingsRow, SettingsSection, icons, toast } from '@facile/muse';
	import { backend, TOKEN_KEY, type AuthUser, type MyceliumStatus } from '$lib/backend';
	import { t } from '$lib/i18n/index.svelte';

	const status = getContext<() => MyceliumStatus | null>('status');

	let me: AuthUser | null = $state(null);
	let loggingOut = $state(false);

	$effect(() => {
		backend
			.authMe()
			.then((u) => (me = u))
			.catch(() => (me = null));
	});

	const meta = $derived.by(() => {
		const s = status?.();
		const rows: { label: string; value: string }[] = [];
		if (s?.machine) rows.push({ label: t('profile.thisMachine'), value: s.machine });
		if (s) rows.push({ label: t('profile.rules'), value: `${s.rules?.length ?? 0}` });
		if (s) rows.push({ label: t('profile.skills'), value: `${s.skills?.length ?? 0}` });
		return rows;
	});

	async function logout() {
		loggingOut = true;
		try {
			if (me) await backend.logout();
		} catch {
			/* The session cookie may already be gone server-side; the local token still has to go. */
		}
		localStorage.removeItem(TOKEN_KEY);
		toast.success(t('profile.loggedOut'));
		goto('/login');
	}
</script>

<div class="flex flex-col gap-10">
	<ProfileCard
		name={me?.name || me?.email || 'admin'}
		email={me?.email}
		avatar={me?.avatar_url || undefined}
		role={me?.admin ? 'admin' : undefined}
		{meta}
	/>

	<SettingsSection
		title={t('profile.account')}
		description={me
			? t('profile.managed')
			: t('profile.sharedPassword')}
	>
		<SettingsRow
			label={t('profile.logout')}
			description={t('profile.logoutHelp')}
		>
			<Button variant="outline" icon={icons.logout} disabled={loggingOut} onclick={logout}>
				{loggingOut ? t('profile.loggingOut') : t('profile.logout')}
			</Button>
		</SettingsRow>
	</SettingsSection>
</div>
