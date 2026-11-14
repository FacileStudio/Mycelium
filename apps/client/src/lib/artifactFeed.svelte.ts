import { toast } from '@facile/muse';
import { backend, type ArtifactHit } from './backend';
import { latestOnly, normaliseQuery } from './artifactSearch';

function message(e: unknown): string {
	return e instanceof Error ? e.message : String(e);
}

/**
 * The artifacts page's list: every artifact while the query is empty, the search route's pages
 * otherwise. Only the newest request writes, so typing fast never shows an older answer.
 */
export class ArtifactFeed {
	items = $state<ArtifactHit[]>([]);
	total = $state(0);
	nextCursor = $state('');
	query = $state('');
	loading = $state(true);
	ready = $state(false);
	loadingMore = $state(false);
	error = $state('');
	#gate = latestOnly();

	async load(raw: string) {
		const q = normaliseQuery(raw);
		const ticket = this.#gate.next();
		this.loading = true;
		this.error = '';
		try {
			const page = q ? await backend.artifactsSearch(q) : null;
			const all = q ? [] : ((await backend.artifactsAll()) ?? []);
			if (!this.#gate.isCurrent(ticket)) return;
			this.query = q;
			this.items = page ? (page.results ?? []) : all;
			this.total = page ? page.total : all.length;
			this.nextCursor = page?.next_cursor ?? '';
		} catch (e) {
			if (this.#gate.isCurrent(ticket)) this.error = message(e);
		} finally {
			if (this.#gate.isCurrent(ticket)) {
				this.loading = false;
				this.ready = true;
			}
		}
	}

	async more() {
		if (!this.nextCursor || this.loadingMore) return;
		const ticket = this.#gate.next();
		this.loadingMore = true;
		try {
			const page = await backend.artifactsSearch(this.query, this.nextCursor);
			if (!this.#gate.isCurrent(ticket)) return;
			this.items = [...this.items, ...(page.results ?? [])];
			this.nextCursor = page.next_cursor ?? '';
		} catch (e) {
			if (this.#gate.isCurrent(ticket)) toast.danger(message(e));
		} finally {
			this.loadingMore = false;
		}
	}
}
