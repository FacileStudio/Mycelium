import type { Row } from '../row';

/** The Antenne bus settings. */
export const bus = {
	'bus.adminOnly': ['Admin only', 'Réservé aux administrateurs', 'Solo administradores', 'Nur für Admins'],
	'bus.adminOnlyText': [
		'Bus settings are visible to administrators. Sign in with an admin account to change them.',
		'Les réglages du bus sont visibles par les administrateurs. Connectez-vous avec un compte administrateur pour les modifier.',
		'Los ajustes del bus son visibles para los administradores. Inicia sesión con una cuenta de administrador para cambiarlos.',
		'Die Bus-Einstellungen sind für Administratoren sichtbar. Melde dich mit einem Admin-Konto an, um sie zu ändern.'
	],
	'bus.loadFailed': [
		'Could not load the bus settings',
		'Impossible de charger les réglages du bus',
		'No se pudieron cargar los ajustes del bus',
		'Bus-Einstellungen konnten nicht geladen werden'
	],
	'bus.saved': ['Bus settings saved.', 'Réglages du bus enregistrés.', 'Ajustes del bus guardados.', 'Bus-Einstellungen gespeichert.'],
	'bus.saveFailed': ['Save failed', 'Échec de l’enregistrement', 'Error al guardar', 'Speichern fehlgeschlagen'],
	'bus.notSaved': ['Not saved', 'Non enregistré', 'No guardado', 'Nicht gespeichert'],
	'bus.disabled': ['Disabled', 'Désactivé', 'Desactivado', 'Deaktiviert'],
	'bus.connected': ['Connected', 'Connecté', 'Conectado', 'Verbunden'],
	'bus.disconnected': ['Disconnected — {error}', 'Déconnecté : {error}', 'Desconectado: {error}', 'Getrennt – {error}'],
	'bus.connecting': ['Connecting…', 'Connexion…', 'Conectando…', 'Verbindung wird aufgebaut…'],
	'bus.title': ['Antenne bus', 'Bus Antenne', 'Bus de Antenne', 'Antenne-Bus'],
	'bus.help': [
		'One socket to Antenne carries every sealed session this instance emits, so Sablier can turn them into time entries. Apps never talk to each other directly.',
		'Un seul socket vers Antenne transporte chaque session terminée que cette instance émet, pour que Sablier en fasse des entrées de temps. Les applications ne se parlent jamais directement.',
		'Un único socket hacia Antenne transporta cada sesión cerrada que emite esta instancia, para que Sablier las convierta en entradas de tiempo. Las aplicaciones nunca hablan entre sí directamente.',
		'Ein einziger Socket zu Antenne überträgt jede abgeschlossene Sitzung dieser Instanz, damit Sablier daraus Zeiteinträge machen kann. Apps sprechen nie direkt miteinander.'
	],
	'bus.status': ['Status', 'État', 'Estado', 'Status'],
	'bus.statusHelp': [
		'Live, from the socket itself — not from the last save.',
		'En direct, depuis le socket lui-même, pas depuis le dernier enregistrement.',
		'En vivo, desde el propio socket, no desde el último guardado.',
		'Live vom Socket selbst – nicht vom letzten Speichern.'
	],
	'bus.emit': ['Emit sessions', 'Émettre les sessions', 'Emitir sesiones', 'Sitzungen senden'],
	'bus.emitHelp': [
		'Off keeps the config but stops the socket.',
		'Désactivé conserve la config mais coupe le socket.',
		'Desactivado conserva la configuración pero detiene el socket.',
		'Aus behält die Konfiguration, beendet aber den Socket.'
	],
	'bus.emitted': ['Emitted', 'Émises', 'Emitidas', 'Gesendet'],
	'bus.emittedHelp': [
		'Sessions this instance has published since it started.',
		'Sessions que cette instance a publiées depuis son démarrage.',
		'Sesiones que esta instancia ha publicado desde que arrancó.',
		'Sitzungen, die diese Instanz seit dem Start veröffentlicht hat.'
	],
	'bus.sent': ['{count} sent', '{count} envoyées', '{count} enviadas', '{count} gesendet'],
	'bus.outbox': ['Outbox', 'File d’envoi', 'Bandeja de salida', 'Postausgang'],
	'bus.outboxHelp': [
		'Sealed sessions held locally until the socket comes back.',
		'Sessions terminées conservées localement jusqu’au retour du socket.',
		'Sesiones cerradas guardadas localmente hasta que vuelva el socket.',
		'Abgeschlossene Sitzungen, lokal gehalten, bis der Socket zurückkehrt.'
	],
	'bus.pending': ['{count} pending', '{count} en attente', '{count} pendientes', '{count} ausstehend'],
	'bus.alertOutbox': ['Alert outbox', 'File d’alertes', 'Bandeja de alertas', 'Alarm-Postausgang'],
	'bus.alertOutboxHelp': [
		'Usage alerts held locally until the socket comes back.',
		'Alertes de consommation conservées localement jusqu’au retour du socket.',
		'Alertas de uso guardadas localmente hasta que vuelva el socket.',
		'Nutzungsalarme, lokal gehalten, bis der Socket zurückkehrt.'
	],
	'bus.since': ['Emitting since', 'Émission depuis', 'Emitiendo desde', 'Sendet seit'],
	'bus.sinceHelp': [
		'Sessions that ended before this are never republished.',
		'Les sessions terminées avant cette date ne sont jamais republiées.',
		'Las sesiones que terminaron antes de esto nunca se vuelven a publicar.',
		'Sitzungen, die vorher endeten, werden nie erneut veröffentlicht.'
	],
	'bus.connection': ['Connection', 'Connexion', 'Conexión', 'Verbindung'],
	'bus.connectionHelp': [
		'Where Antenne lives and the shared secret that registers this instance with it.',
		'Où se trouve Antenne et le secret partagé qui y enregistre cette instance.',
		'Dónde vive Antenne y el secreto compartido que registra esta instancia en ella.',
		'Wo Antenne läuft und das gemeinsame Geheimnis, mit dem sich diese Instanz dort registriert.'
	],
	'bus.url': ['Instance URL', 'URL de l’instance', 'URL de la instancia', 'Instanz-URL'],
	'bus.urlPinned': [
		'Pinned by ANTENNE_URL. Edit the deployment, not this field.',
		'Imposé par ANTENNE_URL. Modifiez le déploiement, pas ce champ.',
		'Fijado por ANTENNE_URL. Edita el despliegue, no este campo.',
		'Durch ANTENNE_URL festgelegt. Ändere das Deployment, nicht dieses Feld.'
	],
	'bus.urlHelp': [
		'Scheme included. The socket URL is derived from it.',
		'Schéma inclus. L’URL du socket en est déduite.',
		'Con el esquema incluido. La URL del socket se deriva de ella.',
		'Mit Schema. Die Socket-URL wird daraus abgeleitet.'
	],
	'bus.secret': ['Shared secret', 'Secret partagé', 'Secreto compartido', 'Gemeinsames Geheimnis'],
	'bus.secretPinned': [
		'Pinned by ANTENNE_SECRET. Edit the deployment, not this field.',
		'Imposé par ANTENNE_SECRET. Modifiez le déploiement, pas ce champ.',
		'Fijado por ANTENNE_SECRET. Edita el despliegue, no este campo.',
		'Durch ANTENNE_SECRET festgelegt. Ändere das Deployment, nicht dieses Feld.'
	],
	'bus.secretHelp': [
		'Posted to Antenne when the socket registers. Required as soon as emitting is on.',
		'Envoyé à Antenne à l’enregistrement du socket. Requis dès que l’émission est activée.',
		'Se envía a Antenne cuando el socket se registra. Obligatorio en cuanto se activa la emisión.',
		'Wird an Antenne gesendet, wenn sich der Socket registriert. Erforderlich, sobald das Senden aktiv ist.'
	],
	'bus.email': ['Attribution email', 'E-mail d’attribution', 'Correo de atribución', 'Zuordnungs-E-Mail'],
	'bus.emailPinned': [
		'Pinned by ANTENNE_USER_EMAIL. Edit the deployment, not this field.',
		'Imposé par ANTENNE_USER_EMAIL. Modifiez le déploiement, pas ce champ.',
		'Fijado por ANTENNE_USER_EMAIL. Edita el despliegue, no este campo.',
		'Durch ANTENNE_USER_EMAIL festgelegt. Ändere das Deployment, nicht dieses Feld.'
	],
	'bus.emailHelp': [
		'Whose sessions these are, as far as the rest of the suite is concerned.',
		'À qui appartiennent ces sessions, du point de vue du reste de la suite.',
		'De quién son estas sesiones, en lo que respecta al resto de la suite.',
		'Wem diese Sitzungen für den Rest der Suite gehören.'
	],
	'bus.saveConnect': ['Save and connect', 'Enregistrer et connecter', 'Guardar y conectar', 'Speichern und verbinden'],
	'alerts.title': ['Usage alerts', 'Alertes de consommation', 'Alertas de uso', 'Nutzungsalarme'],
	'alerts.help': [
		'When a subscription window crosses your threshold, Mycelium publishes one event to the Antenne. It fires once per window, not once per sync tick — the next alert waits for the window to reset. The Antenne owns what happens after that; Mycelium sends nothing itself. The alert shows up in the Antenne’s activity feed — anything that should act on it needs to subscribe to usage_alert.created on its side.',
		'Quand une fenêtre d’abonnement franchit votre seuil, Mycelium publie un événement vers Antenne. Il part une fois par fenêtre, pas à chaque synchronisation : la prochaine alerte attend la réinitialisation de la fenêtre. Antenne décide de la suite ; Mycelium n’envoie rien lui-même. L’alerte apparaît dans le fil d’activité d’Antenne : ce qui doit y réagir doit s’abonner à usage_alert.created de son côté.',
		'Cuando una ventana de suscripción supera tu umbral, Mycelium publica un evento en Antenne. Se emite una vez por ventana, no en cada ciclo de sincronización: la siguiente alerta espera a que la ventana se reinicie. Antenne decide qué ocurre después; Mycelium no envía nada por sí mismo. La alerta aparece en el flujo de actividad de Antenne: lo que deba reaccionar tiene que suscribirse a usage_alert.created por su cuenta.',
		'Überschreitet ein Abo-Zeitfenster deinen Schwellenwert, veröffentlicht Mycelium ein Ereignis an Antenne. Es wird einmal pro Fenster gesendet, nicht bei jedem Sync-Takt – der nächste Alarm wartet, bis das Fenster zurückgesetzt ist. Was danach passiert, entscheidet Antenne; Mycelium sendet selbst nichts. Der Alarm erscheint im Aktivitätsfeed von Antenne – wer darauf reagieren soll, muss dort usage_alert.created abonnieren.'
	],
	'alerts.emittingOff': ['Emitting is off', 'L’émission est désactivée', 'La emisión está desactivada', 'Das Senden ist aus'],
	'alerts.emittingOffText': [
		'Alerts ride the same socket as sessions. Turn emitting on above and they start flowing.',
		'Les alertes passent par le même socket que les sessions. Activez l’émission ci-dessus et elles partent.',
		'Las alertas usan el mismo socket que las sesiones. Activa la emisión arriba y empezarán a fluir.',
		'Alarme laufen über denselben Socket wie Sitzungen. Aktiviere oben das Senden, dann fließen sie.'
	],
	'alerts.on': ['Alert on usage', 'Alerter sur la consommation', 'Alertar por uso', 'Bei Nutzung alarmieren'],
	'alerts.onHelp': [
		'Off by default. Lowering the threshold later re-arms the current window.',
		'Désactivé par défaut. Baisser le seuil plus tard réarme la fenêtre en cours.',
		'Desactivado por defecto. Bajar el umbral más tarde rearma la ventana actual.',
		'Standardmäßig aus. Ein späteres Senken des Schwellenwerts aktiviert das aktuelle Fenster erneut.'
	],
	'alerts.threshold': ['Threshold', 'Seuil', 'Umbral', 'Schwellenwert'],
	'alerts.thresholdHelp': [
		'Percent of a window’s limit. 80 unless you change it; anything outside 1–100 falls back to it.',
		'Pourcentage de la limite d’une fenêtre. 80 sauf modification ; toute valeur hors 1–100 revient à 80.',
		'Porcentaje del límite de una ventana. 80 salvo que lo cambies; cualquier valor fuera de 1–100 vuelve a 80.',
		'Prozent des Limits eines Fensters. 80, sofern du es nicht änderst; Werte außerhalb von 1–100 fallen darauf zurück.'
	],
	'overrides.title': ['Machine overrides', 'Exceptions par machine', 'Excepciones por máquina', 'Maschinen-Ausnahmen'],
	'overrides.help': [
		'Sessions from a machine listed here are attributed to its own email instead of yours.',
		'Les sessions d’une machine listée ici sont attribuées à son propre e-mail plutôt qu’au vôtre.',
		'Las sesiones de una máquina listada aquí se atribuyen a su propio correo en lugar del tuyo.',
		'Sitzungen einer hier aufgeführten Maschine werden deren eigener E-Mail statt deiner zugeordnet.'
	],
	'overrides.row': ['Override {n}', 'Exception {n}', 'Excepción {n}', 'Ausnahme {n}'],
	'overrides.machine': ['machine', 'machine', 'máquina', 'Maschine'],
	'overrides.machineLabel': ['Machine name', 'Nom de la machine', 'Nombre de la máquina', 'Maschinenname'],
	'overrides.email': ['email', 'e-mail', 'correo', 'E-Mail'],
	'overrides.emailLabel': ['Email for this machine', 'E-mail de cette machine', 'Correo de esta máquina', 'E-Mail für diese Maschine'],
	'overrides.remove': ['Remove override {n}', 'Retirer l’exception {n}', 'Quitar la excepción {n}', 'Ausnahme {n} entfernen'],
	'overrides.add': ['Add an override', 'Ajouter une exception', 'Añadir una excepción', 'Ausnahme hinzufügen'],
	'overrides.addHelp': [
		'One machine per row. Leave the name empty to drop a row on save.',
		'Une machine par ligne. Laissez le nom vide pour supprimer une ligne à l’enregistrement.',
		'Una máquina por fila. Deja el nombre vacío para descartar una fila al guardar.',
		'Eine Maschine pro Zeile. Lass den Namen leer, um eine Zeile beim Speichern zu verwerfen.'
	],
	'overrides.addButton': ['Add override', 'Ajouter l’exception', 'Añadir excepción', 'Ausnahme hinzufügen']
} as const satisfies Record<string, Row>;
