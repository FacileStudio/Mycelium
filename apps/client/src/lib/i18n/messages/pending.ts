import type { Row } from '../row';

/** The confirm page for what agents ask over MCP. */
export const pending = {
	'pending.title': ['Waiting for you', 'En attente de vous', 'Esperando tu respuesta', 'Wartet auf dich'],
	'pending.intro': [
		'An agent asked to delete what is listed here. Nothing happens until you confirm.',
		'Un agent a demandé à supprimer ce qui figure ici. Rien ne se passe tant que vous ne confirmez pas.',
		'Un agente pidió borrar lo que aparece aquí. No pasa nada hasta que confirmes.',
		'Ein Agent möchte löschen, was hier steht. Nichts passiert, bis du bestätigst.'
	],
	'pending.none': ['Nothing waiting', 'Rien en attente', 'Nada pendiente', 'Nichts offen'],
	'pending.noneBody': [
		'Requests from your agents show up here when they need your confirmation.',
		'Les demandes de vos agents apparaissent ici quand elles ont besoin de votre confirmation.',
		'Las solicitudes de tus agentes aparecen aquí cuando necesitan tu confirmación.',
		'Anfragen deiner Agenten erscheinen hier, wenn sie deine Bestätigung brauchen.'
	],
	'pending.loadFailed': [
		'Could not load the waiting actions.',
		'Impossible de charger les actions en attente.',
		'No se pudieron cargar las acciones pendientes.',
		'Offene Aktionen konnten nicht geladen werden.'
	],
	'pending.deletePage': ['Delete a page', 'Supprimer une page', 'Borrar una página', 'Seite löschen'],
	'pending.deleteArtifact': ['Delete an artifact', 'Supprimer un artefact', 'Borrar un artefacto', 'Artefakt löschen'],
	'pending.confirmArtifact': ['Delete the artifact', 'Supprimer l’artefact', 'Borrar el artefacto', 'Artefakt löschen'],
	'pending.artifactDeleted': ['Artifact deleted.', 'Artefact supprimé.', 'Artefacto borrado.', 'Artefakt gelöscht.'],
	'pending.askedBy': ['Asked by {client}', 'Demandé par {client}', 'Pedido por {client}', 'Angefragt von {client}'],
	'pending.reason': ['Reason: {log}', 'Motif : {log}', 'Motivo: {log}', 'Grund: {log}'],
	'pending.changed': [
		'The page changed since the agent asked. Dismiss this and let the agent ask again.',
		'La page a changé depuis la demande de l’agent. Ignorez celle-ci et laissez l’agent redemander.',
		'La página cambió desde que el agente lo pidió. Descarta esto y deja que el agente vuelva a pedirlo.',
		'Die Seite hat sich seit der Anfrage geändert. Verwirf sie und lass den Agenten neu fragen.'
	],
	'pending.missing': [
		'This page is gone, or no longer yours.',
		'Cette page n’existe plus, ou n’est plus à vous.',
		'Esta página ya no existe o ya no es tuya.',
		'Diese Seite gibt es nicht mehr, oder sie gehört dir nicht mehr.'
	],
	'pending.dismiss': ['Dismiss', 'Ignorer', 'Descartar', 'Verwerfen'],
	'pending.confirmDelete': ['Delete the page', 'Supprimer la page', 'Borrar la página', 'Seite löschen'],
	'pending.deleted': ['Page deleted.', 'Page supprimée.', 'Página borrada.', 'Seite gelöscht.'],
	'pending.dismissed': ['Request dismissed.', 'Demande ignorée.', 'Solicitud descartada.', 'Anfrage verworfen.'],
	'pending.failed': [
		'That did not work. Try again.',
		'Cela n’a pas fonctionné. Réessayez.',
		'No funcionó. Inténtalo de nuevo.',
		'Das hat nicht geklappt. Versuch es noch einmal.'
	]
} as const satisfies Record<string, Row>;
