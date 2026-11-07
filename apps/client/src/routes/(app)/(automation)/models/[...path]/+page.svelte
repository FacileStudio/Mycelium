<script lang="ts">
	import { page } from '$app/state';
	import { icons, toast } from '@facile/muse';
	import { backend } from '$lib/backend';
	import CodeView from '$lib/components/CodeView.svelte';
	import { t } from '$lib/i18n/index.svelte';

	const path = $derived(page.params.path ?? '');
	const type = $derived('@' + path.replace(/\.ts$/, ''));
	let content = $state('');
	let loading = $state(true);

	$effect(() => {
		const p = path;
		loading = true;
		backend
			.modelGet(p)
			.then((c) => (content = c))
			.catch((e) => {
				content = '';
				toast.danger(e instanceof Error ? e.message : t('models.loadOne', { name: p }));
			})
			.finally(() => (loading = false));
	});
</script>

<CodeView
	title={type}
	path="extensions/models/{path}"
	icon={icons.code}
	backHref="/models"
	backLabel={t('models.title')}
	{content}
	{loading}
/>
