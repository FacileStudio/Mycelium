/*
 * The runnable check for the SSO loop guard, run by scripts/check.sh.
 *
 *   bun src/lib/followRegistre.check.ts
 */
import {
	SSO_START,
	claimSsoAttempt,
	clearSsoAttempt,
	followRegistreSession,
	ssoAttemptPending
} from './followRegistre';

function memoryStore(): Storage {
	const values = new Map<string, string>();
	return {
		get length() {
			return values.size;
		},
		clear: () => values.clear(),
		key: (index: number) => [...values.keys()][index] ?? null,
		getItem: (key: string) => values.get(key) ?? null,
		setItem: (key: string, value: string) => void values.set(key, value),
		removeItem: (key: string) => void values.delete(key)
	};
}

const assigned: string[] = [];
const tab = memoryStore();
Object.assign(globalThis, {
	sessionStorage: tab,
	location: { assign: (url: string) => assigned.push(url) }
});


function want(label: string, got: unknown, expected: unknown) {
	if (JSON.stringify(got) !== JSON.stringify(expected)) {
		throw new Error(`${label}: got ${JSON.stringify(got)}, want ${JSON.stringify(expected)}`);
	}
}

const store = memoryStore();
want('first loss redirects', claimSsoAttempt(1_000, store), true);
want('attempt is pending', ssoAttemptPending(2_000, store), true);
want('second loss gives up', claimSsoAttempt(2_000, store), false);
want('giving up consumes the marker', ssoAttemptPending(2_000, store), false);
want('the next loss redirects again', claimSsoAttempt(3_000, store), true);
clearSsoAttempt(store);
want('a confirmed session resets', claimSsoAttempt(4_000, store), true);
want('a stale attempt is forgotten', claimSsoAttempt(10 * 60 * 1000, store), true);
want('no storage, no redirect', claimSsoAttempt(1, null), false);

want('follow leaves for the SSO login', followRegistreSession('/artifacts/abc?space_id=s'), true);
want('a second call while leaving spends nothing', followRegistreSession('/memory'), true);
want('one navigation', assigned, [SSO_START]);
want('the path survives the trip', tab.getItem('mycelium.redirect'), '/artifacts/abc?space_id=s');

console.log('followRegistre: ok');
