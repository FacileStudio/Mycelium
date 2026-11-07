<script lang="ts">
	import { Badge, Card, icons } from '@facile/muse';
	import { backend, type MemoryIndexState } from '$lib/backend';
	import { formatNumber } from '$lib/i18n/format';
	import { t } from '$lib/i18n/index.svelte';

	/*
	 * Two cadences, because the two states answer different questions. A run in flight is a
	 * moving number a human watches, so it refreshes every 2s; a finished index only changes
	 * when someone writes a page, and 30s is soon enough to notice. Polling stops entirely on
	 * a hidden tab — a dashboard left open in a background tab is the whole reason this file
	 * has a visibility listener rather than a bare interval.
	 */
	const LIVE_MS = 2_000;
	const IDLE_MS = 30_000;

	let index: MemoryIndexState | null = $state(null);
	let error = $state('');
	let inFlight = false;
	let visible = $state(true);

	async function load() {
		if (inFlight) return;
		inFlight = true;
		try {
			index = await backend.memoryIndexStatus();
			error = '';
		} catch (e) {
			error = e instanceof Error ? e.message : t('index.readFailed');
		} finally {
			inFlight = false;
		}
	}

	$effect(() => {
		const sync = () => (visible = !document.hidden);
		sync();
		document.addEventListener('visibilitychange', sync);
		return () => document.removeEventListener('visibilitychange', sync);
	});

	const period = $derived.by(() => (index?.indexing ? LIVE_MS : IDLE_MS));

	$effect(() => {
		if (!visible) return;
		const every = period;
		load();
		const timer = setInterval(load, every);
		return () => clearInterval(timer);
	});

	const percent = $derived.by(() => {
		if (!index || index.total_chunks <= 0) return 0;
		return Math.max(0, Math.min(100, Math.round((index.indexed_chunks / index.total_chunks) * 100)));
	});

	function count(n: number): string {
		return formatNumber(n);
	}

	function rateText(rate: number): string {
		if (!Number.isFinite(rate) || rate <= 0) return t('index.measuring');
		return t('index.rate', { rate: formatNumber(rate, { maximumFractionDigits: rate >= 10 ? 0 : 2 }) });
	}

	/*
	 * An ETA computed from a rate is an estimate, and it should read like one: rounded to
	 * something a human can act on, never to the second the arithmetic produced.
	 */
	function etaText(seconds: number): string {
		if (!Number.isFinite(seconds) || seconds <= 0) return t('index.almostDone');
		if (seconds < 60) return t('index.lessThanMinute');
		const minutes = Math.round(seconds / 60);
		if (minutes < 60) return t('index.minutesLeft', { count: minutes });
		const hours = Math.floor(minutes / 60);
		const rest = minutes % 60;
		return rest === 0
			? t('index.hoursLeft', { hours })
			: t('index.hoursMinutesLeft', { hours, minutes: rest });
	}
</script>

{#if index || error}
	<Card class="flex flex-col gap-3">
		<div class="flex flex-wrap items-center justify-between gap-x-3 gap-y-2">
			<span class="flex min-w-0 items-center gap-2">
				<iconify-icon
					icon={icons.search}
					width="16"
					height="16"
					class="block shrink-0 text-fc-fg-muted"
				></iconify-icon>
				<span class="text-fc-sm font-medium text-fc-fg">{t('index.title')}</span>
			</span>
			{#if index?.indexing}
				<Badge tone="info">{t('index.indexing')}</Badge>
			{:else if index?.enabled}
				<Badge tone="success">{t('index.ready')}</Badge>
			{:else if index}
				<Badge tone="neutral">{t('index.off')}</Badge>
			{/if}
		</div>

		{#if index && !index.enabled}
			<p class="text-fc-sm text-fc-fg-muted">
				{t('index.offText')}
			</p>
		{:else if index}
			{#if index.indexing}
				<div class="flex flex-col gap-1.5">
					<div class="flex items-baseline justify-between gap-3">
						<span class="truncate text-fc-sm text-fc-fg">
							{t('index.progress', { done: count(index.indexed_chunks), total: count(index.total_chunks) })}
						</span>
						<span class="shrink-0 text-fc-sm font-medium tabular-nums text-fc-fg">{percent}%</span>
					</div>
					<div
						class="h-2 w-full overflow-hidden rounded-fc-full bg-fc-surface"
						role="progressbar"
						aria-valuenow={percent}
						aria-valuemin={0}
						aria-valuemax={100}
						aria-label={t('index.embedded')}
					>
						<div
							class="h-full rounded-fc-full bg-fc-accent transition-[width] duration-500"
							style:width="{percent}%"
						></div>
					</div>
					<p class="text-fc-xs text-fc-fg-muted">
						{rateText(index.chunks_per_second)} · {etaText(index.eta_seconds)}{index.pending_paths > 0
							? ` · ${t('index.queued', { count: count(index.pending_paths) })}`
							: ''}
					</p>
				</div>
			{:else}
				<p class="text-fc-sm text-fc-fg-muted">
					{index.pending_paths > 0
						? t('index.doneQueued', { count: count(index.indexed_chunks), queued: count(index.pending_paths) })
						: t('index.doneText', { count: count(index.indexed_chunks) })}
				</p>
			{/if}

			<div class="flex flex-wrap items-center gap-x-3 gap-y-1 font-fc-mono text-fc-xs text-fc-fg-muted">
				{#if index.model?.name}
					<span class="truncate">{index.model.name}</span>
				{/if}
				{#if index.model?.dims}
					<span>{t('index.dims', { count: index.model.dims })}</span>
				{/if}
				{#if index.store}
					<span>{t('index.store', { name: index.store })}</span>
				{/if}
			</div>
		{/if}

		{#if index?.last_error || error}
			<p
				class="flex items-start gap-2 rounded-fc-sm bg-fc-danger/10 px-3 py-2 text-fc-xs text-fc-danger"
			>
				<iconify-icon
					icon={icons.warning}
					width="14"
					height="14"
					class="mt-0.5 block shrink-0"
				></iconify-icon>
				<span class="min-w-0 break-words">{index?.last_error || error}</span>
			</p>
		{/if}
	</Card>
{/if}
