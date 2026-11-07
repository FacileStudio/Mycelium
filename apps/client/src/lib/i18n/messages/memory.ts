import type { Row } from '../row';

/** Memory pages, search, the semantic index and the shared document editor. */
export const memory = {
	'memory.title': ['Memory', 'Mémoire', 'Memoria', 'Gedächtnis'],
	'memory.subtitle': [
		'Browse, search and curate the wiki your agents read from and write back to.',
		'Parcourez, recherchez et organisez le wiki que vos agents lisent et complètent.',
		'Explora, busca y organiza la wiki que tus agentes leen y completan.',
		'Durchsuche und pflege das Wiki, das deine Agenten lesen und ergänzen.'
	],
	'memory.newPage': ['New page', 'Nouvelle page', 'Nueva página', 'Neue Seite'],
	'memory.loadFailed': [
		'Could not load the memory tree.',
		'Impossible de charger l’arborescence de la mémoire.',
		'No se pudo cargar el árbol de memoria.',
		'Der Gedächtnisbaum konnte nicht geladen werden.'
	],
	'memory.emptyTitle': ['No memories yet', 'Aucune mémoire pour l’instant', 'Aún no hay memorias', 'Noch keine Einträge'],
	'memory.emptyText': [
		'Create a page, or let your agents fill this in as they learn.',
		'Créez une page, ou laissez vos agents la remplir au fil de leur apprentissage.',
		'Crea una página o deja que tus agentes la llenen a medida que aprenden.',
		'Erstelle eine Seite oder lass deine Agenten sie beim Lernen füllen.'
	],
	'memory.createTitle': ['New memory page', 'Nouvelle page de mémoire', 'Nueva página de memoria', 'Neue Gedächtnisseite'],
	'memory.folder': ['Folder', 'Dossier', 'Carpeta', 'Ordner'],
	'memory.folderHelp': [
		'Leave it on root for a page that fits nowhere else.',
		'Laissez la racine pour une page qui ne va nulle part ailleurs.',
		'Déjalo en la raíz para una página que no encaja en otro sitio.',
		'Lass es auf der Wurzel für eine Seite, die sonst nirgends passt.'
	],
	'memory.nameHelp': [
		'The .md extension is added if you leave it off.',
		'L’extension .md est ajoutée si vous l’omettez.',
		'La extensión .md se añade si la omites.',
		'Die Endung .md wird ergänzt, wenn du sie weglässt.'
	],
	'memory.createPage': ['Create page', 'Créer la page', 'Crear página', 'Seite erstellen'],
	'memory.createFailed': [
		'Could not create the page',
		'Impossible de créer la page',
		'No se pudo crear la página',
		'Die Seite konnte nicht erstellt werden'
	],
	'memory.created': ['Created {path}.', '{path} créé.', '{path} creado.', '{path} erstellt.'],
	'memory.deleted': ['Deleted {path}.', '{path} supprimé.', '{path} eliminado.', '{path} gelöscht.'],
	'memory.deleteTitle': ['Delete “{title}”?', 'Supprimer « {title} » ?', '¿Eliminar «{title}»?', '„{title}“ löschen?'],
	'memory.deleteText': [
		'The page is removed from every machine on the next sync, and nothing keeps a copy. It cannot be undone.',
		'La page est retirée de chaque machine à la prochaine synchronisation, et aucune copie n’est conservée. Action irréversible.',
		'La página se elimina de todas las máquinas en la próxima sincronización y no se conserva ninguna copia. No se puede deshacer.',
		'Die Seite wird beim nächsten Sync von jeder Maschine entfernt, und es bleibt keine Kopie. Das lässt sich nicht rückgängig machen.'
	],
	'memory.loadFile': ['Could not load {path}.', 'Impossible de charger {path}.', 'No se pudo cargar {path}.', '{path} konnte nicht geladen werden.'],
	'search.placeholder': ['Search memory…', 'Rechercher dans la mémoire…', 'Buscar en la memoria…', 'Im Gedächtnis suchen…'],
	'search.label': ['Search memory', 'Rechercher dans la mémoire', 'Buscar en la memoria', 'Im Gedächtnis suchen'],
	'search.failed': ['Search failed.', 'La recherche a échoué.', 'La búsqueda ha fallado.', 'Die Suche ist fehlgeschlagen.'],
	'search.degraded': [
		'Matched words only — the semantic half is unavailable.',
		'Mots exacts uniquement : la partie sémantique est indisponible.',
		'Solo coincidencia de palabras: la parte semántica no está disponible.',
		'Nur Wortübereinstimmungen – die semantische Hälfte ist nicht verfügbar.'
	],
	'search.none': ['Nothing matched “{query}”.', 'Rien ne correspond à « {query} ».', 'Nada coincide con «{query}».', 'Nichts passt zu „{query}“.'],
	'stats.pages': ['Pages', 'Pages', 'Páginas', 'Seiten'],
	'stats.folders': ['Folders', 'Dossiers', 'Carpetas', 'Ordner'],
	'stats.size': ['Size', 'Taille', 'Tamaño', 'Größe'],
	'stats.lastWritten': ['Last written', 'Dernière écriture', 'Última escritura', 'Zuletzt geschrieben'],
	'stats.pagesTouched': ['Pages touched', 'Pages modifiées', 'Páginas modificadas', 'Geänderte Seiten'],
	'stats.perDay': [
		'Pages written per day · last {days} days',
		'Pages écrites par jour · {days} derniers jours',
		'Páginas escritas por día · últimos {days} días',
		'Geschriebene Seiten pro Tag · letzte {days} Tage'
	],
	'stats.perFolder': ['Pages per folder', 'Pages par dossier', 'Páginas por carpeta', 'Seiten pro Ordner'],
	'index.title': ['Semantic index', 'Index sémantique', 'Índice semántico', 'Semantischer Index'],
	'index.indexing': ['Indexing', 'Indexation', 'Indexando', 'Indizierung'],
	'index.ready': ['Ready', 'Prêt', 'Listo', 'Bereit'],
	'index.off': ['Off', 'Désactivé', 'Desactivado', 'Aus'],
	'index.offText': [
		'Semantic search is off. Search matches words literally until a model is configured.',
		'La recherche sémantique est désactivée. La recherche compare les mots tels quels tant qu’aucun modèle n’est configuré.',
		'La búsqueda semántica está desactivada. La búsqueda compara palabras literalmente hasta que se configure un modelo.',
		'Die semantische Suche ist aus. Die Suche vergleicht Wörter wörtlich, bis ein Modell konfiguriert ist.'
	],
	'index.progress': ['{done} / {total} chunks', '{done} / {total} fragments', '{done} / {total} fragmentos', '{done} / {total} Abschnitte'],
	'index.embedded': ['Chunks embedded', 'Fragments vectorisés', 'Fragmentos vectorizados', 'Abschnitte eingebettet'],
	'index.rate': ['{rate} chunks/sec', '{rate} fragments/s', '{rate} fragmentos/s', '{rate} Abschnitte/s'],
	'index.measuring': ['measuring rate…', 'mesure du débit…', 'midiendo la velocidad…', 'Rate wird gemessen…'],
	'index.almostDone': ['almost done', 'presque terminé', 'casi terminado', 'fast fertig'],
	'index.lessThanMinute': ['less than a minute left', 'moins d’une minute restante', 'queda menos de un minuto', 'weniger als eine Minute übrig'],
	'index.minutesLeft': ['about {count} min left', 'environ {count} min restantes', 'quedan unos {count} min', 'noch etwa {count} Min.'],
	'index.hoursLeft': ['about {hours} h left', 'environ {hours} h restantes', 'quedan unas {hours} h', 'noch etwa {hours} Std.'],
	'index.hoursMinutesLeft': [
		'about {hours} h {minutes} min left',
		'environ {hours} h {minutes} min restantes',
		'quedan unas {hours} h {minutes} min',
		'noch etwa {hours} Std. {minutes} Min.'
	],
	'index.queued': ['{count} queued', '{count} en attente', '{count} en cola', '{count} in der Warteschlange'],
	'index.doneText': ['{count} chunks embedded', '{count} fragments vectorisés', '{count} fragmentos vectorizados', '{count} Abschnitte eingebettet'],
	'index.doneQueued': [
		'{count} chunks embedded, {queued} queued',
		'{count} fragments vectorisés, {queued} en attente',
		'{count} fragmentos vectorizados, {queued} en cola',
		'{count} Abschnitte eingebettet, {queued} in der Warteschlange'
	],
	'index.dims': ['{count} dims', '{count} dim.', '{count} dim.', '{count} Dim.'],
	'index.store': ['{name} store', 'stockage {name}', 'almacén {name}', 'Speicher {name}'],
	'index.readFailed': [
		'Could not read the index state.',
		'Impossible de lire l’état de l’index.',
		'No se pudo leer el estado del índice.',
		'Der Indexstatus konnte nicht gelesen werden.'
	],
	'rules.title': ['Rules', 'Règles', 'Reglas', 'Regeln'],
	'rules.subtitle': [
		'Modular instructions, concatenated into every agent config in filename order.',
		'Instructions modulaires, concaténées dans la config de chaque agent par ordre de nom de fichier.',
		'Instrucciones modulares, concatenadas en la configuración de cada agente por orden de nombre de archivo.',
		'Modulare Anweisungen, in Dateinamen-Reihenfolge in jede Agentenkonfiguration eingefügt.'
	],
	'rules.new': ['New rule', 'Nouvelle règle', 'Nueva regla', 'Neue Regel'],
	'rules.loadFailed': ['Could not load rules', 'Impossible de charger les règles', 'No se pudieron cargar las reglas', 'Regeln konnten nicht geladen werden'],
	'rules.emptyTitle': ['No rules yet', 'Aucune règle pour l’instant', 'Aún no hay reglas', 'Noch keine Regeln'],
	'rules.emptyText': [
		'Create one to shape how every agent behaves.',
		'Créez-en une pour guider le comportement de chaque agent.',
		'Crea una para definir cómo se comporta cada agente.',
		'Erstelle eine, um das Verhalten jedes Agenten zu steuern.'
	],
	'rules.nameHelp': [
		'Prefix it to control the order agents read it in — 30-testing, 40-git.',
		'Préfixez-la pour contrôler l’ordre de lecture : 30-testing, 40-git.',
		'Ponle un prefijo para controlar el orden de lectura: 30-testing, 40-git.',
		'Setze ein Präfix, um die Lesereihenfolge zu steuern: 30-testing, 40-git.'
	],
	'rules.create': ['Create rule', 'Créer la règle', 'Crear regla', 'Regel erstellen'],
	'rules.createFailed': ['Could not create the rule', 'Impossible de créer la règle', 'No se pudo crear la regla', 'Die Regel konnte nicht erstellt werden'],
	'rules.created': ['Rule “{name}” created.', 'Règle « {name} » créée.', 'Regla «{name}» creada.', 'Regel „{name}“ erstellt.'],
	'rules.deleted': ['Rule “{name}” deleted.', 'Règle « {name} » supprimée.', 'Regla «{name}» eliminada.', 'Regel „{name}“ gelöscht.'],
	'rules.loadOne': ['Could not load rule “{name}”.', 'Impossible de charger la règle « {name} ».', 'No se pudo cargar la regla «{name}».', 'Regel „{name}“ konnte nicht geladen werden.'],
	'rules.deleteText': [
		'Every agent config is rebuilt without this rule on the next sync. It cannot be undone.',
		'La config de chaque agent est reconstruite sans cette règle à la prochaine synchronisation. Action irréversible.',
		'La configuración de cada agente se reconstruye sin esta regla en la próxima sincronización. No se puede deshacer.',
		'Jede Agentenkonfiguration wird beim nächsten Sync ohne diese Regel neu aufgebaut. Das lässt sich nicht rückgängig machen.'
	],
	'skills.title': ['Skills', 'Compétences', 'Habilidades', 'Skills'],
	'skills.subtitle': [
		'Agent-agnostic capabilities, installed into each agent’s own skill format.',
		'Capacités indépendantes de l’agent, installées dans le format de skill propre à chaque agent.',
		'Capacidades independientes del agente, instaladas en el formato de habilidad propio de cada agente.',
		'Agent-unabhängige Fähigkeiten, installiert im jeweiligen Skill-Format des Agenten.'
	],
	'skills.new': ['New skill', 'Nouveau skill', 'Nueva habilidad', 'Neuer Skill'],
	'skills.loadFailed': ['Could not load skills', 'Impossible de charger les skills', 'No se pudieron cargar las habilidades', 'Skills konnten nicht geladen werden'],
	'skills.emptyTitle': ['No skills yet', 'Aucun skill pour l’instant', 'Aún no hay habilidades', 'Noch keine Skills'],
	'skills.emptyText': [
		'Add one to teach every agent the same trick.',
		'Ajoutez-en un pour apprendre le même tour à chaque agent.',
		'Añade una para enseñar el mismo truco a todos los agentes.',
		'Füge einen hinzu, um jedem Agenten denselben Trick beizubringen.'
	],
	'skills.nameHelp': [
		'Becomes the slash command agents invoke it with.',
		'Devient la commande slash avec laquelle les agents l’invoquent.',
		'Se convierte en el comando slash con el que los agentes lo invocan.',
		'Wird zum Slash-Befehl, mit dem Agenten ihn aufrufen.'
	],
	'skills.create': ['Create skill', 'Créer le skill', 'Crear habilidad', 'Skill erstellen'],
	'skills.createFailed': ['Could not create the skill', 'Impossible de créer le skill', 'No se pudo crear la habilidad', 'Der Skill konnte nicht erstellt werden'],
	'skills.created': ['Skill “{name}” created.', 'Skill « {name} » créé.', 'Habilidad «{name}» creada.', 'Skill „{name}“ erstellt.'],
	'skills.deleted': ['Skill “{name}” deleted.', 'Skill « {name} » supprimé.', 'Habilidad «{name}» eliminada.', 'Skill „{name}“ gelöscht.'],
	'skills.loadOne': ['Could not load skill “{name}”.', 'Impossible de charger le skill « {name} ».', 'No se pudo cargar la habilidad «{name}».', 'Skill „{name}“ konnte nicht geladen werden.'],
	'skills.deleteText': [
		'The skill disappears from every agent on the next sync. It cannot be undone.',
		'Le skill disparaît de chaque agent à la prochaine synchronisation. Action irréversible.',
		'La habilidad desaparece de todos los agentes en la próxima sincronización. No se puede deshacer.',
		'Der Skill verschwindet beim nächsten Sync bei jedem Agenten. Das lässt sich nicht rückgängig machen.'
	],
	'instructions.title': ['Instructions', 'Instructions', 'Instrucciones', 'Anweisungen'],
	'instructions.subtitle': [
		'What every agent reads before it starts working — installed into each agent’s own config format.',
		'Ce que chaque agent lit avant de se mettre au travail, installé dans le format de config propre à chaque agent.',
		'Lo que cada agente lee antes de empezar a trabajar, instalado en el formato de configuración propio de cada agente.',
		'Was jeder Agent liest, bevor er loslegt – installiert im jeweiligen Konfigurationsformat des Agenten.'
	],
	'instructions.sections': ['Instruction sections', 'Sections d’instructions', 'Secciones de instrucciones', 'Abschnitte der Anweisungen']
} as const satisfies Record<string, Row>;
