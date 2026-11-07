<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { icons, toast } from '@facile/muse';
	import { backend } from '$lib/backend';
	import DocumentEditor from '$lib/components/DocumentEditor.svelte';
	import { t } from '$lib/i18n/index.svelte';

	const name = $derived(page.params.name ?? '');
	let content = $state('');
	let loading = $state(true);

	$effect(() => {
		const n = name;
		loading = true;
		backend
			.skillGet(n)
			.then((c) => (content = c))
			.catch((e) => {
				content = '';
				toast.danger(e instanceof Error ? e.message : t('skills.loadOne', { name: n }));
			})
			.finally(() => (loading = false));
	});

	async function save(next: string) {
		await backend.skillSave(name, next);
		content = next;
	}

	async function remove() {
		await backend.skillDelete(name);
		toast.success(t('skills.deleted', { name }));
		goto('/skills');
	}
</script>

<DocumentEditor
	title={name}
	path="skills/{name}.md"
	icon={icons.bolt}
	backHref="/skills"
	backLabel={t('skills.title')}
	{content}
	{loading}
	deleteTitle={t('memory.deleteTitle', { title: name })}
	deleteDescription={t('skills.deleteText')}
	onSave={save}
	onDelete={remove}
/>
