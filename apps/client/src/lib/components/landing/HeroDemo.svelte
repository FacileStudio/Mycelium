<script lang="ts">
	import { Caret, Stage, Stream, Terminal, grow, type Act } from '@facile/mise-en-scene';
	import { t } from '$lib/i18n/index.svelte';

	const COMMAND = 'mycelium memory search "backup docker group"';

	const RESULTS = [
		{ path: 'conventions/backups.md:16', text: 'The dump runs through docker exec, so it needs the docker group.' },
		{ path: 'projects/ops-backups.md:49', text: 'A daily cron at 03:15 runs the dump, then the check.' },
		{ path: 'tools/postgres.md:8', text: 'Test a dump with pg_restore into a throwaway container.' }
	];

	const prompt = $derived(t('landing.demo.prompt'));
	const answer = $derived(t('landing.demo.answer'));

	let asked = $state('');
	let command = $state('');
	let calling = $state(false);
	let shown = $state(0);
	let words = $state(0);

	function reset() {
		asked = '';
		command = '';
		calling = false;
		shown = 0;
		words = 0;
	}

	function poster() {
		asked = prompt;
		calling = true;
		command = COMMAND;
		shown = RESULTS.length;
		words = answer.split(' ').length;
	}

	async function script(act: Act) {
		await act.wait(700);
		await act.type(prompt, (v) => (asked = v), 45);
		await act.wait(700);
		calling = true;
		await act.type(COMMAND, (v) => (command = v), 28);
		await act.wait(500);
		for (let i = 1; i <= RESULTS.length; i++) {
			shown = i;
			await act.wait(320);
		}
		await act.wait(700);
		await act.count(answer.split(' ').length, (v) => (words = v), 70);
		await act.wait(2600);
		reset();
		await act.wait(400);
	}
</script>

<Stage
	width={960}
	height={300}
	compact={{ below: 640, width: 358, height: 410 }}
	label={t('landing.demo.label')}
	{script}
	{reset}
	{poster}
>
	{#snippet children({ compact })}
		<Terminal class="hd {compact ? 'hd-compact' : ''}" style={compact ? 'top:12px;left:4px;right:4px;bottom:20px' : 'top:12px;left:0;right:0;bottom:36px'}>
			<div class="hd-head"><span class="hd-dot"></span>claude <span class="hd-dim">~/Code/api</span></div>
			<div class="hd-line"><span class="hd-dim">&gt;</span> {asked}{#if !calling}<Caret />{/if}</div>
			{#if calling}
				<div class="hd-block" transition:grow={{ duration: 260 }}>
					<div class="hd-line hd-call"><span class="hd-ok">●</span> <span>{command}</span>{#if shown === 0}<Caret />{/if}</div>
					{#each RESULTS.slice(0, shown) as result, i (result.path)}
						<div class="hd-result" transition:grow={{ duration: 280 }}>
							<span class="hd-dim">{i === 0 ? '⎿' : ''}</span>
							<span class={i === 0 ? 'hd-path' : 'hd-dim'}>{result.path}</span>
							<span class="hd-text {i === 0 ? '' : 'hd-dim'}">{result.text}</span>
						</div>
					{/each}
				</div>
			{/if}
			{#if words > 0}
				<div class="hd-answer" transition:grow={{ duration: 260 }}>
					<span class="hd-mark">●</span>
					<Stream text={answer} count={words} class="hd-stream" />
				</div>
			{/if}
		</Terminal>
	{/snippet}
</Stage>

<style>
	:global(.hd) {
		padding: 20px 24px;
		font-size: 14px;
		line-height: 24px;
	}
	:global(.hd.hd-compact) {
		padding: 14px 14px;
		font-size: 12px;
		line-height: 19px;
	}
	.hd-head {
		display: flex;
		align-items: center;
		gap: 8px;
		margin-bottom: 14px;
	}
	.hd-dot {
		width: 8px;
		height: 8px;
		border-radius: 999px;
		background: var(--color-fc-success);
	}
	.hd-line {
		white-space: pre-wrap;
		overflow-wrap: anywhere;
	}
	.hd-block {
		padding-top: 12px;
	}
	.hd-call {
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}
	:global(.hd-compact) .hd-call {
		white-space: pre-wrap;
		overflow-wrap: anywhere;
	}
	.hd-result {
		display: grid;
		grid-template-columns: 2ch 30ch minmax(0, 1fr);
		gap: 8px;
		white-space: nowrap;
	}
	:global(.hd-compact) .hd-result {
		grid-template-columns: 2ch minmax(0, 1fr);
		row-gap: 0;
		padding-bottom: 6px;
	}
	:global(.hd-compact) .hd-result .hd-text {
		grid-column: 2;
	}
	.hd-text {
		overflow: hidden;
		text-overflow: ellipsis;
	}
	.hd-answer {
		display: grid;
		grid-template-columns: 2ch minmax(0, 1fr);
		gap: 8px;
		padding-top: 16px;
	}
	.hd-answer :global(.hd-stream) {
		margin: 0;
		font-family: var(--font-fc-body);
	}
	.hd-path {
		font-weight: 700;
	}
	.hd-dim {
		color: var(--color-fc-on-inverse-muted);
	}
	.hd-ok {
		color: var(--color-fc-success);
	}
	.hd-mark {
		color: var(--color-fc-on-inverse);
	}
</style>
