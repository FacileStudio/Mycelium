import type { Row } from '../row';

/** Machines, live sessions and session history. */
export const sessions = {
	'range.7d': ['7 days', '7 jours', '7 días', '7 Tage'],
	'range.30d': ['30 days', '30 jours', '30 días', '30 Tage'],
	'range.all': ['All time', 'Tout', 'Todo el tiempo', 'Gesamter Zeitraum'],
	'range.label': ['Time range', 'Période', 'Intervalo de tiempo', 'Zeitraum'],
	'group.project': ['Project', 'Projet', 'Proyecto', 'Projekt'],
	'group.machine': ['Machine', 'Machine', 'Máquina', 'Maschine'],
	'group.agent': ['Agent', 'Agent', 'Agente', 'Agent'],
	'group.model': ['Model', 'Modèle', 'Modelo', 'Modell'],
	'group.label': ['Group sessions by', 'Grouper les sessions par', 'Agrupar sesiones por', 'Sitzungen gruppieren nach'],
	'machines.title': ['Machines', 'Machines', 'Máquinas', 'Maschinen'],
	'machines.subtitle': [
		'Machines and agents syncing with this brain. Connected means a sync in the last ten minutes.',
		'Machines et agents synchronisés avec ce cerveau. « Connectée » signifie une synchronisation dans les dix dernières minutes.',
		'Máquinas y agentes sincronizados con este cerebro. Conectada significa una sincronización en los últimos diez minutos.',
		'Maschinen und Agenten, die mit diesem Gehirn synchronisieren. Verbunden heißt: Sync in den letzten zehn Minuten.'
	],
	'machines.fleet': ['Fleet activity', 'Activité de la flotte', 'Actividad de la flota', 'Aktivität der Flotte'],
	'machines.count': ['Machines', 'Machines', 'Máquinas', 'Maschinen'],
	'machines.workingNow': ['Working now', 'Au travail', 'Trabajando ahora', 'Aktiv jetzt'],
	'machines.synced': [
		'{count} synced in the last 10 minutes',
		'{count} synchronisées ces 10 dernières minutes',
		'{count} sincronizadas en los últimos 10 minutos',
		'{count} in den letzten 10 Minuten synchronisiert'
	],
	'machines.activeTime': ['Active time', 'Temps actif', 'Tiempo activo', 'Aktive Zeit'],
	'machines.tokensOut': ['Tokens out', 'Tokens générés', 'Tokens de salida', 'Ausgabe-Tokens'],
	'machines.perBucket': [
		'Active time per {unit}, by machine',
		'Temps actif par {unit}, par machine',
		'Tiempo activo por {unit}, por máquina',
		'Aktive Zeit pro {unit}, nach Maschine'
	],
	'machines.unit.day': ['day', 'jour', 'día', 'Tag'],
	'machines.unit.month': ['month', 'mois', 'mes', 'Monat'],
	'machines.share': ['Share of active time', 'Part du temps actif', 'Parte del tiempo activo', 'Anteil der aktiven Zeit'],
	'machines.active': ['Active', 'Actif', 'Activo', 'Aktiv'],
	'machines.noActivity': [
		'No recorded activity in this range yet.',
		'Aucune activité enregistrée sur cette période.',
		'Aún no hay actividad registrada en este intervalo.',
		'In diesem Zeitraum wurde noch keine Aktivität aufgezeichnet.'
	],
	'machines.emptyTitle': ['No machines connected yet', 'Aucune machine connectée', 'Aún no hay máquinas conectadas', 'Noch keine Maschinen verbunden'],
	'machines.adminOnly': [
		'Only an admin can list paired machines — the activity above still reflects the whole fleet.',
		'Seul un administrateur peut lister les machines appairées : l’activité ci-dessus reflète toujours toute la flotte.',
		'Solo un administrador puede listar las máquinas emparejadas: la actividad de arriba sigue reflejando toda la flota.',
		'Nur ein Admin kann gekoppelte Maschinen auflisten – die Aktivität oben zeigt weiterhin die ganze Flotte.'
	],
	'machines.pairBefore': ['Run', 'Exécutez', 'Ejecuta', 'Führe'],
	'machines.pairAfter': ['on a machine to pair it.', 'sur une machine pour l’appairer.', 'en una máquina para emparejarla.', 'auf einer Maschine aus, um sie zu koppeln.'],
	'machines.paired': ['Paired machines', 'Machines appairées', 'Máquinas emparejadas', 'Gekoppelte Maschinen'],
	'machines.never': ['never', 'jamais', 'nunca', 'nie'],
	'machines.connected': ['connected', 'connectée', 'conectada', 'verbunden'],
	'machines.neverSynced': ['never synced', 'jamais synchronisée', 'nunca sincronizada', 'nie synchronisiert'],
	'machines.idle': ['idle · {ago}', 'inactive · {ago}', 'inactiva · {ago}', 'inaktiv · {ago}'],
	'machines.lastSync': ['Last sync', 'Dernière synchro', 'Última sincronización', 'Letzter Sync'],
	'machines.activeInRange': ['Active in range', 'Actif sur la période', 'Activo en el intervalo', 'Aktiv im Zeitraum'],
	'machines.pairedOn': ['Paired', 'Appairée', 'Emparejada', 'Gekoppelt'],
	'machines.secondsAgo': ['{count}s ago', 'il y a {count} s', 'hace {count} s', 'vor {count} Sek.'],
	'sessions.title': ['Sessions', 'Sessions', 'Sesiones', 'Sitzungen'],
	'sessions.subtitle': [
		'Agent work sessions recorded across every machine that syncs with this brain.',
		'Sessions de travail des agents enregistrées sur chaque machine synchronisée avec ce cerveau.',
		'Sesiones de trabajo de los agentes registradas en cada máquina que sincroniza con este cerebro.',
		'Arbeitssitzungen der Agenten, aufgezeichnet auf jeder Maschine, die mit diesem Gehirn synchronisiert.'
	],
	'sessions.totals': ['Totals', 'Totaux', 'Totales', 'Summen'],
	'sessions.totalsHelp': [
		'Sealed sessions only — whatever is running now lands here when it ends.',
		'Sessions terminées uniquement : ce qui tourne maintenant arrive ici à la fin.',
		'Solo sesiones cerradas: lo que se está ejecutando ahora llega aquí cuando termina.',
		'Nur abgeschlossene Sitzungen – was gerade läuft, erscheint hier nach dem Ende.'
	],
	'sessions.tab.overview': ['Overview', 'Vue d’ensemble', 'Resumen', 'Übersicht'],
	'sessions.tab.history': ['History', 'Historique', 'Historial', 'Verlauf'],
	'sessions.tab.usage': ['Usage', 'Consommation', 'Uso', 'Nutzung'],
	'sessions.emptyTitle': ['No sessions recorded yet', 'Aucune session enregistrée', 'Aún no hay sesiones registradas', 'Noch keine Sitzungen aufgezeichnet'],
	'sessions.emptyText': [
		'Machines record agent activity automatically once they run mycelium v0.5 or later.',
		'Les machines enregistrent l’activité des agents automatiquement dès qu’elles exécutent mycelium v0.5 ou plus récent.',
		'Las máquinas registran la actividad de los agentes automáticamente cuando ejecutan mycelium v0.5 o posterior.',
		'Maschinen zeichnen die Agentenaktivität automatisch auf, sobald sie mycelium v0.5 oder neuer ausführen.'
	],
	'sessions.windows': ['Subscription windows', 'Fenêtres d’abonnement', 'Ventanas de suscripción', 'Abo-Zeitfenster'],
	'sessions.windowsHelp': [
		'Usage against subscription limits recorded by each agent.',
		'Consommation par rapport aux limites d’abonnement, enregistrée par chaque agent.',
		'Uso frente a los límites de suscripción registrado por cada agente.',
		'Nutzung gegenüber den Abo-Limits, von jedem Agenten aufgezeichnet.'
	],
	'live.title': ['Running now', 'En cours', 'En ejecución ahora', 'Läuft gerade'],
	'live.refresh': ['Refreshed every thirty seconds.', 'Actualisé toutes les trente secondes.', 'Se actualiza cada treinta segundos.', 'Wird alle dreißig Sekunden aktualisiert.'],
	'live.none': ['No session is running.', 'Aucune session en cours.', 'No hay ninguna sesión en ejecución.', 'Es läuft keine Sitzung.'],
	'live.offline': ['machine offline', 'machine hors ligne', 'máquina sin conexión', 'Maschine offline'],
	'live.active': ['active', 'active', 'activa', 'aktiv'],
	'live.idle': ['idle {minutes}m', 'inactive {minutes} min', 'inactiva {minutes} min', 'inaktiv {minutes} Min.'],
	'live.out': ['{tokens} out', '{tokens} générés', '{tokens} de salida', '{tokens} Ausgabe'],
	'overview.sessions': ['Sessions', 'Sessions', 'Sesiones', 'Sitzungen'],
	'overview.cacheRead': ['Cache read', 'Lecture du cache', 'Lectura de caché', 'Cache-Lesezugriffe'],
	'overview.perBucketBy': [
		'Active time per {unit}, by {group}',
		'Temps actif par {unit}, par {group}',
		'Tiempo activo por {unit}, por {group}',
		'Aktive Zeit pro {unit}, nach {group}'
	],
	'overview.shareBy': ['Share by {group}', 'Part par {group}', 'Parte por {group}', 'Anteil nach {group}'],
	'overview.activeBy': ['Active time by {group}', 'Temps actif par {group}', 'Tiempo activo por {group}', 'Aktive Zeit nach {group}'],
	'claims.title': ['Active claims', 'Réservations actives', 'Reclamaciones activas', 'Aktive Reservierungen'],
	'claims.subtitle': [
		'In-flight task leases across every machine — release one to let another agent take over.',
		'Réservations de tâches en cours sur chaque machine : libérez-en une pour qu’un autre agent reprenne.',
		'Reservas de tareas en curso en cada máquina: libera una para que otro agente la retome.',
		'Laufende Aufgabenreservierungen auf jeder Maschine – gib eine frei, damit ein anderer Agent übernimmt.'
	],
	'claims.none': ['No active claim.', 'Aucune réservation active.', 'No hay ninguna reclamación activa.', 'Keine aktive Reservierung.'],
	'claims.since': ['since {elapsed}', 'depuis {elapsed}', 'desde hace {elapsed}', 'seit {elapsed}'],
	'claims.idle': ['idle', 'inactive', 'inactiva', 'inaktiv'],
	'claims.release': ['Release', 'Libérer', 'Liberar', 'Freigeben'],
	'claims.releasing': ['Releasing…', 'Libération…', 'Liberando…', 'Wird freigegeben…'],
	'claims.releaseLabel': [
		'Release claim on {project} by {agent}',
		'Libérer la réservation sur {project} par {agent}',
		'Liberar la reclamación en {project} de {agent}',
		'Reservierung für {project} von {agent} freigeben'
	],
	'history.sessions': ['Sessions', 'Sessions', 'Sesiones', 'Sitzungen'],
	'history.active': ['Active', 'Actif', 'Activo', 'Aktiv'],
	'history.tokensIn': ['Tokens in', 'Tokens reçus', 'Tokens de entrada', 'Eingabe-Tokens'],
	'history.tokensOut': ['Tokens out', 'Tokens générés', 'Tokens de salida', 'Ausgabe-Tokens'],
	'history.cost': ['Est. cost', 'Coût est.', 'Coste est.', 'Gesch. Kosten'],
	'history.ended': ['Ended', 'Terminée', 'Finalizada', 'Beendet'],
	'history.project': ['Project', 'Projet', 'Proyecto', 'Projekt'],
	'history.machineAgent': ['Machine / agent', 'Machine / agent', 'Máquina / agente', 'Maschine / Agent'],
	'history.model': ['Model', 'Modèle', 'Modelo', 'Modell'],
	'history.duration': ['Duration', 'Durée', 'Duración', 'Dauer']
} as const satisfies Record<string, Row>;
