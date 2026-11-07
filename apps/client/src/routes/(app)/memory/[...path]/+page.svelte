<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { icons, toast } from '@facile/muse';
	import { backend } from '$lib/backend';
	import DocumentEditor from '$lib/components/DocumentEditor.svelte';
	import { t } from '$lib/i18n/index.svelte';

	const rel = $derived(page.params.path ?? '');
	const fullPath = $derived(`memory/${rel}`);
	const title = $derived(rel.split('/').pop()?.replace(/\.md$/, '') ?? rel);
	let content = $state('');
	let loading = $state(true);

	$effect(() => {
		const p = fullPath;
		loading = true;
		backend
			.syncFile(p)
			.then((c) => (content = c))
			.catch((e) => {
				content = '';
				toast.danger(e instanceof Error ? e.message : t('memory.loadFile', { path: p }));
			})
			.finally(() => (loading = false));
	});

	async function save(next: string) {
		await backend.syncFilePut(fullPath, next);
		content = next;
	}

	async function remove() {
		await backend.syncFileDelete(fullPath);
		toast.success(t('memory.deleted', { path: fullPath }));
		goto('/memory');
	}
</script>

<DocumentEditor
	{title}
	path={fullPath}
	icon={icons.folder}
	backHref="/memory"
	backLabel={t('memory.title')}
	{content}
	{loading}
	deleteTitle={t('memory.deleteTitle', { title })}
	deleteDescription={t('memory.deleteText')}
	onSave={save}
	onDelete={remove}
/>
