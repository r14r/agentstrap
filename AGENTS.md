# Agent Instructions — AgentStrap

## Zweck und Quellen

- Dieses Repository entwickelt **AgentStrap**, eine plattformübergreifende Go-CLI für die Einrichtung von AI-Coding-Umgebungen.
- `README.md` beschreibt den aktuellen Funktionsumfang. `doc/ai/CONTEXT.md` dient als Kontextindex.
- Maßgeblich sind aktueller Quellcode, Tests, `VERSION`, `doc/ai/PROJECT.md` und `doc/ai/ARCHITECTURE.md`.
- Geplante Befehle aus der Roadmap nicht als bereits implementiert bezeichnen.

## Technische Regeln

- Go-Modul: `github.com/r14r/agentstrap`.
- CLI-Einstiegspunkt: `cmd/agentstrap`; Anwendungslogik und Tests unter `internal/`.
- Standardbibliothek bevorzugen, Drittanbieter-Abhängigkeiten nur mit dokumentiertem Nutzen hinzufügen.
- `VERSION` ist kanonisch; `just build` injiziert die Version über Go-`ldflags`.
- Keine Secrets, Tokens, privaten Schlüssel, oder `.env`-Dateien committen.
- Spätere Konfigurationsänderungen an AI-Clients: Standardmäßig Vorschau, Rückfrage und Backup, keine stillen Überschreibungen.

## Änderung und Validierung

- Vor Änderungen Branch, HEAD und betroffene Dateien prüfen.
- Bei Features/Bugfixes passende Go-Tests ergänzen.
- `just ci` umfasst Format-Check, Vet, Tests, Projektvalidierung und Race-Tests.
- `just build` umfasst bereits `just ci`; Checks nicht unnötig mehrfach ausführen.
- `just next` zeigt nächste Schritte; `just next --run` führt lokale Checks aus.
- Für Repository-Änderungen Feature-Branch und Pull Request verwenden.
- **Kein Merge, Release oder Deployment ohne explizite Freigabe.**
- `VERSION`, `CHANGELOG.md` und betroffene Dokumentation synchron halten.
- Tests nicht als erfolgreich melden, bevor sie tatsächlich ausgeführt wurden.

## Dokumentation

- [PROJECT](doc/ai/PROJECT.md): Ziel, Zielgruppe und Umfang
- [ARCHITECTURE](doc/ai/ARCHITECTURE.md): Go- und Adapter-Architektur
- [CONVENTIONS](doc/ai/CONVENTIONS.md): Entwicklungsregeln
- [ROADMAP](doc/ai/ROADMAP.md): geplante Features
- [STATUS](doc/ai/STATUS.md): Momentaufnahme, keine Live-Daten
- [TOOLS](doc/ai/TOOLS.md): MCP- und Client-Sicherheit

## Übernommene Agent Skills

- `.agents/skills/implement-feature/SKILL.md`
- `.agents/skills/code-review/SKILL.md`
- `.agents/skills/release/SKILL.md`
- `.agents/skills/documentation/SKILL.md`
