import { getActiveSpaceId } from './space.svelte';
import { t } from './i18n/index.svelte';
import { followRegistreSession } from './followRegistre';

const BASE = '/api';

export const TOKEN_KEY = 'mycelium.token';

/**
 * A 401 on any other call means the session is gone, for example revoked by
 * Registre's back-channel logout, so the browser follows Registre's session.
 * Signing out must not: a dead token would sign the person straight back in.
 */
const SIGN_OUT_PATH = '/auth/logout';

export function spaceQuery(sep: '?' | '&' = '?'): string {
	const id = getActiveSpaceId();
	return id ? `${sep}space_id=${encodeURIComponent(id)}` : '';
}

function getToken(): string | null {
	if (typeof window === 'undefined') return null;
	return localStorage.getItem(TOKEN_KEY);
}

type ApiErrorPayload = {
	error?: { message?: string };
};

/**
 * ApiError carries the HTTP status alongside the message.
 *
 * Without it a caller can only see a string, so every failure looks alike —
 * which is how the pool settings page came to answer "you are not an admin"
 * to a server that was simply restarting. A page that wants to special-case
 * one status has to be able to read it.
 */
export class ApiError extends Error {
	readonly status: number;

	constructor(message: string, status: number) {
		super(message);
		this.name = 'ApiError';
		this.status = status;
	}
}

function errorMessage(text: string, status: number): string {
	try {
		const payload = JSON.parse(text) as ApiErrorPayload;
		if (payload?.error?.message) return payload.error.message;
	} catch {}
	if (text && text.length <= 200) return text;
	return t('common.requestFailed', { status });
}

export async function request<T>(method: string, path: string, body?: unknown): Promise<T> {
	const headers: Record<string, string> = {};
	const token = getToken();
	if (token) headers['Authorization'] = `Bearer ${token}`;

	if (body && !(body instanceof FormData)) {
		headers['Content-Type'] = 'application/json';
	}

	const res = await fetch(`${BASE}${path}`, {
		method,
		headers,
		body: body ? (body instanceof FormData ? body : JSON.stringify(body)) : undefined
	});

	if (!res.ok) {
		const text = await res.text();
		if (res.status === 401 && typeof window !== 'undefined') {
			localStorage.removeItem(TOKEN_KEY);
			if (token && path !== SIGN_OUT_PATH) followRegistreSession(location.pathname + location.search);
		}
		throw new ApiError(errorMessage(text.trim(), res.status), res.status);
	}

	const contentType = res.headers.get('content-type');
	if (contentType?.includes('application/json')) {
		return res.json();
	}
	return res.text() as unknown as T;
}
