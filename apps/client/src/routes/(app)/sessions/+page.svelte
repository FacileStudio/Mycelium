<script lang="ts">
	import { EmptyState, Tabs, icons } from '@facile/muse';
	import {
		backend,
		type SessionBlock,
		type SessionStats,
		type SessionTimeline,
		type UsageHistory,
		type UsageSnapshot
	} from '$lib/backend';
	import UsageMeter from '$lib/components/UsageMeter.svelte';
	import SessionsClaims from '$lib/components/SessionsClaims.svelte';
	import SessionsHistoryTab from '$lib/components/SessionsHistoryTab.svelte';
	import SessionsLive from '$lib/components/SessionsLive.svelte';
	import SessionsOverviewTab from '$lib/components/SessionsOverviewTab.svelte';
	import { t } from '$lib/i18n/index.svelte';

	const RANGES = $derived([
		{ id: '7d', label: t('range.7d') },
		{ id: '30d', label: t('range.30d') },
		{ id: 'all', label: t('range.all') }
	]);
	const GROUPS = $derived([
		{ id: 'project', label: t('group.project'), icon: icons.folder },
		{ id: 'machine', label: t('group.machine'), icon: icons.server },
		{ id: 'agent', label: t('group.agent'), icon: icons.bolt },
		{ id: 'model', label: t('group.model'), icon: icons.dashboard }
	]);
	const PAGE_TABS = $derived([
		{ id: 'overview', label: t('sessions.tab.overview') },
		{ id: 'history', label: t('sessions.tab.history') },
		{ id: 'usage', label: t('sessions.tab.usage') }
	]);

	let tab = $state('overview');
	let since = $state('7d');
	let by = $state('project');
	let stats: SessionStats | null = $state(null);
	let recent: SessionBlock[] = $state([]);
	let timeline: SessionTimeline | null = $state(null);
	let usage: UsageSnapshot[] = $state([]);
	let usageLog: UsageHistory | null = $state(null);

	/* All time over daily buckets is a thousand pixels of noise; months keep the axis honest. */
	const bucket = $derived(since === 'all' ? 'month' : 'day');
	const bucketUnit = $derived(bucket === 'month' ? 'month' : 'day');

	$effect(() => {
		backend
			.sessionsStats(since, by)
			.then((s) => (stats = s))
			.catch(() => {});
	});

	$effect(() => {
		backend
			.sessionsRecent(15)
			.then((r) => (recent = r ?? []))
			.catch(() => {});
	});

	$effect(() => {
		backend
			.sessionsTimeline(since, bucket, by)
			.then((tl) => (timeline = tl))
			.catch(() => (timeline = null));
	});

	$effect(() => {
		const load = () => {
			backend
				.usageCurrent()
				.then((u) => (usage = u ?? []))
				.catch(() => (usage = []));
			backend
				.usageHistory('7d')
				.then((h) => (usageLog = h))
				.catch(() => (usageLog = null));
		};
		load();
		const timer = setInterval(load, 60_000);
		return () => clearInterval(timer);
	});

	/* $derived.by so `stats` reads as its declared type rather than the null it was
	   initialised with — see the same note in (app)/+layout.svelte. */
	const rows = $derived.by(() => stats?.rows ?? []);
</script>

<div class="flex flex-col gap-10">
	<div class="flex flex-col gap-2">
		<h1 class="text-fc-2xl font-semibold text-fc-fg">{t('sessions.title')}</h1>
		<p class="text-fc-sm text-fc-fg-muted">
			{t('sessions.subtitle')}
		</p>
	</div>

	<SessionsLive />

	<SessionsClaims />

	<section class="flex flex-col gap-4">
		<div class="flex flex-col gap-1">
			<h2 class="text-fc-lg font-semibold text-fc-fg">{t('sessions.totals')}</h2>
			<p class="text-fc-sm text-fc-fg-muted">
				{t('sessions.totalsHelp')}
			</p>
		</div>

		<Tabs items={PAGE_TABS} bind:value={tab} />

		{#if rows.length === 0}
			<EmptyState
				icon={icons.history}
				title={t('sessions.emptyTitle')}
				description={t('sessions.emptyText')}
			/>
		{:else if tab === 'overview'}
			<div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
				<Tabs items={RANGES} bind:value={since} label={t('range.label')} />
				<Tabs items={GROUPS} bind:value={by} label={t('group.label')} />
			</div>

			<SessionsOverviewTab {rows} {timeline} {by} {bucketUnit} />
		{:else if tab === 'history'}
			<div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
				<Tabs items={RANGES} bind:value={since} label={t('range.label')} />
				<Tabs items={GROUPS} bind:value={by} label={t('group.label')} />
			</div>

			<SessionsHistoryTab {rows} {recent} {by} groups={GROUPS} />
		{:else if tab === 'usage'}
			<h2 class="text-fc-lg font-semibold text-fc-fg">{t('sessions.windows')}</h2>
			<p class="text-fc-sm text-fc-fg-muted">
				{t('sessions.windowsHelp')}
			</p>
			<UsageMeter snapshots={usage} history={usageLog} />
		{/if}
	</section>
</div>
