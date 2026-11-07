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
			.ruleGet(n)
			.then((c) => (content = c))
			.catch((e) => {
				content = '';
				toast.danger(e instanceof Error ? e.message : t('rules.loadOne', { name: n }));
			})
			.finally(() => (loading = false));
	});

	async function save(next: string) {
		await backend.ruleSave(name, next);
		content = next;
	}

	async function remove() {
		await backend.ruleDelete(name);
		toast.success(t('rules.deleted', { name }));
		goto('/rules');
	}
</script>

<DocumentEditor
	title={name}
	path="rules/{name}.md"
	icon={icons.shield}
	backHref="/rules"
	backLabel={t('rules.title')}
	{content}
	{loading}
	deleteTitle={t('memory.deleteTitle', { title: name })}
	deleteDescription={t('rules.deleteText')}
	onSave={save}
	onDelete={remove}
/>
