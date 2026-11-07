import type { Row } from '../row';

/** The markdown viewer, diagrams and the device authorization page. */
export const viewer = {
	'md.note': ['Note', 'Note', 'Nota', 'Hinweis'],
	'md.tip': ['Tip', 'Astuce', 'Consejo', 'Tipp'],
	'md.advantages': ['Advantages', 'Avantages', 'Ventajas', 'Vorteile'],
	'md.important': ['Important', 'Important', 'Importante', 'Wichtig'],
	'md.warning': ['Warning', 'Attention', 'Advertencia', 'Warnung'],
	'md.caution': ['Caution', 'Prudence', 'Precaución', 'Vorsicht'],
	'md.drawbacks': ['Drawbacks', 'Inconvénients', 'Inconvenientes', 'Nachteile'],
	'md.codeCopied': ['Code copied to clipboard', 'Code copié dans le presse-papiers', 'Código copiado al portapapeles', 'Code in die Zwischenablage kopiert'],
	'mermaid.flowchart': ['Flowchart', 'Diagramme', 'Diagrama de flujo', 'Flussdiagramm'],
	'mermaid.syntaxError': ['Syntax Error', 'Erreur de syntaxe', 'Error de sintaxis', 'Syntaxfehler'],
	'mermaid.zoomOut': ['Zoom out', 'Zoom arrière', 'Alejar', 'Verkleinern'],
	'mermaid.zoomIn': ['Zoom in', 'Zoom avant', 'Acercar', 'Vergrößern'],
	'mermaid.resetZoom': ['Reset zoom', 'Réinitialiser le zoom', 'Restablecer zoom', 'Zoom zurücksetzen'],
	'mermaid.fit': ['Fit to screen', 'Ajuster à l’écran', 'Ajustar a la pantalla', 'An Bildschirm anpassen'],
	'mermaid.fitShort': ['Fit', 'Ajuster', 'Ajustar', 'Anpassen'],
	'mermaid.fullscreen': ['Fullscreen', 'Plein écran', 'Pantalla completa', 'Vollbild'],
	'mermaid.fullscreenView': ['Flowchart fullscreen view', 'Diagramme en plein écran', 'Diagrama a pantalla completa', 'Flussdiagramm im Vollbild'],
	'mermaid.canvas': ['Fullscreen flowchart canvas', 'Zone du diagramme en plein écran', 'Lienzo del diagrama a pantalla completa', 'Vollbild-Zeichenfläche des Flussdiagramms'],
	'mermaid.interactive': ['Interactive flowchart diagram', 'Diagramme interactif', 'Diagrama de flujo interactivo', 'Interaktives Flussdiagramm'],
	'mermaid.viewDiagram': ['View Diagram', 'Voir le diagramme', 'Ver diagrama', 'Diagramm anzeigen'],
	'mermaid.viewSource': ['View Source', 'Voir la source', 'Ver código fuente', 'Quelltext anzeigen'],
	'mermaid.rendering': ['Rendering diagram…', 'Rendu du diagramme…', 'Renderizando diagrama…', 'Diagramm wird gerendert…'],
	'authorize.title': ['Authorize a machine — Mycelium', 'Autoriser une machine — Mycelium', 'Autorizar una máquina — Mycelium', 'Maschine autorisieren — Mycelium'],
	'authorize.heading': ['Authorize a machine', 'Autoriser une machine', 'Autorizar una máquina', 'Maschine autorisieren'],
	'authorize.intro': [
		'Confirm the code shown in your terminal to let this machine sync.',
		'Confirmez le code affiché dans votre terminal pour que cette machine puisse se synchroniser.',
		'Confirma el código que aparece en tu terminal para que esta máquina pueda sincronizar.',
		'Bestätige den im Terminal angezeigten Code, damit diese Maschine synchronisieren kann.'
	],
	'authorize.approved': ['Machine authorized', 'Machine autorisée', 'Máquina autorizada', 'Maschine autorisiert'],
	'authorize.approvedText': [
		'can now sync. Return to your terminal — it will finish automatically.',
		'peut maintenant se synchroniser. Retournez dans votre terminal : il se terminera tout seul.',
		'ya puede sincronizar. Vuelve a tu terminal: terminará automáticamente.',
		'kann jetzt synchronisieren. Kehre zu deinem Terminal zurück – es schließt automatisch ab.'
	],
	'authorize.denied': ['Request denied', 'Demande refusée', 'Solicitud denegada', 'Anfrage abgelehnt'],
	'authorize.deniedText': ['No token was issued.', 'Aucun token n’a été émis.', 'No se emitió ningún token.', 'Es wurde kein Token ausgestellt.'],
	'authorize.requesting': ['Machine requesting access', 'Machine demandant l’accès', 'Máquina que solicita acceso', 'Maschine fordert Zugriff an'],
	'authorize.code': ['code', 'code', 'código', 'Code'],
	'authorize.from': ['from', 'depuis', 'desde', 'von'],
	'authorize.warnBefore': ['Only authorize if you just ran', 'N’autorisez que si vous venez d’exécuter', 'Autoriza solo si acabas de ejecutar', 'Autorisiere nur, wenn du gerade'],
	'authorize.warnAfter': [
		'on this machine. Approving grants it ongoing sync access to your shared brain.',
		'sur cette machine. L’autorisation lui donne un accès permanent à la synchronisation de votre cerveau partagé.',
		'en esta máquina. Aprobar le concede acceso continuo de sincronización a tu cerebro compartido.',
		'auf dieser Maschine ausgeführt hast. Die Freigabe gewährt ihr dauerhaften Sync-Zugriff auf dein gemeinsames Gehirn.'
	],
	'authorize.deny': ['Deny', 'Refuser', 'Denegar', 'Ablehnen'],
	'authorize.approve': ['Authorize', 'Autoriser', 'Autorizar', 'Autorisieren'],
	'authorize.authorizing': ['Authorizing…', 'Autorisation…', 'Autorizando…', 'Wird autorisiert…'],
	'authorize.codeLabel': ['Device code', 'Code de l’appareil', 'Código del dispositivo', 'Gerätecode'],
	'authorize.continue': ['Continue', 'Continuer', 'Continuar', 'Weiter'],
	'authorize.checking': ['Checking…', 'Vérification…', 'Comprobando…', 'Wird geprüft…'],
	'authorize.hintDroppedOne': [
		'{list} is not used in device codes, so nothing was added.',
		'{list} n’est pas utilisé dans les codes d’appareil : rien n’a été ajouté.',
		'{list} no se usa en los códigos de dispositivo, así que no se añadió nada.',
		'{list} kommt in Gerätecodes nicht vor, daher wurde nichts hinzugefügt.'
	],
	'authorize.hintDroppedMany': [
		'{list} are not used in device codes, so nothing was added.',
		'{list} ne sont pas utilisés dans les codes d’appareil : rien n’a été ajouté.',
		'{list} no se usan en los códigos de dispositivo, así que no se añadió nada.',
		'{list} kommen in Gerätecodes nicht vor, daher wurde nichts hinzugefügt.'
	],
	'authorize.hintEmpty': [
		'Eight characters, like {example}.',
		'Huit caractères, comme {example}.',
		'Ocho caracteres, como {example}.',
		'Acht Zeichen, wie {example}.'
	],
	'authorize.hintCount': ['{count} of {total} characters.', '{count} caractères sur {total}.', '{count} de {total} caracteres.', '{count} von {total} Zeichen.'],
	'authorize.alreadyUsed': [
		'This code was already {status}. Run mycelium login again for a fresh one.',
		'Ce code a déjà été {status}. Exécutez à nouveau mycelium login pour en obtenir un nouveau.',
		'Este código ya fue {status}. Ejecuta mycelium login de nuevo para obtener uno nuevo.',
		'Dieser Code wurde bereits {status}. Führe mycelium login erneut aus, um einen neuen zu erhalten.'
	],
	'authorize.status.approved': ['approved', 'approuvé', 'aprobado', 'bestätigt'],
	'authorize.status.denied': ['denied', 'refusé', 'denegado', 'abgelehnt'],
	'authorize.status.expired': ['expired', 'expiré', 'caducado', 'abgelaufen'],
	'authorize.needsAdmin': [
		'Approving a machine needs an admin account, and this session is not one. Ask an admin of this Mycelium to approve it.',
		'Autoriser une machine demande un compte administrateur, et cette session n’en est pas un. Demandez à un administrateur de ce Mycelium de l’autoriser.',
		'Aprobar una máquina requiere una cuenta de administrador y esta sesión no lo es. Pide a un administrador de este Mycelium que la apruebe.',
		'Zum Freigeben einer Maschine ist ein Admin-Konto nötig, und diese Sitzung ist keines. Bitte einen Admin dieses Mycelium, sie freizugeben.'
	],
	'authorize.unknownCode': [
		'This server is not holding that code. It may have expired, been used already, or been issued before the server last restarted. Run mycelium login again for a fresh one.',
		'Ce serveur ne connaît pas ce code. Il a pu expirer, être déjà utilisé, ou avoir été émis avant le dernier redémarrage du serveur. Exécutez à nouveau mycelium login pour en obtenir un nouveau.',
		'Este servidor no tiene ese código. Puede haber caducado, haberse usado ya o haberse emitido antes del último reinicio del servidor. Ejecuta mycelium login de nuevo para obtener uno nuevo.',
		'Dieser Server kennt diesen Code nicht. Er kann abgelaufen sein, bereits verwendet worden sein oder vor dem letzten Neustart des Servers ausgestellt worden sein. Führe mycelium login erneut aus, um einen neuen zu erhalten.'
	],
	'authorize.unfinished': [
		'That code is not finished. Device codes are {total} characters, like {example}.',
		'Ce code n’est pas complet. Les codes d’appareil font {total} caractères, comme {example}.',
		'Ese código no está completo. Los códigos de dispositivo tienen {total} caracteres, como {example}.',
		'Dieser Code ist unvollständig. Gerätecodes haben {total} Zeichen, wie {example}.'
	]
} as const satisfies Record<string, Row>;
