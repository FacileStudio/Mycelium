import type { Row } from '../row';

/** Strings shared by every page: buttons, states and relative times. */
export const common = {
	'common.cancel': ['Cancel', 'Annuler', 'Cancelar', 'Abbrechen'],
	'common.save': ['Save', 'Enregistrer', 'Guardar', 'Speichern'],
	'common.saving': ['Saving…', 'Enregistrement…', 'Guardando…', 'Wird gespeichert…'],
	'common.saved': ['Saved.', 'Enregistré.', 'Guardado.', 'Gespeichert.'],
	'common.saveFailed': ['Save failed', 'Échec de l’enregistrement', 'Error al guardar', 'Speichern fehlgeschlagen'],
	'common.create': ['Create', 'Créer', 'Crear', 'Erstellen'],
	'common.creating': ['Creating…', 'Création…', 'Creando…', 'Wird erstellt…'],
	'common.delete': ['Delete', 'Supprimer', 'Eliminar', 'Löschen'],
	'common.deleteFailed': ['Delete failed', 'Échec de la suppression', 'Error al eliminar', 'Löschen fehlgeschlagen'],
	'common.keep': ['Keep it', 'Conserver', 'Conservar', 'Behalten'],
	'common.edit': ['Edit', 'Modifier', 'Editar', 'Bearbeiten'],
	'common.close': ['Close', 'Fermer', 'Cerrar', 'Schließen'],
	'common.revoke': ['Revoke', 'Révoquer', 'Revocar', 'Widerrufen'],
	'common.loading': ['Loading', 'Chargement', 'Cargando', 'Wird geladen'],
	'common.loadingDots': ['Loading…', 'Chargement…', 'Cargando…', 'Wird geladen…'],
	'common.genericError': [
		'Something went wrong.',
		'Une erreur est survenue.',
		'Se ha producido un error.',
		'Es ist ein Fehler aufgetreten.'
	],
	'common.name': ['Name', 'Nom', 'Nombre', 'Name'],
	'common.copy': ['Copy', 'Copier', 'Copiar', 'Kopieren'],
	'common.copied': ['Copied', 'Copié', 'Copiado', 'Kopiert'],
	'common.root': ['root', 'racine', 'raíz', 'Wurzel'],
	'common.search': ['Search', 'Rechercher', 'Buscar', 'Suchen'],
	'common.searching': ['Searching…', 'Recherche…', 'Buscando…', 'Suche läuft…'],
	'time.justNow': ['just now', 'à l’instant', 'ahora mismo', 'gerade eben'],
	'time.minutesAgo': ['{count}m ago', 'il y a {count} min', 'hace {count} min', 'vor {count} Min.'],
	'time.hoursAgo': ['{count}h ago', 'il y a {count} h', 'hace {count} h', 'vor {count} Std.'],
	'time.daysAgo': ['{count}d ago', 'il y a {count} j', 'hace {count} d', 'vor {count} T.'],
	'time.now': ['now', 'maintenant', 'ahora', 'jetzt'],
	'time.unknown': ['unknown', 'inconnu', 'desconocido', 'unbekannt']
} as const satisfies Record<string, Row>;
