# Agent Instructions

## Zweck und Quellen
- Dieses Repository ist die Vorlage für neue KI-unterstützte Projekte.
- `README.md` beschreibt die Benutzung; `doc/ai/CONTEXT.md` ist der Kontext-Wegweiser.
- Lies nur die für die Aufgabe relevanten Dokumente und Quelldateien.
- Das Repository und die Tests sind maßgeblich; Chatverläufe sind keine verifizierte Ist-Dokumentation.

## Änderungen
- Prüfe vor jeder Änderung den aktuellen Branch, HEAD-Commit und die betroffenen Dateien.
- Formuliere bei größeren Aufgaben einen Plan; halte Änderungen klein und prüfbar.
- Erhalte vorhandene Funktionen und Schnittstellen, sofern nichts anderes beauftragt ist.
- Neue Abhängigkeiten, Frameworks oder MCP-Server nur mit nachvollziehbarem Bedarf.
- Keine Secrets, persönlichen Zugangsdaten, privaten Schlüssel oder `.env`-Dateien committen.

## Konventionen
- Versionsstand in `VERSION` nach SemVer pflegen; Änderungshistorie in `CHANGELOG.md`.
- Dokumentationsänderungen und Codeänderungen synchron halten.
- Architekturentscheidungen unter `doc/ai/decisions/` mit Status und Begründung festhalten.
- Deutsch für die fachliche Dokumentation; Englisch für Quellcode und technische Bezeichner.
- Keine erfundenen Testergebnisse, Commits, Fähigkeiten oder externen Daten.

## Validierung und Lieferung
- Nach Änderungen `just validate` oder `python3 scripts/validate_template.py` ausführen.
- Für Änderungen am Repository nach Möglichkeit Branch → Commit → Pull Request verwenden.
- Kein Merge, Release oder externes Deployment ohne ausdrückliche Freigabe.
- Im Ergebnis Änderungen, ausgeführte Prüfungen, nicht ausgeführte Prüfungen und offene Punkte nennen.

## Kontext-Routing
- Ziel und Scope: `doc/ai/PROJECT.md`
- Kontextindex: `doc/ai/CONTEXT.md`
- Architektur und Konventionen: `doc/ai/ARCHITECTURE.md`, `doc/ai/CONVENTIONS.md`
- Verifizierter Arbeitsstand: `doc/ai/STATUS.md`
- Prioritäten: `doc/ai/ROADMAP.md`
- Tools/MCP: `doc/ai/TOOLS.md`
- ChatGPT-Projektanweisungen: `doc/ai/CHATGPT_PROJECT.md`

## Skills
Bedarfsweise verwenden, sofern vom Agent-Client unterstützt:
- `implement-feature`: Feature planen, implementieren und prüfen.
- `code-review`: Risiko- und Code-Review.
- `release`: SemVer und Release-Vorbereitung.
- `documentation`: Dokumentation und Kontext aktualisieren.
