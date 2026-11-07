import type { Row } from '../row';

/** Settings: profile, appearance, API tokens, the Antenne bus and the instance. */
export const settings = {
	'settings.title': ['Settings', 'Réglages', 'Ajustes', 'Einstellungen'],
	'settings.subtitle': [
		'Your account, this instance, and everything wired to it.',
		'Votre compte, cette instance et tout ce qui y est relié.',
		'Tu cuenta, esta instancia y todo lo que está conectado a ella.',
		'Dein Konto, diese Instanz und alles, was daran hängt.'
	],
	'settings.sections': ['Settings sections', 'Sections des réglages', 'Secciones de ajustes', 'Abschnitte der Einstellungen'],
	'settings.tab.profile': ['Profile', 'Profil', 'Perfil', 'Profil'],
	'settings.tab.appearance': ['Appearance', 'Apparence', 'Apariencia', 'Darstellung'],
	'settings.tab.api': ['API', 'API', 'API', 'API'],
	'settings.tab.bus': ['Bus', 'Bus', 'Bus', 'Bus'],
	'settings.tab.advanced': ['Advanced', 'Avancé', 'Avanzado', 'Erweitert'],
	'profile.thisMachine': ['This machine', 'Cette machine', 'Esta máquina', 'Diese Maschine'],
	'profile.rules': ['Rules', 'Règles', 'Reglas', 'Regeln'],
	'profile.skills': ['Skills', 'Skills', 'Habilidades', 'Skills'],
	'profile.account': ['Account', 'Compte', 'Cuenta', 'Konto'],
	'profile.managed': [
		'Managed by Facile SSO. Change your name, email or picture at sso.facile.studio/account.',
		'Géré par Facile SSO. Modifiez votre nom, votre e-mail ou votre photo sur sso.facile.studio/account.',
		'Gestionado por Facile SSO. Cambia tu nombre, correo o foto en sso.facile.studio/account.',
		'Verwaltet über Facile SSO. Ändere Name, E-Mail oder Bild unter sso.facile.studio/account.'
	],
	'profile.sharedPassword': [
		'This instance is running on a shared password, so there is no account to manage.',
		'Cette instance utilise un mot de passe partagé : il n’y a pas de compte à gérer.',
		'Esta instancia usa una contraseña compartida, así que no hay ninguna cuenta que gestionar.',
		'Diese Instanz läuft mit einem gemeinsamen Passwort, es gibt also kein Konto zu verwalten.'
	],
	'profile.logout': ['Log out', 'Se déconnecter', 'Cerrar sesión', 'Abmelden'],
	'profile.logoutHelp': [
		'Ends this session in this browser. Paired machines keep syncing.',
		'Termine cette session dans ce navigateur. Les machines appairées continuent à se synchroniser.',
		'Cierra esta sesión en este navegador. Las máquinas emparejadas siguen sincronizando.',
		'Beendet diese Sitzung in diesem Browser. Gekoppelte Maschinen synchronisieren weiter.'
	],
	'profile.loggingOut': ['Signing out…', 'Déconnexion…', 'Cerrando sesión…', 'Abmeldung läuft…'],
	'profile.loggedOut': ['Signed out.', 'Déconnecté.', 'Sesión cerrada.', 'Abgemeldet.'],
	'appearance.theme': ['Theme', 'Thème', 'Tema', 'Design'],
	'appearance.themeHelp': [
		'Applied to this browser. Every device you sign in from keeps its own choice.',
		'Appliqué à ce navigateur. Chaque appareil depuis lequel vous vous connectez garde son propre choix.',
		'Se aplica a este navegador. Cada dispositivo desde el que inicias sesión conserva su propia elección.',
		'Gilt für diesen Browser. Jedes Gerät, mit dem du dich anmeldest, behält seine eigene Auswahl.'
	],
	'appearance.scheme': ['Colour scheme', 'Jeu de couleurs', 'Esquema de color', 'Farbschema'],
	'appearance.schemeHelp': [
		'System follows your OS for as long as the app is open.',
		'« Système » suit votre OS tant que l’application est ouverte.',
		'«Sistema» sigue tu sistema operativo mientras la aplicación esté abierta.',
		'„System“ folgt deinem Betriebssystem, solange die App geöffnet ist.'
	],
	'appearance.system': ['System', 'Système', 'Sistema', 'System'],
	'appearance.light': ['Light', 'Clair', 'Claro', 'Hell'],
	'appearance.dark': ['Dark', 'Sombre', 'Oscuro', 'Dunkel'],
	'advanced.auth.unknown': ['Unknown', 'Inconnu', 'Desconocido', 'Unbekannt'],
	'advanced.auth.sso': ['SSO only', 'SSO uniquement', 'Solo SSO', 'Nur SSO'],
	'advanced.auth.password': ['Shared password only', 'Mot de passe partagé uniquement', 'Solo contraseña compartida', 'Nur gemeinsames Passwort'],
	'advanced.instance': ['Instance', 'Instance', 'Instancia', 'Instanz'],
	'advanced.instanceHelp': [
		'Useful when you file a bug against this self-hosted install.',
		'Utile pour signaler un bug sur cette installation auto-hébergée.',
		'Útil cuando reportas un error de esta instalación autoalojada.',
		'Nützlich, wenn du einen Fehler dieser selbst gehosteten Installation meldest.'
	],
	'advanced.signIn': ['Sign-in', 'Connexion', 'Inicio de sesión', 'Anmeldung'],
	'advanced.signInHelp': [
		'How this instance authenticates people.',
		'Comment cette instance authentifie les personnes.',
		'Cómo autentica esta instancia a las personas.',
		'Wie diese Instanz Personen authentifiziert.'
	],
	'advanced.machineHelp': [
		'The name the server reports for itself.',
		'Le nom que le serveur donne de lui-même.',
		'El nombre que el servidor da de sí mismo.',
		'Der Name, den der Server von sich meldet.'
	],
	'advanced.rulesLoaded': ['Rules loaded', 'Règles chargées', 'Reglas cargadas', 'Geladene Regeln'],
	'advanced.rulesHelp': [
		'Concatenated into every agent config.',
		'Concaténées dans la config de chaque agent.',
		'Concatenadas en la configuración de cada agente.',
		'In jede Agentenkonfiguration eingefügt.'
	],
	'advanced.skillsLoaded': ['Skills loaded', 'Skills chargés', 'Habilidades cargadas', 'Geladene Skills'],
	'advanced.skillsHelp': [
		'Installed into each agent’s skill format.',
		'Installés dans le format de skill de chaque agent.',
		'Instaladas en el formato de habilidad de cada agente.',
		'Im Skill-Format jedes Agenten installiert.'
	],
	'advanced.cli': ['Command line', 'Ligne de commande', 'Línea de comandos', 'Kommandozeile'],
	'advanced.cliHelp': [
		'The CLI is the same binary as the server, so a machine only needs the install script and a token.',
		'Le CLI est le même binaire que le serveur : une machine n’a besoin que du script d’installation et d’un token.',
		'La CLI es el mismo binario que el servidor, así que una máquina solo necesita el script de instalación y un token.',
		'Die CLI ist dasselbe Binary wie der Server, eine Maschine braucht also nur das Installationsskript und einen Token.'
	],
	'advanced.install': ['Install', 'Installation', 'Instalación', 'Installation'],
	'advanced.pair': ['Pair this machine', 'Appairer cette machine', 'Emparejar esta máquina', 'Diese Maschine koppeln'],
	'advanced.pairHelp': [
		'Sign in with your Facile account when it asks.',
		'Connectez-vous avec votre compte Facile quand il le demande.',
		'Inicia sesión con tu cuenta de Facile cuando te lo pida.',
		'Melde dich mit deinem Facile-Konto an, wenn danach gefragt wird.'
	],
	'usage.plan': ['Plan usage', 'Consommation du forfait', 'Uso del plan', 'Plannutzung'],
	'usage.none': [
		'Nothing recorded yet. A machine reports its subscription windows once',
		'Rien d’enregistré pour l’instant. Une machine signale ses fenêtres d’abonnement dès que',
		'Aún no se ha registrado nada. Una máquina informa de sus ventanas de suscripción cuando',
		'Noch nichts aufgezeichnet. Eine Maschine meldet ihre Abo-Zeitfenster, sobald'
	],
	'usage.noneAfter': [
		'has wired the status line and an agent session has made its first request.',
		'a configuré la ligne d’état et qu’une session d’agent a fait sa première requête.',
		'haya configurado la línea de estado y una sesión de agente haya hecho su primera petición.',
		'die Statuszeile eingerichtet hat und eine Agentensitzung ihre erste Anfrage gestellt hat.'
	],
	'usage.lastReported': [
		'Last reported {age}. The status line only reports while Claude Code is running, so a gap means nobody reported — not that nothing was used.',
		'Dernier rapport {age}. La ligne d’état ne rapporte que pendant l’exécution de Claude Code : un trou signifie que personne n’a rapporté, pas que rien n’a été utilisé.',
		'Último informe {age}. La línea de estado solo informa mientras Claude Code se ejecuta, así que un hueco significa que nadie informó, no que no se usara nada.',
		'Zuletzt gemeldet {age}. Die Statuszeile meldet nur, während Claude Code läuft; eine Lücke heißt also, dass niemand gemeldet hat – nicht, dass nichts genutzt wurde.'
	],
	'usage.noWindow': ['No window reported.', 'Aucune fenêtre signalée.', 'No se informó de ninguna ventana.', 'Kein Zeitfenster gemeldet.'],
	'usage.burnDown': ['Burn-down', 'Consommation dans le temps', 'Consumo en el tiempo', 'Verbrauchsverlauf'],
	'usage.expired': [
		'{pct} when last seen · window has since reset',
		'{pct} à la dernière mesure · la fenêtre s’est réinitialisée depuis',
		'{pct} en la última lectura · la ventana se ha reiniciado desde entonces',
		'{pct} beim letzten Stand · das Fenster wurde inzwischen zurückgesetzt'
	],
	'usage.resetsIn': ['resets in {span}', 'se réinitialise dans {span}', 'se reinicia en {span}', 'wird in {span} zurückgesetzt'],
	'usage.resetUnknown': ['reset time unknown', 'heure de réinitialisation inconnue', 'hora de reinicio desconocida', 'Zeitpunkt der Rücksetzung unbekannt'],
	'usage.meterExpired': [
		'{name}: {pct} when last seen, window has since reset',
		'{name} : {pct} à la dernière mesure, la fenêtre s’est réinitialisée depuis',
		'{name}: {pct} en la última lectura, la ventana se ha reiniciado desde entonces',
		'{name}: {pct} beim letzten Stand, das Fenster wurde inzwischen zurückgesetzt'
	],
	'usage.meterUsed': ['{name} used', '{name} utilisé', '{name} usado', '{name} genutzt'],
	'usage.ago': ['{span} ago', 'il y a {span}', 'hace {span}', 'vor {span}'],
	'api.endpoint': ['Endpoint', 'Point d’accès', 'Endpoint', 'Endpunkt'],
	'api.endpointHelp': [
		'Point the CLI or any HTTP client here. Not a secret — the token is.',
		'Pointez le CLI ou tout client HTTP ici. Ce n’est pas un secret : le token, si.',
		'Apunta aquí la CLI o cualquier cliente HTTP. No es un secreto: el token sí.',
		'Richte die CLI oder jeden HTTP-Client hierher. Das ist kein Geheimnis – der Token schon.'
	],
	'api.baseUrl': ['Base URL', 'URL de base', 'URL base', 'Basis-URL'],
	'api.connect': ['Connect your agents', 'Connectez vos agents', 'Conecta tus agentes', 'Verbinde deine Agenten'],
	'api.connectHelp': [
		'Paste this into an agent’s master prompt (CLAUDE.md, AGENTS.md, GEMINI.md…) so it knows how to read, write and sync the shared brain.',
		'Collez ceci dans le prompt principal d’un agent (CLAUDE.md, AGENTS.md, GEMINI.md…) pour qu’il sache lire, écrire et synchroniser le cerveau partagé.',
		'Pega esto en el prompt principal de un agente (CLAUDE.md, AGENTS.md, GEMINI.md…) para que sepa leer, escribir y sincronizar el cerebro compartido.',
		'Füge das in den Hauptprompt eines Agenten ein (CLAUDE.md, AGENTS.md, GEMINI.md…), damit er das gemeinsame Gehirn lesen, schreiben und synchronisieren kann.'
	],
	'api.copyPrompt': ['Copy prompt', 'Copier le prompt', 'Copiar prompt', 'Prompt kopieren'],
	'api.promptCopied': ['Agent prompt copied.', 'Prompt copié.', 'Prompt copiado.', 'Prompt kopiert.'],
	'api.clipboardDenied': [
		'The browser refused clipboard access. Select the text and copy it manually.',
		'Le navigateur a refusé l’accès au presse-papiers. Sélectionnez le texte et copiez-le à la main.',
		'El navegador rechazó el acceso al portapapeles. Selecciona el texto y cópialo manualmente.',
		'Der Browser hat den Zugriff auf die Zwischenablage verweigert. Markiere den Text und kopiere ihn manuell.'
	],
	'api.promptNote': [
		'It assumes the mycelium CLI is installed and logged in with one of the tokens below.',
		'Il suppose que le CLI mycelium est installé et connecté avec l’un des tokens ci-dessous.',
		'Supone que la CLI de mycelium está instalada y con la sesión iniciada con uno de los tokens de abajo.',
		'Er setzt voraus, dass die mycelium-CLI installiert und mit einem der Token unten angemeldet ist.'
	],
	'api.tokens': ['API tokens', 'Tokens d’API', 'Tokens de API', 'API-Token'],
	'api.tokensHelp': [
		'One token per machine. Revoking one never touches the others.',
		'Un token par machine. En révoquer un ne touche jamais les autres.',
		'Un token por máquina. Revocar uno nunca afecta a los demás.',
		'Ein Token pro Maschine. Das Widerrufen eines Tokens berührt die anderen nie.'
	],
	'api.newToken': ['New token', 'Nouveau token', 'Nuevo token', 'Neuer Token'],
	'api.adminOnly': [
		'Only instance administrators can view and issue sync tokens.',
		'Seuls les administrateurs de l’instance peuvent voir et émettre des tokens de synchronisation.',
		'Solo los administradores de la instancia pueden ver y emitir tokens de sincronización.',
		'Nur Instanz-Administratoren können Sync-Token sehen und ausstellen.'
	],
	'api.noTokens': [
		'No tokens yet. The CLI needs one before a machine can sync with this brain.',
		'Aucun token pour l’instant. Le CLI en a besoin avant qu’une machine puisse se synchroniser avec ce cerveau.',
		'Aún no hay tokens. La CLI necesita uno antes de que una máquina pueda sincronizar con este cerebro.',
		'Noch keine Token. Die CLI braucht einen, bevor eine Maschine mit diesem Gehirn synchronisieren kann.'
	],
	'api.scope': ['Scope', 'Portée', 'Ámbito', 'Bereich'],
	'api.created': ['Created', 'Créé', 'Creado', 'Erstellt'],
	'api.lastSeen': ['Last seen', 'Dernière activité', 'Última actividad', 'Zuletzt gesehen'],
	'api.machine': ['machine', 'machine', 'máquina', 'Maschine'],
	'api.never': ['never', 'jamais', 'nunca', 'nie'],
	'api.revokeLabel': ['Revoke {name}', 'Révoquer {name}', 'Revocar {name}', '{name} widerrufen'],
	'api.drawerTitle': ['New API token', 'Nouveau token d’API', 'Nuevo token de API', 'Neuer API-Token'],
	'api.copyNow': ['Copy it now', 'Copiez-le maintenant', 'Cópialo ahora', 'Jetzt kopieren'],
	'api.copyNowText': [
		'This is the only time the token is shown. Mycelium stores a hash, so it cannot be shown again — losing it means issuing a new one.',
		'Le token n’est affiché qu’une seule fois. Mycelium n’en conserve qu’un hash et ne peut donc pas le réafficher : le perdre oblige à en émettre un nouveau.',
		'Esta es la única vez que se muestra el token. Mycelium guarda un hash, así que no puede mostrarse de nuevo: perderlo significa emitir uno nuevo.',
		'Der Token wird nur dieses eine Mal angezeigt. Mycelium speichert nur einen Hash und kann ihn nicht erneut zeigen – wer ihn verliert, muss einen neuen ausstellen.'
	],
	'api.tokenLabel': ['Token', 'Token', 'Token', 'Token'],
	'api.tokenHelp': [
		'Run `mycelium login --token` on the machine and paste it when asked.',
		'Exécutez `mycelium login --token` sur la machine et collez-le quand il est demandé.',
		'Ejecuta `mycelium login --token` en la máquina y pégalo cuando te lo pida.',
		'Führe `mycelium login --token` auf der Maschine aus und füge ihn ein, wenn danach gefragt wird.'
	],
	'api.done': ['Done', 'Terminé', 'Hecho', 'Fertig'],
	'api.nameHelp': [
		'Where the token will live — a machine, a pipeline, a script.',
		'Où vivra le token : une machine, un pipeline, un script.',
		'Dónde vivirá el token: una máquina, un pipeline, un script.',
		'Wo der Token liegen wird – eine Maschine, eine Pipeline, ein Skript.'
	],
	'api.createToken': ['Create token', 'Créer le token', 'Crear token', 'Token erstellen'],
	'api.createFailed': ['Could not create the token.', 'Impossible de créer le token.', 'No se pudo crear el token.', 'Der Token konnte nicht erstellt werden.'],
	'api.notCreated': ['Token not created', 'Token non créé', 'Token no creado', 'Token nicht erstellt'],
	'api.loadFailed': ['Could not load tokens.', 'Impossible de charger les tokens.', 'No se pudieron cargar los tokens.', 'Token konnten nicht geladen werden.'],
	'api.revoked': ['Token “{name}” revoked.', 'Token « {name} » révoqué.', 'Token «{name}» revocado.', 'Token „{name}“ widerrufen.'],
	'api.revokeFailed': ['Could not revoke that token.', 'Impossible de révoquer ce token.', 'No se pudo revocar ese token.', 'Der Token konnte nicht widerrufen werden.'],
	'api.revokeTitle': ['Revoke “{name}”?', 'Révoquer « {name} » ?', '¿Revocar «{name}»?', '„{name}“ widerrufen?'],
	'api.revokeText': [
		'That machine stops syncing immediately and any script using the token starts failing. It cannot be un-revoked — issue a new token and log in again.',
		'Cette machine cesse de se synchroniser immédiatement et tout script utilisant le token échoue. Impossible d’annuler : émettez un nouveau token et reconnectez-vous.',
		'Esa máquina deja de sincronizar de inmediato y cualquier script que use el token empieza a fallar. No se puede deshacer: emite un token nuevo e inicia sesión de nuevo.',
		'Diese Maschine synchronisiert sofort nicht mehr, und jedes Skript mit dem Token schlägt fehl. Das lässt sich nicht rückgängig machen – stelle einen neuen Token aus und melde dich neu an.'
	],
	'api.revokeConfirm': ['Revoke token', 'Révoquer le token', 'Revocar token', 'Token widerrufen']
} as const satisfies Record<string, Row>;
