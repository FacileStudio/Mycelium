import type { CliCommand } from '@facile/muse';

export const CLI_RELEASES = 'https://github.com/FacileStudio/Mycelium/releases';

export const CLI_COMMANDS: CliCommand[] = [
	{ label: 'curl', command: 'curl -fsSL https://raw.githubusercontent.com/FacileStudio/Mycelium/main/install.sh | bash' },
	{ label: 'facile', command: 'facile install mycelium' },
	{ label: 'Homebrew', command: 'brew install --cask FacileStudio/tap/mycelium' }
];
