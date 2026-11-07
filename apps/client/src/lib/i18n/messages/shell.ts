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
		'Mycelium — One brain for your agents',
		'Mycelium — Un cerveau pour vos agents',
		'Mycelium — Un cerebro para tus agentes',
		'Mycelium — Ein Gehirn für deine Agenten'
	],
	'landing.description': [
		'Shared memory, rules and skills for your AI agents, synced across pi, Claude, Gemini, Codex and every machine.',
		'Mémoire, règles et compétences partagées pour vos agents IA, synchronisées entre pi, Claude, Gemini, Codex et toutes vos machines.',
		'Memoria, reglas y habilidades compartidas para tus agentes de IA, sincronizadas entre pi, Claude, Gemini, Codex y todas tus máquinas.',
		'Geteiltes Gedächtnis, Regeln und Skills für deine KI-Agenten, synchronisiert zwischen pi, Claude, Gemini, Codex und allen Maschinen.'
	],
	'landing.badge': [
		'Local-first · Multi-agent · Open source',
		'Local-first · Multi-agent · Open source',
		'Local-first · Multiagente · Código abierto',
		'Local-first · Multi-Agent · Open Source'
	],
	'landing.heroTitle': [
		'One brain.|All your agents.',
		'Un cerveau.|Tous vos agents.',
		'Un cerebro.|Todos tus agentes.',
		'Ein Gehirn.|Alle deine Agenten.'
	],
	'landing.heroText': [
		'Memory, rules and skills shared between pi, Claude, Gemini, Codex, Cursor and all your machines. One place, zero friction.',
		'Mémoire, règles et compétences partagées entre pi, Claude, Gemini, Codex, Cursor et toutes vos machines. Un seul endroit, zéro friction.',
		'Memoria, reglas y habilidades compartidas entre pi, Claude, Gemini, Codex, Cursor y todas tus máquinas. Un solo lugar, cero fricción.',
		'Gedächtnis, Regeln und Skills, geteilt zwischen pi, Claude, Gemini, Codex, Cursor und allen deinen Maschinen. Ein Ort, keine Reibung.'
	],
	'landing.start': ['Get started', 'Commencer', 'Empezar', 'Loslegen'],
	'landing.adaptersLabel': ['Adapters', 'Adaptateurs', 'Adaptadores', 'Adapter'],
	'landing.adaptersTitle': [
		'Write once.|Deploy everywhere.',
		'Écrivez une fois.|Déployez partout.',
		'Escribe una vez.|Despliega en todas partes.',
		'Einmal schreiben.|Überall ausrollen.'
	],
	'landing.adaptersText': [
		'Your rules and skills are written in markdown. Mycelium generates the configuration files for every agent.',
		'Vos règles et compétences sont écrites en markdown. Mycelium génère automatiquement les fichiers de configuration pour chaque agent.',
		'Tus reglas y habilidades se escriben en markdown. Mycelium genera los archivos de configuración de cada agente.',
		'Deine Regeln und Skills sind in Markdown geschrieben. Mycelium erzeugt die Konfigurationsdateien für jeden Agenten.'
	],
	'landing.layersTitle': [
		'Three layers.|Each one independent.',
		'Trois couches.|Chacune indépendante.',
		'Tres capas.|Cada una independiente.',
		'Drei Schichten.|Jede für sich.'
	],
	'landing.layersText': [
		'No dependencies. Each piece works alone and gets richer with the others.',
		'Pas de dépendance. Chaque pièce fonctionne seule et s’enrichit avec les autres.',
		'Sin dependencias. Cada pieza funciona sola y se enriquece con las demás.',
		'Keine Abhängigkeiten. Jedes Teil funktioniert allein und wird mit den anderen besser.'
	],
	'landing.layer.memory.title': ['Memory', 'Mémoire', 'Memoria', 'Gedächtnis'],
	'landing.layer.memory.body': [
		'A shared markdown wiki. Bugs, tools, projects, conventions: your agents learn from every session and share what they know.',
		'Wiki partagé en markdown. Bugs, outils, projets, conventions : vos agents apprennent de chaque session et partagent ce savoir entre eux.',
		'Una wiki compartida en markdown. Errores, herramientas, proyectos, convenciones: tus agentes aprenden de cada sesión y comparten lo que saben.',
		'Ein gemeinsames Markdown-Wiki. Bugs, Tools, Projekte, Konventionen: deine Agenten lernen aus jeder Sitzung und teilen ihr Wissen.'
	],
	'landing.layer.rules.title': ['Rules', 'Règles', 'Reglas', 'Regeln'],
	'landing.layer.rules.body': [
		'Modular rules for your agents. Code style, git conventions, engineering ladder: write them once, every agent follows them.',
		'Règles modulaires pour vos agents. Style de code, conventions git, engineering ladder : écrivez-les une fois, tous vos agents les suivent.',
		'Reglas modulares para tus agentes. Estilo de código, convenciones de git, engineering ladder: escríbelas una vez y todos tus agentes las siguen.',
		'Modulare Regeln für deine Agenten. Code-Stil, Git-Konventionen, Engineering Ladder: einmal schreiben, alle Agenten befolgen sie.'
	],
	'landing.layer.skills.title': ['Skills', 'Compétences', 'Habilidades', 'Skills'],
	'landing.layer.skills.body': [
		'Agent-agnostic skills with portable definitions. One skill, six agents. No vendor lock-in.',
		'Compétences agnostiques à l’agent, avec des définitions portables. Un skill, six agents. Pas de vendor lock-in.',
		'Habilidades independientes del agente con definiciones portables. Una habilidad, seis agentes. Sin dependencia de proveedor.',
		'Agent-unabhängige Skills mit portablen Definitionen. Ein Skill, sechs Agenten. Kein Vendor-Lock-in.'
	],
	'landing.cellsLabel': ['Cells', 'Cellules', 'Celdas', 'Zellen'],
	'landing.cellsTitle': ['Personal vs. team.', 'Perso vs. équipe.', 'Personal frente a equipo.', 'Privat oder Team.'],
	'landing.cellsText': [
		'A personal profile, a team profile, a client profile: each cell has its own memory, rules and skills. Stack them: personal rules win, memory adds up.',
		'Un profil personnel, un profil équipe, un profil client : chaque cellule a sa propre mémoire, ses règles et ses skills. Superposez-les : les règles perso gagnent, la mémoire s’additionne.',
		'Un perfil personal, un perfil de equipo, un perfil de cliente: cada celda tiene su propia memoria, reglas y habilidades. Apílalas: las reglas personales ganan, la memoria se suma.',
		'Ein persönliches Profil, ein Teamprofil, ein Kundenprofil: jede Zelle hat eigenes Gedächtnis, eigene Regeln und Skills. Schichte sie: persönliche Regeln gewinnen, das Gedächtnis addiert sich.'
	],
	'landing.cell.personal': [
		'Your memory, your preferences, your shortcuts',
		'Votre mémoire, vos préférences, vos raccourcis',
		'Tu memoria, tus preferencias, tus atajos',
		'Dein Gedächtnis, deine Vorlieben, deine Abkürzungen'
	],
	'landing.cell.team': [
		'Team conventions, tech stack, shared projects',
		'Conventions d’équipe, stack technique, projets partagés',
		'Convenciones del equipo, stack técnico, proyectos compartidos',
		'Team-Konventionen, Tech-Stack, gemeinsame Projekte'
	],
	'landing.cell.client': [
		'Project context, client rules',
		'Contexte spécifique au projet, règles du client',
		'Contexto del proyecto, reglas del cliente',
		'Projektkontext, Kundenregeln'
	],
	'landing.syncLabel': ['Sync', 'Synchro', 'Sincronización', 'Sync'],
	'landing.syncTitle': [
		'Built-in sync.|Same binary.',
		'Sync intégré.|Même binaire.',
		'Sincronización integrada.|El mismo binario.',
		'Integrierter Sync.|Dasselbe Binary.'
	],
	'landing.syncText': [
		'mycelium serve runs an HTTP sync server. Deploy it on your VPS and connect your machines with a token.',
		'mycelium serve lance un serveur de sync HTTP. Déployez-le sur votre VPS, connectez vos machines avec un token.',
		'mycelium serve lanza un servidor de sincronización HTTP. Despliégalo en tu VPS y conecta tus máquinas con un token.',
		'mycelium serve startet einen HTTP-Sync-Server. Betreibe ihn auf deinem VPS und verbinde deine Maschinen mit einem Token.'
	],
	'landing.syncOnServer': ['# on the server', '# sur le serveur', '# en el servidor', '# auf dem Server'],
	'landing.syncOnMachine': ['# on your machine', '# sur votre machine', '# en tu máquina', '# auf deiner Maschine'],
	'landing.syncDone': ['Synced 2 file(s).', '2 fichier(s) synchronisé(s).', '2 archivo(s) sincronizado(s).', '2 Datei(en) synchronisiert.'],
	'landing.ctaTitle': [
		'Open source. Local-first. Free.',
		'Open source. Local-first. Gratuit.',
		'Código abierto. Local-first. Gratis.',
		'Open Source. Local-first. Kostenlos.'
	],
	'landing.ctaText': [
		'One Go binary. Zero dependencies. Your data stays with you.',
		'Un binaire Go. Zéro dépendance. Vos données restent chez vous.',
		'Un binario de Go. Cero dependencias. Tus datos se quedan contigo.',
		'Ein Go-Binary. Keine Abhängigkeiten. Deine Daten bleiben bei dir.'
	],
	'landing.code': ['View the code', 'Voir le code', 'Ver el código', 'Code ansehen'],
	'landing.pricing': ['Pricing', 'Tarifs', 'Precios', 'Preise'],
	'landing.footerTagline': [
		'Shared memory for your AI agents.',
		'La mémoire partagée de vos agents IA.',
		'La memoria compartida de tus agentes de IA.',
		'Das gemeinsame Gedächtnis deiner KI-Agenten.'
	]
} as const satisfies Record<string, Row>;
