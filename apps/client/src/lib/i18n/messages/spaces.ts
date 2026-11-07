import type { Row } from '../row';

/** Spaces and their members. */
export const spaces = {
	'role.owner': ['owner', 'propriétaire', 'propietario', 'Eigentümer'],
	'role.admin': ['admin', 'admin', 'admin', 'Admin'],
	'role.member': ['member', 'membre', 'miembro', 'Mitglied'],
	'spaces.title': ['Spaces', 'Espaces', 'Espacios', 'Bereiche'],
	'spaces.subtitleRegistre': [
		'Shared memory trees. Members and roles are managed in Registre, for every Facile app at once.',
		'Arbres de mémoire partagés. Les membres et les rôles se gèrent dans Registre, pour toutes les applications Facile à la fois.',
		'Árboles de memoria compartidos. Los miembros y los roles se gestionan en Registre, para todas las aplicaciones de Facile a la vez.',
		'Gemeinsame Gedächtnisbäume. Mitglieder und Rollen werden in Registre verwaltet, für alle Facile-Apps gleichzeitig.'
	],
	'spaces.subtitleLocal': [
		'Shared memory trees. Invite teammates and scope memory, rules and skills per space.',
		'Arbres de mémoire partagés. Invitez vos collègues et limitez mémoire, règles et skills à chaque espace.',
		'Árboles de memoria compartidos. Invita a tus compañeros y limita memoria, reglas y habilidades por espacio.',
		'Gemeinsame Gedächtnisbäume. Lade Kollegen ein und begrenze Gedächtnis, Regeln und Skills pro Bereich.'
	],
	'spaces.manage': ['Manage in Registre', 'Gérer dans Registre', 'Gestionar en Registre', 'In Registre verwalten'],
	'spaces.new': ['New space', 'Nouvel espace', 'Nuevo espacio', 'Neuer Bereich'],
	'spaces.loadFailed': ['Could not load spaces.', 'Impossible de charger les espaces.', 'No se pudieron cargar los espacios.', 'Bereiche konnten nicht geladen werden.'],
	'spaces.emptyTitle': ['No spaces yet', 'Aucun espace pour l’instant', 'Aún no hay espacios', 'Noch keine Bereiche'],
	'spaces.emptyRegistre': [
		'Spaces are shared across the Facile suite. Create one in Registre and it appears here.',
		'Les espaces sont partagés dans toute la suite Facile. Créez-en un dans Registre et il apparaît ici.',
		'Los espacios se comparten en toda la suite Facile. Crea uno en Registre y aparecerá aquí.',
		'Bereiche werden in der gesamten Facile-Suite geteilt. Erstelle einen in Registre, und er erscheint hier.'
	],
	'spaces.emptyLocal': [
		'Create one to share a memory tree with the people you work with.',
		'Créez-en un pour partager un arbre de mémoire avec les personnes avec qui vous travaillez.',
		'Crea uno para compartir un árbol de memoria con las personas con las que trabajas.',
		'Erstelle einen, um einen Gedächtnisbaum mit deinen Kollegen zu teilen.'
	],
	'spaces.noDescription': ['No description', 'Aucune description', 'Sin descripción', 'Keine Beschreibung'],
	'spaces.description': ['Description', 'Description', 'Descripción', 'Beschreibung'],
	'spaces.descriptionHelp': [
		'Optional — what this memory tree is for.',
		'Facultatif : à quoi sert cet arbre de mémoire.',
		'Opcional: para qué sirve este árbol de memoria.',
		'Optional – wofür dieser Gedächtnisbaum gedacht ist.'
	],
	'spaces.descriptionPlaceholder': [
		'Team conventions, stack and shared projects.',
		'Conventions d’équipe, stack et projets partagés.',
		'Convenciones del equipo, stack y proyectos compartidos.',
		'Team-Konventionen, Stack und gemeinsame Projekte.'
	],
	'spaces.create': ['Create space', 'Créer l’espace', 'Crear espacio', 'Bereich erstellen'],
	'spaces.createFailed': ['Could not create the space', 'Impossible de créer l’espace', 'No se pudo crear el espacio', 'Der Bereich konnte nicht erstellt werden'],
	'spaces.created': ['Space “{name}” created.', 'Espace « {name} » créé.', 'Espacio «{name}» creado.', 'Bereich „{name}“ erstellt.'],
	'spaces.loadOne': ['Could not load this space.', 'Impossible de charger cet espace.', 'No se pudo cargar este espacio.', 'Dieser Bereich konnte nicht geladen werden.'],
	'spaces.missing': [
		'This space does not exist, or you left it.',
		'Cet espace n’existe pas, ou vous l’avez quitté.',
		'Este espacio no existe o lo has abandonado.',
		'Dieser Bereich existiert nicht, oder du hast ihn verlassen.'
	],
	'spaces.edit': ['Edit space', 'Modifier l’espace', 'Editar espacio', 'Bereich bearbeiten'],
	'spaces.saveChanges': ['Save changes', 'Enregistrer les modifications', 'Guardar cambios', 'Änderungen speichern'],
	'spaces.updated': ['Space updated.', 'Espace mis à jour.', 'Espacio actualizado.', 'Bereich aktualisiert.'],
	'spaces.updateFailed': ['Update failed', 'Échec de la mise à jour', 'Error al actualizar', 'Aktualisierung fehlgeschlagen'],
	'spaces.leave': ['Leave', 'Quitter', 'Abandonar', 'Verlassen'],
	'spaces.left': ['You left the space.', 'Vous avez quitté l’espace.', 'Has abandonado el espacio.', 'Du hast den Bereich verlassen.'],
	'spaces.leaveFailed': ['Could not leave this space.', 'Impossible de quitter cet espace.', 'No se pudo abandonar este espacio.', 'Der Bereich konnte nicht verlassen werden.'],
	'spaces.deleted': ['Space deleted.', 'Espace supprimé.', 'Espacio eliminado.', 'Bereich gelöscht.'],
	'spaces.deleteFailed': ['Could not delete this space.', 'Impossible de supprimer cet espace.', 'No se pudo eliminar este espacio.', 'Der Bereich konnte nicht gelöscht werden.'],
	'spaces.thisSpace': ['this space', 'cet espace', 'este espacio', 'diesen Bereich'],
	'spaces.leaveTitle': ['Leave {name}?', 'Quitter {name} ?', '¿Abandonar {name}?', '{name} verlassen?'],
	'spaces.leaveText': [
		'This space disappears from your machines on the next sync. An owner has to invite you back.',
		'Cet espace disparaît de vos machines à la prochaine synchronisation. Un propriétaire doit vous réinviter.',
		'Este espacio desaparece de tus máquinas en la próxima sincronización. Un propietario tiene que volver a invitarte.',
		'Dieser Bereich verschwindet beim nächsten Sync von deinen Maschinen. Ein Eigentümer muss dich erneut einladen.'
	],
	'spaces.leaveConfirm': ['Leave space', 'Quitter l’espace', 'Abandonar espacio', 'Bereich verlassen'],
	'spaces.stay': ['Stay', 'Rester', 'Quedarme', 'Bleiben'],
	'spaces.deleteTitle': ['Delete {name}?', 'Supprimer {name} ?', '¿Eliminar {name}?', '{name} löschen?'],
	'spaces.deleteText': [
		'Every member loses it, and the memory, rules and skills stored in it are gone. This cannot be undone.',
		'Tous les membres le perdent, et la mémoire, les règles et les skills qu’il contient disparaissent. Action irréversible.',
		'Todos los miembros lo pierden y la memoria, las reglas y las habilidades que contiene desaparecen. No se puede deshacer.',
		'Alle Mitglieder verlieren ihn, und Gedächtnis, Regeln und Skills darin sind weg. Das lässt sich nicht rückgängig machen.'
	],
	'spaces.deleteConfirm': ['Delete space', 'Supprimer l’espace', 'Eliminar espacio', 'Bereich löschen'],
	'members.title': ['Members', 'Membres', 'Miembros', 'Mitglieder'],
	'members.subtitle': [
		'Everyone here reads and writes the same memory, rules and skills.',
		'Toutes les personnes ici lisent et écrivent la même mémoire, les mêmes règles et les mêmes skills.',
		'Todas las personas aquí leen y escriben la misma memoria, reglas y habilidades.',
		'Alle hier lesen und schreiben dasselbe Gedächtnis, dieselben Regeln und Skills.'
	],
	'members.email': ['Email', 'E-mail', 'Correo electrónico', 'E-Mail'],
	'members.role': ['Role', 'Rôle', 'Rol', 'Rolle'],
	'members.actions': ['Actions', 'Actions', 'Acciones', 'Aktionen'],
	'members.roleFor': ['Role for {email}', 'Rôle de {email}', 'Rol de {email}', 'Rolle für {email}'],
	'members.remove': ['Remove', 'Retirer', 'Quitar', 'Entfernen'],
	'members.removeLabel': ['Remove {email}', 'Retirer {email}', 'Quitar a {email}', '{email} entfernen'],
	'members.add': ['Add a member', 'Ajouter un membre', 'Añadir un miembro', 'Mitglied hinzufügen'],
	'members.user': ['User', 'Utilisateur', 'Usuario', 'Benutzer'],
	'members.userHelp': [
		'Only people with a Mycelium account can be added.',
		'Seules les personnes ayant un compte Mycelium peuvent être ajoutées.',
		'Solo se pueden añadir personas con una cuenta de Mycelium.',
		'Nur Personen mit einem Mycelium-Konto können hinzugefügt werden.'
	],
	'members.select': ['Select a user…', 'Choisir un utilisateur…', 'Selecciona un usuario…', 'Benutzer auswählen…'],
	'members.addButton': ['Add', 'Ajouter', 'Añadir', 'Hinzufügen'],
	'members.adding': ['Adding…', 'Ajout…', 'Añadiendo…', 'Wird hinzugefügt…'],
	'members.added': ['{email} added to {space}.', '{email} ajouté à {space}.', '{email} añadido a {space}.', '{email} zu {space} hinzugefügt.'],
	'members.theSpace': ['the space', 'l’espace', 'el espacio', 'den Bereich'],
	'members.addFailed': ['Could not add that member.', 'Impossible d’ajouter ce membre.', 'No se pudo añadir a ese miembro.', 'Das Mitglied konnte nicht hinzugefügt werden.'],
	'members.roleChanged': ['{email} is now {role}.', '{email} est maintenant {role}.', '{email} ahora es {role}.', '{email} ist jetzt {role}.'],
	'members.roleFailed': ['Could not change that role.', 'Impossible de changer ce rôle.', 'No se pudo cambiar ese rol.', 'Die Rolle konnte nicht geändert werden.'],
	'members.removed': ['{email} removed.', '{email} retiré.', '{email} quitado.', '{email} entfernt.'],
	'members.removeFailed': ['Could not remove that member.', 'Impossible de retirer ce membre.', 'No se pudo quitar a ese miembro.', 'Das Mitglied konnte nicht entfernt werden.'],
	'members.thisMember': ['this member', 'ce membre', 'este miembro', 'dieses Mitglied'],
	'members.removeTitle': ['Remove {email}?', 'Retirer {email} ?', '¿Quitar a {email}?', '{email} entfernen?'],
	'members.removeText': [
		'They lose access to this space’s memory, rules and skills on their next sync. Anything they wrote stays.',
		'Cette personne perd l’accès à la mémoire, aux règles et aux skills de cet espace à sa prochaine synchronisation. Ce qu’elle a écrit reste.',
		'Pierde el acceso a la memoria, las reglas y las habilidades de este espacio en su próxima sincronización. Lo que escribió permanece.',
		'Die Person verliert beim nächsten Sync den Zugriff auf Gedächtnis, Regeln und Skills dieses Bereichs. Was sie geschrieben hat, bleibt.'
	],
	'members.removeConfirm': ['Remove member', 'Retirer le membre', 'Quitar miembro', 'Mitglied entfernen'],
	'members.keep': ['Keep them', 'Les garder', 'Mantenerlos', 'Behalten']
} as const satisfies Record<string, Row>;
