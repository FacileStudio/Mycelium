import type { Row } from '../row';

/** The sign-in page, the app shell and the public landing page. */
export const shell = {
	'login.title': ['Sign in — Mycelium', 'Connexion — Mycelium', 'Iniciar sesión — Mycelium', 'Anmelden — Mycelium'],
	'login.tagline': [
		'One brain. All your agents.',
		'Un cerveau. Tous vos agents.',
		'Un cerebro. Todos tus agentes.',
		'Ein Gehirn. Alle deine Agenten.'
	],
	'login.failed': [
		'Could not sign in with Facile. Try again, or start a new sign-in.',
		'Connexion impossible avec Facile. Réessayez, ou commencez une nouvelle connexion.',
		'No se pudo iniciar sesión con Facile. Inténtalo de nuevo o empieza un nuevo inicio de sesión.',
		'Anmeldung mit Facile nicht möglich. Versuche es erneut oder starte eine neue Anmeldung.'
	],
	'callback.signingIn': ['Signing in…', 'Connexion…', 'Iniciando sesión…', 'Anmeldung läuft…'],
	'nav.memory': ['Memory', 'Mémoire', 'Memoria', 'Gedächtnis'],
	'nav.instructions': ['Instructions', 'Instructions', 'Instrucciones', 'Anweisungen'],
	'nav.automation': ['Automation', 'Automatisation', 'Automatización', 'Automatisierung'],
	'nav.artifacts': ['Artifacts', 'Artefacts', 'Artefactos', 'Artefakte'],
	'nav.machines': ['Machines', 'Machines', 'Máquinas', 'Maschinen'],
	'nav.sessions': ['Sessions', 'Sessions', 'Sesiones', 'Sitzungen'],
	'nav.spaces': ['Spaces', 'Espaces', 'Espacios', 'Bereiche'],
	'nav.personal': ['Personal', 'Personnel', 'Personal', 'Persönlich'],
	'nav.user': ['user', 'utilisateur', 'usuario', 'Benutzer'],
	'landing.title': [
		'Mycelium · One brain for your coding agents',
		'Mycelium · Un cerveau pour vos agents de code',
		'Mycelium · Un cerebro para tus agentes de código',
		'Mycelium · Ein Gehirn für deine Coding-Agenten'
	],
	'landing.description': [
		'One store of memory, rules and skills for every AI coding agent, synced across all your machines.',
		'Une seule réserve de mémoire, de règles et de compétences pour chaque agent IA de code, synchronisée sur toutes vos machines.',
		'Un único almacén de memoria, reglas y habilidades para cada agente de IA de código, sincronizado en todas tus máquinas.',
		'Ein Speicher für Gedächtnis, Regeln und Skills aller KI-Coding-Agenten, synchronisiert auf allen deinen Maschinen.'
	],
	'landing.nav.features': ['Features', 'Fonctionnalités', 'Funciones', 'Funktionen'],
	'landing.nav.downloads': ['Downloads', 'Téléchargements', 'Descargas', 'Downloads'],
	'landing.hero.title': [
		'One brain for every coding agent.',
		'Un cerveau pour tous vos agents de code.',
		'Un cerebro para todos tus agentes de código.',
		'Ein Gehirn für alle deine Coding-Agenten.'
	],
	'landing.hero.line': [
		'Memory, rules, skills and flows, shared by every agent on every machine.',
		'Mémoire, règles, skills et flows, partagés par chaque agent sur chaque machine.',
		'Memoria, reglas, skills y flows, compartidos por cada agente en cada máquina.',
		'Gedächtnis, Regeln, Skills und Flows, geteilt von jedem Agenten auf jeder Maschine.'
	],
	'landing.hero.cli': ['Install the CLI', 'Installer la CLI', 'Instalar la CLI', 'CLI installieren'],
	'landing.demo.label': [
		'An agent searches the shared memory and answers from a page another agent wrote.',
		'Un agent cherche dans la mémoire partagée et répond avec une page écrite par un autre agent.',
		'Un agente busca en la memoria compartida y responde con una página que escribió otro agente.',
		'Ein Agent durchsucht das gemeinsame Gedächtnis und antwortet mit einer Seite, die ein anderer Agent schrieb.'
	],
	'landing.demo.prompt': [
		'Why did the nightly backup stop?',
		'Pourquoi la sauvegarde de nuit s’est arrêtée ?',
		'¿Por qué se detuvo la copia de seguridad nocturna?',
		'Warum ist das nächtliche Backup ausgefallen?'
	],
	'landing.demo.answer': [
		'The backup user lost the docker group on 2026-08-12, so every dump failed. The fix is in conventions/backups.md.',
		'L’utilisateur de sauvegarde a perdu le groupe docker le 2026-08-12, donc chaque dump échoue. Le correctif est dans conventions/backups.md.',
		'El usuario de copias perdió el grupo docker el 2026-08-12, así que cada volcado falla. La solución está en conventions/backups.md.',
		'Der Backup-Benutzer hat am 2026-08-12 die Gruppe docker verloren, deshalb scheitert jeder Dump. Die Lösung steht in conventions/backups.md.'
	],
	'landing.features.title': ['How it works', 'Comment ça marche', 'Cómo funciona', 'So funktioniert es'],
	'landing.feature.memory.title': ['Shared memory', 'Mémoire partagée', 'Memoria compartida', 'Gemeinsames Gedächtnis'],
	'landing.feature.memory.body': [
		'A markdown wiki your agents search before they start and update when they finish.',
		'Un wiki markdown que vos agents consultent avant de commencer et complètent à la fin.',
		'Una wiki en markdown que tus agentes consultan antes de empezar y actualizan al terminar.',
		'Ein Markdown-Wiki, das deine Agenten vor dem Start durchsuchen und am Ende ergänzen.'
	],
	'landing.feature.rules.title': ['Rules and skills', 'Règles et skills', 'Reglas y skills', 'Regeln und Skills'],
	'landing.feature.rules.body': [
		'Write them once. Mycelium generates CLAUDE.md, AGENTS.md, GEMINI.md and the Cursor rules.',
		'Écrivez-les une fois. Mycelium génère CLAUDE.md, AGENTS.md, GEMINI.md et les règles Cursor.',
		'Escríbelas una vez. Mycelium genera CLAUDE.md, AGENTS.md, GEMINI.md y las reglas de Cursor.',
		'Einmal schreiben. Mycelium erzeugt CLAUDE.md, AGENTS.md, GEMINI.md und die Cursor-Regeln.'
	],
	'landing.feature.flows.title': ['Flows', 'Flows', 'Flows', 'Flows'],
	'landing.feature.flows.body': [
		'Shell procedures that already worked, replayed in order, with every run recorded.',
		'Des procédures shell qui ont déjà marché, rejouées dans l’ordre, chaque exécution enregistrée.',
		'Procedimientos de shell que ya funcionaron, repetidos en orden, con cada ejecución registrada.',
		'Shell-Abläufe, die schon funktioniert haben, der Reihe nach wiederholt, jeder Lauf protokolliert.'
	],
	'landing.feature.sync.title': ['Sync and spaces', 'Synchro et espaces', 'Sincronización y espacios', 'Sync und Bereiche'],
	'landing.feature.sync.body': [
		'Every machine syncs in the background. Personal, team and client spaces stay apart.',
		'Chaque machine se synchronise en arrière-plan. Les espaces perso, équipe et client restent séparés.',
		'Cada máquina se sincroniza en segundo plano. Los espacios personal, de equipo y de cliente quedan separados.',
		'Jede Maschine synchronisiert im Hintergrund. Private, Team- und Kundenbereiche bleiben getrennt.'
	],
	'landing.cli.title': ['From the terminal', 'Depuis le terminal', 'Desde la terminal', 'Aus dem Terminal'],
	'landing.cli.flow': [
		'Each step reports its exit code and time. The run is saved as JSON.',
		'Chaque étape donne son code de sortie et sa durée. L’exécution est enregistrée en JSON.',
		'Cada paso muestra su código de salida y su duración. La ejecución se guarda en JSON.',
		'Jeder Schritt meldet Exit-Code und Dauer. Der Lauf wird als JSON gespeichert.'
	],
	'landing.cta.title': [
		'Give your agents a memory.',
		'Donnez une mémoire à vos agents.',
		'Dales memoria a tus agentes.',
		'Gib deinen Agenten ein Gedächtnis.'
	],
	'landing.code': ['View the code', 'Voir le code', 'Ver el código', 'Code ansehen'],
	'downloads.title': ['Download Mycelium', 'Télécharger Mycelium', 'Descargar Mycelium', 'Mycelium herunterladen'],
	'downloads.line': [
		'Install the CLI on macOS or Linux, then run mycelium login.',
		'Installez la CLI sur macOS ou Linux, puis lancez mycelium login.',
		'Instala la CLI en macOS o Linux y ejecuta mycelium login.',
		'Installiere die CLI auf macOS oder Linux und führe dann mycelium login aus.'
	],
	'landing.pricing': ['Pricing', 'Tarifs', 'Precios', 'Preise'],
	'landing.footerTagline': [
		'Shared memory for your AI agents.',
		'La mémoire partagée de vos agents IA.',
		'La memoria compartida de tus agentes de IA.',
		'Das gemeinsame Gedächtnis deiner KI-Agenten.'
	]
} as const satisfies Record<string, Row>;
