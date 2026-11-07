<script lang="ts">
	import {
		Card,
		DonutChart,
		EmptyState,
		LineChart,
		Sparkline,
		StatCard,
		Tabs,
		icons,
		type ChartSeries
	} from '@facile/muse';
	import {
		backend,
		type SessionStats,
		type SessionTimeline,
		type TokenInfo,
		type LiveSession
	} from '$lib/backend';
	import MachineList from '$lib/components/MachineList.svelte';
	import { t } from '$lib/i18n/index.svelte';
	import {
		bucketLabel,
		columnActive,
		columnTotals,
		formatDuration,
		formatTokens,
		hours,
		periodDelta
	} from '$lib/metrics';

	const ONLINE_WINDOW_MS = 11 * 60 * 1000;

	const RANGES = $derived([
		{ id: '7d', label: t('range.7d') },
		{ id: '30d', label: t('range.30d') },
		{ id: 'all', label: t('range.all') }
	]);

	let tokens: TokenInfo[] = $state([]);
	/*
	 * `/api/tokens` is admin-only, so a plain user gets an error there while the sessions
	 * endpoints answer fine. The pairing list degrades; the charts must not.
	 */
	let tokensDenied = $state(false);
	let since = $state('7d');
	let timeline: SessionTimeline | null = $state(null);
	let stats: SessionStats | null = $state(null);
	let live: LiveSession[] = $state([]);

	const bucket = $derived(since === 'all' ? 'month' : 'day');
	const bucketUnit = $derived(bucket === 'month' ? 'month' : 'day');

	const machines = $derived(
		tokens.filter(
			(m) => m.scope !== 'admin' && m.scope !== 'user' && !m.name.startsWith('session')
		)
	);
	const online = $derived(machines.filter((m) => machineState(m) === 'connected').length);

	$effect(() => {
		const load = () =>
			backend
				.tokensList()
				.then((list) => {
					tokens = list ?? [];
					tokensDenied = false;
				})
				.catch(() => (tokensDenied = true));
		load();
		const id = setInterval(load, 30_000);
		return () => clearInterval(id);
	});

	$effect(() => {
		const load = () =>
			backend
				.sessionsLive()
				.then((l) => (live = l ?? []))
				.catch(() => {});
		load();
		const id = setInterval(load, 30_000);
		return () => clearInterval(id);
	});

	$effect(() => {
		backend
			.sessionsTimeline(since, bucket, 'machine')
			.then((tl) => (timeline = tl))
			.catch(() => (timeline = null));
	});

	$effect(() => {
		backend
			.sessionsStats(since, 'machine')
			.then((s) => (stats = s))
			.catch(() => (stats = null));
	});

	const labels = $derived.by(() => timeline?.labels ?? []);
	const tSeries = $derived.by(() => timeline?.series ?? []);

	const secondsPerBucket = $derived(columnTotals(tSeries.map((s) => s.seconds)));
	const sessionsPerBucket = $derived(columnTotals(tSeries.map((s) => s.sessions)));
	const tokensPerBucket = $derived(columnTotals(tSeries.map((s) => s.tokens_out)));
	const activePerBucket = $derived(columnActive(tSeries.map((s) => s.seconds)));

	const hoursSeries: ChartSeries[] = $derived(
		tSeries.map((s) => ({ name: s.key, data: s.seconds.map(hours) }))
	);

	const rows = $derived.by(() => stats?.rows ?? []);
	const totalSeconds = $derived(rows.reduce((total, r) => total + r.seconds, 0));
	const totalTokensOut = $derived(rows.reduce((total, r) => total + r.tokens_out, 0));
	const shareSlices = $derived(
		[...rows]
			.sort((a, b) => b.seconds - a.seconds)
			.slice(0, 6)
			.map((r) => ({ label: r.key, value: hours(r.seconds) }))
	);

	/* A machine can be working without holding an admin token we are allowed to list, so the
	   headline count falls back to whoever shows up in the session data. */
	const knownMachines = $derived.by(() => {
		if (machines.length > 0) return machines.length;
		const names = new Set<string>([...tSeries.map((s) => s.key), ...live.map((l) => l.machine)]);
		names.delete('Other');
		return names.size;
	});
	const workingNow = $derived(
		new Set(live.filter((l) => l.live && l.machine_online).map((l) => l.machine)).size
	);

	const timeDelta = $derived(periodDelta(secondsPerBucket, bucketUnit));
	const tokensDelta = $derived(periodDelta(tokensPerBucket, bucketUnit));
	const hasTimeline = $derived(labels.length > 0 && tSeries.length > 0);

	/* Three states, not a boolean: a machine that has never checked in needs a different fix
	   from one that checked in this morning and stopped. */
	function machineState(m: TokenInfo): 'connected' | 'idle' | 'never' {
		const then = new Date(m.last_seen).getTime();
		if (!m.last_seen || isNaN(then)) return 'never';
		return Date.now() - then < ONLINE_WINDOW_MS ? 'connected' : 'idle';
	}

