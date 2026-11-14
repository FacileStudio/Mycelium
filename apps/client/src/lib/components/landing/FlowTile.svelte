<script lang="ts">
	import { Caret, Stage, Terminal, Tile, grow, type Act } from '@facile/mise-en-scene';
	import { t } from '$lib/i18n/index.svelte';

	const COMMAND = 'mycelium flow run ship-check';

	const STEPS = [
		{ name: 'git', out: 'main, clean, in step with origin', time: '0.4s' },
		{ name: 'test', out: 'go test ./... ok', time: '38.2s' },
		{ name: 'probe', out: 'mycelium.facile.studio 200', time: '0.3s' }
	];

	let typed = $state('');
	let running = $state(false);
	let started = $state(0);
	let finished = $state(0);
	let done = $state(false);

	function reset() {
		typed = '';
		running = false;
		started = 0;
		finished = 0;
		done = false;
	}

	function poster() {
		typed = COMMAND;
		running = true;
		started = STEPS.length;
		finished = STEPS.length;
		done = true;
	}

	async function script(act: Act) {
		await act.wait(600);
		await act.type(COMMAND, (v) => (typed = v), 60);
		await act.wait(400);
		running = true;
		for (let i = 1; i <= STEPS.length; i++) {
			await act.wait(500);
			started = i;
			await act.wait(i === 2 ? 1300 : 600);
			finished = i;
		}
		await act.wait(400);
		done = true;
		await act.wait(2400);
		reset();
		await act.wait(400);
	}
</script>

<Tile title={COMMAND} body={t('landing.cli.flow')} level={3}>
	<Stage width={560} height={230} compact={{ below: 480, width: 340, height: 230 }} {script} {reset} {poster}>
		{#snippet children({ compact })}
			<Terminal style="top:24px;left:24px;right:24px;bottom:24px">
				<div class="ft-line"><span class="ft-dim">$</span> {typed}{#if !running}<Caret />{/if}</div>
				{#if running}
					<div class="ft-line" transition:grow={{ duration: 260 }}>
						<span class="ft-step">▸</span> Running ship-check <span class="ft-dim">(3 steps)</span>
					</div>
				{/if}
				{#each STEPS.slice(0, started) as step, i (step.name)}
					<div class="ft-row" class:ft-compact={compact} transition:grow={{ duration: 300 }}>
						<span class="ft-dim">[{step.name}]</span>
						<span class="ft-out">{step.out}</span>
						{#if i < finished}
							<span class="ft-ok">✓ {step.time}</span>
						{:else}
							<span class="ft-dim">…</span>
						{/if}
					</div>
				{/each}
				{#if done}
					<div class="ft-line" transition:grow={{ duration: 260 }}>
						<span class="ft-ok">✓</span> ship-check: 3 steps in 38.9s
					</div>
				{/if}
			</Terminal>
		{/snippet}
	</Stage>
</Tile>

<style>
	.ft-line,
	.ft-row {
		line-height: 24px;
		white-space: nowrap;
	}
	.ft-row {
		display: grid;
		grid-template-columns: 7ch minmax(0, 1fr) 7ch;
		gap: 12px;
	}
	.ft-compact {
		gap: 8px;
	}
	.ft-out {
		overflow: hidden;
		text-overflow: ellipsis;
	}
	.ft-dim {
		color: var(--color-fc-on-inverse-muted);
	}
	.ft-step {
		color: var(--color-fc-info);
	}
	.ft-ok {
		color: var(--color-fc-success);
		text-align: right;
	}
</style>
