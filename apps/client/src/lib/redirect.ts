const KEY = 'mycelium.redirect';
const FALLBACK = '/memory';

export function safeRedirect(target: string | null): string {
	if (!target || !target.startsWith('/') || target.startsWith('//') || target.startsWith('/\\')) return FALLBACK;
	return target;
}

export function rememberRedirect(target: string) {
	try {
		sessionStorage.setItem(KEY, target);
	} catch {}
}

export function takeRedirect(): string {
	try {
		const target = safeRedirect(sessionStorage.getItem(KEY));
		sessionStorage.removeItem(KEY);
		return target;
	} catch {
		return FALLBACK;
	}
}
