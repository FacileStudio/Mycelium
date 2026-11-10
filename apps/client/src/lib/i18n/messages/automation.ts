import type { Row } from '../row';

/** Flows, models and artifacts. */
export const automation = {
	'automation.title': ['Automation', 'Automatisation', 'Automatización', 'Automatisierung'],
	'automation.subtitle': [
		'Recorded procedures and the typed steps they run — synced across every machine.',
		'Procédures enregistrées et étapes typées qu’elles exécutent, synchronisées sur chaque machine.',
		'Procedimientos grabados y los pasos tipados que ejecutan, sincronizados en cada máquina.',
		'Aufgezeichnete Abläufe und die typisierten Schritte, die sie ausführen – auf allen Maschinen synchronisiert.'
	],
	'automation.sections': ['Automation sections', 'Sections d’automatisation', 'Secciones de automatización', 'Abschnitte der Automatisierung'],
	'flows.title': ['Flows', 'Flows', 'Flujos', 'Flows'],
	'flows.subtitle': [
		'Recorded shell procedures, synced across every machine. Trust and run history stay on the machine that runs them.',
		'Procédures shell enregistrées, synchronisées sur chaque machine. La confiance et l’historique d’exécution restent sur la machine qui les exécute.',
		'Procedimientos de shell grabados, sincronizados en cada máquina. La confianza y el historial de ejecución se quedan en la máquina que los ejecuta.',
		'Aufgezeichnete Shell-Abläufe, auf allen Maschinen synchronisiert. Vertrauen und Ausführungsverlauf bleiben auf der ausführenden Maschine.'
	],
	'flows.loadFailed': ['Could not load flows', 'Impossible de charger les flows', 'No se pudieron cargar los flujos', 'Flows konnten nicht geladen werden'],
	'flows.emptyTitle': ['No flows yet', 'Aucun flow pour l’instant', 'Aún no hay flujos', 'Noch keine Flows'],
	'flows.emptyText': [
		'Create one with `mycelium flow add <name>` — it appears here once it syncs.',
		'Créez-en un avec `mycelium flow add <name>` : il apparaît ici après la synchronisation.',
		'Crea uno con `mycelium flow add <name>`: aparece aquí cuando se sincroniza.',
		'Erstelle einen mit `mycelium flow add <name>` – er erscheint hier nach dem Sync.'
	],
	'flows.loadOne': ['Could not load flow “{name}”.', 'Impossible de charger le flow « {name} ».', 'No se pudo cargar el flujo «{name}».', 'Flow „{name}“ konnte nicht geladen werden.'],
	'flows.parseError': ['This flow does not parse', 'Ce flow ne peut pas être lu', 'Este flujo no se puede analizar', 'Dieser Flow lässt sich nicht lesen'],
	'flows.after': ['after {steps}', 'après {steps}', 'después de {steps}', 'nach {steps}'],
	'flows.needs': ['needs {names}', 'nécessite {names}', 'necesita {names}', 'benötigt {names}'],
	'flows.shell': ['shell', 'shell', 'shell', 'Shell'],
	'models.title': ['Models', 'Modèles', 'Modelos', 'Modelle'],
	'models.subtitle': [
		'Typed step extensions a flow can run with type:. Mycelium never executes these here — there is no bun on the server — so this is the source only.',
		'Extensions d’étapes typées qu’un flow peut exécuter avec type:. Mycelium ne les exécute jamais ici (il n’y a pas de bun sur le serveur) : c’est uniquement la source.',
		'Extensiones de pasos tipados que un flujo puede ejecutar con type:. Mycelium nunca las ejecuta aquí (no hay bun en el servidor): solo se muestra el código fuente.',
		'Typisierte Schritterweiterungen, die ein Flow mit type: ausführen kann. Mycelium führt sie hier nie aus (auf dem Server gibt es kein bun) – es ist nur die Quelle.'
	],
	'models.loadFailed': ['Could not load models', 'Impossible de charger les modèles', 'No se pudieron cargar los modelos', 'Modelle konnten nicht geladen werden'],
	'models.emptyTitle': ['No models yet', 'Aucun modèle pour l’instant', 'Aún no hay modelos', 'Noch keine Modelle'],
	'models.emptyText': [
		'Add a TypeScript file under extensions/models — it appears here once it syncs.',
		'Ajoutez un fichier TypeScript dans extensions/models : il apparaît ici après la synchronisation.',
		'Añade un archivo TypeScript en extensions/models: aparece aquí cuando se sincroniza.',
		'Füge unter extensions/models eine TypeScript-Datei hinzu – sie erscheint hier nach dem Sync.'
	],
	'models.loadOne': ['Could not load model “{name}”.', 'Impossible de charger le modèle « {name} ».', 'No se pudo cargar el modelo «{name}».', 'Modell „{name}“ konnte nicht geladen werden.'],
	'artifacts.title': ['Artifacts', 'Artefacts', 'Artefactos', 'Artefakte'],
	'artifacts.heading': ['Artifacts & Reports', 'Artefacts et rapports', 'Artefactos e informes', 'Artefakte und Berichte'],
	'artifacts.subtitle': [
		'Rendered markdown documents, reports, and architecture blueprints generated across your fleet.',
		'Documents markdown rendus, rapports et plans d’architecture générés sur l’ensemble de vos machines.',
		'Documentos markdown renderizados, informes y planos de arquitectura generados en toda tu flota.',
		'Gerenderte Markdown-Dokumente, Berichte und Architekturpläne, erzeugt auf all deinen Maschinen.'
	],
	'artifacts.loading': ['Loading artifacts…', 'Chargement des artefacts…', 'Cargando artefactos…', 'Artefakte werden geladen…'],
	'artifacts.loadFailed': ['Could not load artifacts', 'Impossible de charger les artefacts', 'No se pudieron cargar los artefactos', 'Artefakte konnten nicht geladen werden'],
	'artifacts.emptyTitle': ['No artifacts yet', 'Aucun artefact pour l’instant', 'Aún no hay artefactos', 'Noch keine Artefakte'],
	'artifacts.emptyText': [
		'Publish one with `mycelium artifact add <file>` or the `publish_artifact` tool.',
		'Publiez-en un avec `mycelium artifact add <file>` ou l’outil `publish_artifact`.',
		'Publica uno con `mycelium artifact add <file>` o la herramienta `publish_artifact`.',
		'Veröffentliche einen mit `mycelium artifact add <file>` oder dem Tool `publish_artifact`.'
	],
	'artifacts.filter': [
		'Filter artifacts by title, machine…',
		'Filtrer les artefacts par titre, machine…',
		'Filtrar artefactos por título, máquina…',
		'Artefakte nach Titel, Maschine filtern…'
	],
	'artifacts.noMatch': ['No artifacts match “{query}”.', 'Aucun artefact ne correspond à « {query} ».', 'Ningún artefacto coincide con «{query}».', 'Kein Artefakt passt zu „{query}“.'],
	'artifacts.updated': ['Updated {date}', 'Mis à jour le {date}', 'Actualizado el {date}', 'Aktualisiert am {date}'],
	'artifacts.created': ['Created {date}', 'Créé le {date}', 'Creado el {date}', 'Erstellt am {date}'],
	'artifacts.loadOne': ['Could not load artifact', 'Impossible de charger l’artefact', 'No se pudo cargar el artefacto', 'Artefakt konnte nicht geladen werden'],
	'artifacts.loadingOne': ['Loading artifact…', 'Chargement de l’artefact…', 'Cargando artefacto…', 'Artefakt wird geladen…'],
	'artifacts.deleted': ['Deleted artifact “{title}”.', 'Artefact « {title} » supprimé.', 'Artefacto «{title}» eliminado.', 'Artefakt „{title}“ gelöscht.'],
	'artifacts.preview': ['Preview', 'Aperçu', 'Vista previa', 'Vorschau'],
	'artifacts.source': ['Source', 'Source', 'Código fuente', 'Quelltext'],
	'artifacts.degradedTitle': ['Shown as best it can be', 'Affiché au mieux', 'Mostrado lo mejor posible', 'So gut wie möglich angezeigt'],
	'artifacts.degraded': ['This artifact does not follow the Facile document format ({reason}). Saving it again from Bulle or the CLI would be refused until it is fixed.', 'Cet artefact ne respecte pas le format de document Facile ({reason}). Le réenregistrer depuis Bulle ou la CLI sera refusé tant qu’il n’est pas corrigé.', 'Este artefacto no sigue el formato de documento Facile ({reason}). Volver a guardarlo desde Bulle o la CLI se rechazará hasta que se corrija.', 'Dieses Artefakt entspricht nicht dem Facile-Dokumentformat ({reason}). Ein erneutes Speichern aus Bulle oder der CLI wird abgelehnt, bis es korrigiert ist.'],
	'artifacts.codeCopied': ['Code copied to clipboard', 'Code copié dans le presse-papiers', 'Código copiado al portapapeles', 'Code in die Zwischenablage kopiert'],
	'artifacts.deleteTitle': ['Delete “{title}”?', 'Supprimer « {title} » ?', '¿Eliminar «{title}»?', '„{title}“ löschen?'],
	'artifacts.deleteText': [
		'The artifact will be permanently deleted from this machine and from sync.',
		'L’artefact sera supprimé définitivement de cette machine et de la synchronisation.',
		'El artefacto se eliminará de forma permanente de esta máquina y de la sincronización.',
		'Das Artefakt wird dauerhaft von dieser Maschine und aus dem Sync gelöscht.'
	],
	'artifacts.unknownMachine': ['unknown', 'inconnue', 'desconocida', 'unbekannt'],
} as const satisfies Record<string, Row>;
