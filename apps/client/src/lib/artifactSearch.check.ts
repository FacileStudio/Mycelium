/*
 * The runnable check for the artifact search helpers, run by `bun run test`.
 *
 *   bun src/lib/artifactSearch.check.ts
 */
import { artifactSearchPath, highlightParts, isEditableTarget, latestOnly, MAX_QUERY } from './artifactSearch';

function want(label: string, got: unknown, expected: unknown) {
	const same = JSON.stringify(got) === JSON.stringify(expected);
	if (!same) throw new Error(`${label}: got ${JSON.stringify(got)}, want ${JSON.stringify(expected)}`);
}

want('path', artifactSearchPath('  encre table '), '/artifacts/search?all=true&q=encre+table');
want('cursor', artifactSearchPath('a&b', 'c2'), '/artifacts/search?all=true&q=a%26b&cursor=c2');
want('length cap', new URLSearchParams(artifactSearchPath('x'.repeat(300)).split('?')[1]).get('q')?.length, MAX_QUERY);

want('highlight', highlightParts('Encre tables scroll', 'encre SCROLL'), [
	{ text: 'Encre', hit: true },
	{ text: ' tables ', hit: false },
	{ text: 'scroll', hit: true }
]);
want('regexp chars are literal', highlightParts('a (b) c', '(b)'), [
	{ text: 'a ', hit: false },
	{ text: '(b)', hit: true },
	{ text: ' c', hit: false }
]);
want('longest word wins', highlightParts('encre', 'enc encre'), [{ text: 'encre', hit: true }]);
want('no query', highlightParts('plain', ' '), [{ text: 'plain', hit: false }]);
want('no text', highlightParts('', 'x'), []);

want('input is editable', isEditableTarget({ tagName: 'INPUT' } as unknown as EventTarget), true);
want('button is not', isEditableTarget({ tagName: 'BUTTON' } as unknown as EventTarget), false);
want('contenteditable', isEditableTarget({ tagName: 'DIV', isContentEditable: true } as unknown as EventTarget), true);
want('window is not', isEditableTarget(null), false);

const gate = latestOnly();
const slow = gate.next();
const fast = gate.next();
want('stale answer dropped', [gate.isCurrent(slow), gate.isCurrent(fast)], [false, true]);

console.log('artifact search helpers ok');
