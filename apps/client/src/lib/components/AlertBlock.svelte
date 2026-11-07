<script lang="ts">
	import { richBlock } from './MarkdownMuseRich';
	import type { AlertData } from './MarkdownMuseUtils';

	let { alert }: { alert: AlertData } = $props();

	const toneConfig = {
		pros: { bg: 'bg-fc-success/10', text: 'text-fc-success', emoji: '+' },
		tip: { bg: 'bg-fc-success/10', text: 'text-fc-success', emoji: '+' },
		cons: { bg: 'bg-fc-danger/10', text: 'text-fc-danger', emoji: '-' },
		caution: { bg: 'bg-fc-danger/10', text: 'text-fc-danger', emoji: '-' },
		warning: { bg: 'bg-fc-warning/10', text: 'text-fc-warning', emoji: '!' },
		important: { bg: 'bg-fc-accent/10', text: 'text-fc-accent', emoji: '★' },
		note: { bg: 'bg-fc-info/10', text: 'text-fc-info', emoji: 'i' }
	} as const;

	function getTone(tone: AlertData['tone']) {
		return toneConfig[tone] ?? toneConfig.note;
	}
</script>

<div class="my-3 flex flex-col gap-1.5 rounded-fc-lg {getTone(alert.tone).bg} p-4 text-fc-sm shadow-fc-xs">
	<div class="flex items-center gap-2 font-semibold {getTone(alert.tone).text}">
		<span class="inline-flex size-5 items-center justify-center rounded-full {getTone(alert.tone).bg} text-xs font-bold">{getTone(alert.tone).emoji}</span>
		<span>{alert.title}</span>
	</div>
	<div class="text-fc-xs text-fc-fg leading-relaxed pl-7 flex flex-col gap-1.5 [&>h1]:text-fc-base [&>h1]:font-bold [&>h2]:text-fc-sm [&>h2]:font-bold [&>h3]:text-fc-sm [&>h3]:font-semibold [&>h3]:text-fc-accent [&>ul]:list-disc [&>ul]:pl-5 [&>ul]:space-y-1 [&>ol]:list-decimal [&>ol]:pl-5 [&>ol]:space-y-1">
		{@html richBlock(alert.body)}
	</div>
</div>
