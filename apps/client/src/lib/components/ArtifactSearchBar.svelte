<script lang="ts">
	import { IconButton, Input, Spinner, icons } from '@facile/muse';
	import { isEditableTarget, MAX_QUERY } from '$lib/artifactSearch';
	import { t } from '$lib/i18n/index.svelte';

	let { value = $bindable(''), busy = false }: { value?: string; busy?: boolean } = $props();

	const id = 'artifact-search';
	let focused = $state(false);

	function field(): HTMLInputElement | null {
		return document.getElementById(id) as HTMLInputElement | null;
	}

	function onWindowKey(e: KeyboardEvent) {
		if (e.key !== '/' || e.metaKey || e.ctrlKey || e.altKey || isEditableTarget(e.target)) return;
		e.preventDefault();
		field()?.focus();
	}

	function onFieldKey(e: KeyboardEvent) {
		if (e.key !== 'Escape') return;
		e.preventDefault();
		if (value) value = '';
		else field()?.blur();
	}

	function clear() {
		value = '';
		field()?.focus();
	}
</script>

<svelte:window onkeydown={onWindowKey} />

<div role="search" class="relative w-full max-w-xl">
	<label for={id} class="sr-only">{t('artifacts.searchLabel')}</label>
	<iconify-icon
		icon={icons.search}
		width="18"
		height="18"
		class="pointer-events-none absolute top-1/2 left-3 block -translate-y-1/2 text-fc-fg-muted"
	></iconify-icon>
	<Input
		{id}
		type="search"
		bind:value
		maxlength={MAX_QUERY}
		autocomplete="off"
		spellcheck="false"
		enterkeyhint="search"
		placeholder={t('artifacts.searchPlaceholder')}
		aria-keyshortcuts="/"
		class="pr-12 pl-10 [&::-webkit-search-cancel-button]:appearance-none"
		onkeydown={onFieldKey}
		onfocus={() => (focused = true)}
		onblur={() => (focused = false)}
	/>
	<div class="absolute top-0 right-0 flex h-11 w-11 items-center justify-center">
		{#if busy}
			<Spinner size="sm" label={t('artifacts.searching')} />
		{:else if value}
			<IconButton variant="ghost" aria-label={t('artifacts.clearSearch')} onclick={clear}>
				<iconify-icon icon={icons.close} width="18" height="18" class="block"></iconify-icon>
			</IconButton>
		{:else if !focused}
			<kbd
				class="hidden rounded-fc-sm bg-fc-surface px-1.5 font-fc-mono text-fc-xs text-fc-fg-muted sm:inline-block"
				aria-hidden="true">/</kbd
			>
		{/if}
	</div>
</div>
