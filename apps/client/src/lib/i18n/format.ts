import { locale } from './index.svelte';

/** Formats a date or instant with `Intl.DateTimeFormat` in the current locale. */
export function formatDate(value: Date | string | number, options: Intl.DateTimeFormatOptions = {}): string {
	return new Intl.DateTimeFormat(locale(), options).format(new Date(value));
}

/** Formats a number with `Intl.NumberFormat` in the current locale. */
export function formatNumber(value: number, options: Intl.NumberFormatOptions = {}): string {
	return new Intl.NumberFormat(locale(), options).format(value);
}
