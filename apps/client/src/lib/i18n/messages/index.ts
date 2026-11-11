import type { Catalogs, Message } from '@facile/muse';
import { automation } from './automation';
import { sessions } from './sessions';
import { spaces } from './spaces';
import { bus } from './bus';
import { settings } from './settings';
import { viewer } from './viewer';
import { common } from './common';
import { memory } from './memory';
import { shell } from './shell';
import { pending } from './pending';
import type { Row } from '../row';

const rows = {
	...common,
	...shell,
	...memory,
	...automation,
	...sessions,
	...spaces,
	...settings,
	...bus,
	...viewer,
	...pending
};

export type AppKey = keyof typeof rows;
type AppCatalog = Record<AppKey, Message>;

function column(index: number): AppCatalog {
	const entries = Object.entries(rows as Record<string, Row>).map(([key, row]) => [key, row[index]]);
	return Object.fromEntries(entries) as AppCatalog;
}

/** The app's catalogs: English is the typed base, and every row carries all four locales. */
export const catalogs: Catalogs<AppCatalog> = {
	en: column(0),
	fr: column(1),
	es: column(2),
	de: column(3)
};
