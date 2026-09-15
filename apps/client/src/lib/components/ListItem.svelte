<script lang="ts">
	import { richInline } from './MarkdownMuseRich';

	let { kind, clean }: { kind: 'pro' | 'con' | 'warn' | 'info' | 'normal'; clean: string } = $props();

	const cfg: Record<string, { li: string; badge: string; icon: string }> = {
		pro: { li: 'flex items-start gap-2.5 rounded-fc-md bg-fc-success/5 px-3 py-1.5 text-fc-sm border border-fc-success/20', badge: 'mt-0.5 inline-flex size-4 shrink-0 items-center justify-center rounded-full bg-fc-success/20 text-fc-success text-xs font-bold', icon: '+' },
		con: { li: 'flex items-start gap-2.5 rounded-fc-md bg-fc-danger/5 px-3 py-1.5 text-fc-sm border border-fc-danger/20', badge: 'mt-0.5 inline-flex size-4 shrink-0 items-center justify-center rounded-full bg-fc-danger/20 text-fc-danger text-xs font-bold', icon: '-' },
		warn: { li: 'flex items-start gap-2.5 rounded-fc-md bg-fc-warning/5 px-3 py-1.5 text-fc-sm border border-fc-warning/20', badge: 'mt-0.5 inline-flex size-4 shrink-0 items-center justify-center rounded-full bg-fc-warning/20 text-fc-warning text-xs font-bold', icon: '!' },
		info: { li: 'flex items-start gap-2.5 rounded-fc-md bg-fc-info/5 px-3 py-1.5 text-fc-sm border border-fc-info/20', badge: 'mt-0.5 inline-flex size-4 shrink-0 items-center justify-center rounded-full bg-fc-info/20 text-fc-info text-xs font-bold', icon: 'i' }
	} as const;
</script>

{#if kind === 'normal'}
	<li class="flex items-start gap-2.5 pl-2 text-fc-sm text-fc-fg-muted">
		<span class="mt-2 size-1.5 shrink-0 rounded-full bg-fc-primary/80"></span>
		<div class="leading-relaxed">
			{@html richInline(clean)}
		</div>
	</li>
{:else}
	{@const c = cfg[kind]}
	<li class={c.li}>
		<span class={c.badge}>{c.icon}</span>
		<div class="leading-relaxed text-fc-fg">
			{@html richInline(clean)}
		</div>
	</li>
{/if}