</script>

<div class="flex flex-col gap-10">
	<div class="flex flex-col gap-2">
		<h1 class="text-fc-2xl font-semibold text-fc-fg">{t('machines.title')}</h1>
		<p class="text-fc-sm text-fc-fg-muted">
			{t('machines.subtitle')}
		</p>
	</div>

	<section class="flex flex-col gap-4">
		<div class="flex flex-wrap items-center justify-between gap-3">
			<h2 class="text-fc-lg font-semibold text-fc-fg">{t('machines.fleet')}</h2>
			<Tabs items={RANGES} bind:value={since} label={t('range.label')} />
		</div>

		<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4 [&>*]:min-w-0">
			<StatCard label={t('machines.count')} value={knownMachines}>
				<Sparkline data={activePerBucket} class="mt-3" valueFormat={(n) => `${n}`} />
			</StatCard>
			<StatCard
				label={t('machines.workingNow')}
				value={workingNow}
				delta={tokensDenied ? undefined : t('machines.synced', { count: online })}
			>
				<Sparkline
					data={sessionsPerBucket}
					class="mt-3"
					color="var(--color-fc-chart-3)"
					valueFormat={(n) => `${n}`}
				/>
			</StatCard>
			<StatCard label={t('machines.activeTime')} value={formatDuration(totalSeconds)} delta={timeDelta}>
				<Sparkline
					data={secondsPerBucket.map(hours)}
					class="mt-3"
					color="var(--color-fc-chart-2)"
					showLast
				/>
			</StatCard>
			<StatCard label={t('machines.tokensOut')} value={formatTokens(totalTokensOut)} delta={tokensDelta}>
				<Sparkline data={tokensPerBucket} class="mt-3" color="var(--color-fc-chart-5)" />
			</StatCard>
		</div>

		{#if hasTimeline}
			<div class="grid gap-4 lg:grid-cols-3 [&>*]:min-w-0">
				<Card class="flex flex-col gap-4 lg:col-span-2">
					<p class="text-fc-sm font-medium text-fc-fg">
						{t('machines.perBucket', { unit: t(bucketUnit === 'month' ? 'machines.unit.month' : 'machines.unit.day') })}
					</p>
					<LineChart
						series={hoursSeries}
						{labels}
						area
						height={240}
						class="flex-1"
						yFormat={(n) => `${n} h`}
						xFormat={(l) => bucketLabel(l)}
					/>
				</Card>
				<Card class="flex flex-col gap-4">
					<p class="text-fc-sm font-medium text-fc-fg">{t('machines.share')}</p>
					<DonutChart
						data={shareSlices}
						centerLabel={t('machines.active')}
						centerValue={formatDuration(totalSeconds)}
						valueFormat={(n) => `${n} h`}
						class="flex-1"
					/>
				</Card>
			</div>
		{:else}
			<Card class="text-fc-sm text-fc-fg-muted">
				{t('machines.noActivity')}
			</Card>
		{/if}
	</section>

	{#if machines.length === 0}
		<EmptyState icon={icons.server} title={t('machines.emptyTitle')}>
			<p class="text-fc-sm text-fc-fg-muted">
				{#if tokensDenied}
					{t('machines.adminOnly')}
				{:else}
					{t('machines.pairBefore')} <code class="rounded-fc-xs bg-fc-surface px-1.5 py-0.5 font-fc-mono text-fc-xs"
						>mycelium login https://mycelium.facile.studio</code
					> {t('machines.pairAfter')}
				{/if}
			</p>
		</EmptyState>
	{:else}
		<MachineList {machines} {rows} {tSeries} />
	{/if}
</div>
