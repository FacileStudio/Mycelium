<script lang="ts">
	import '../app.css';
	import { browser } from '$app/environment';
	import { Toaster, provideLocale } from '@facile/muse';
	import { bindLocale, localeSource } from '$lib/i18n/index.svelte';
	import { applyStoredTheme } from '$lib/theme.svelte';

	let { children } = $props();

	const i18n = provideLocale(localeSource);
	bindLocale(() => i18n.locale);

	/*
	 * muse renders <iconify-icon> custom elements; they stay inert until the element is
	 * registered. Loading it here rather than from a CDN keeps the SPA self-hosted, which is
	 * the whole point of Mycelium.
	 */
	if (browser) {
		import('iconify-icon');
		applyStoredTheme();
	}
</script>

{@render children()}

<!-- One Toaster for the whole app, outside the router, so a navigation cannot unmount a
     toast mid-flight. The extra bottom padding clears MobileNav on a phone. -->
<Toaster class="pb-28 md:pb-6" />
