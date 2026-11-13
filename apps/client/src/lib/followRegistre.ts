import { rememberRedirect } from './redirect';

export const SSO_START = '/api/auth/oidc';

const ATTEMPT_KEY = 'mycelium.sso-attempt';
const ATTEMPT_WINDOW_MS = 5 * 60 * 1000;

type Store = Pick<Storage, 'getItem' | 'setItem' | 'removeItem'>;

let leaving = false;

function tabStore(): Store | null {
	try {
		return sessionStorage;
	} catch {
		return null;
	}
}

/** Whether this tab went to the SSO login moments ago and has not come back signed in. */
export function ssoAttemptPending(now = Date.now(), store = tabStore()): boolean {
	try {
		const last = Number(store?.getItem(ATTEMPT_KEY));
		return Boolean(last) && now - last < ATTEMPT_WINDOW_MS;
	} catch {
		return false;
	}
}

/**
 * Claims the one silent sign-in a lost session gets. It answers false, and
 * consumes the marker, when the previous attempt came back still signed out,
 * so the caller shows the sign-in page instead of bouncing through Registre
 * forever. Without storage there is no loop guard, so it answers false too.
 */
export function claimSsoAttempt(now = Date.now(), store = tabStore()): boolean {
	if (!store) return false;
	if (ssoAttemptPending(now, store)) {
		clearSsoAttempt(store);
		return false;
	}
	try {
		store.setItem(ATTEMPT_KEY, String(now));
		return true;
	} catch {
		return false;
	}
}

/** Forgets the attempt once a session is confirmed, so the next loss gets its own. */
export function clearSsoAttempt(store = tabStore()) {
	try {
		store?.removeItem(ATTEMPT_KEY);
	} catch {
		return;
	}
}

/**
 * Sends the browser to the SSO login so Registre's live session, whoever it
 * belongs to now, signs the person straight back in and lands on returnTo.
 * Answers false without navigating when the attempt is refused. A second call
 * while the page is already leaving answers true and spends nothing, so a 401
 * and the layout's guard firing together share the one attempt.
 */
export function followRegistreSession(returnTo: string): boolean {
	if (leaving) return true;
	if (!claimSsoAttempt()) return false;
	leaving = true;
	rememberRedirect(returnTo);
	location.assign(SSO_START);
	return true;
}
