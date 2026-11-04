# Projektstatus

**Stand:** 2026-11-04

- **Repository:** `r14r/agentstrap`
- **Entwicklungsbranch:** `feat/go-cli-bootstrap-v0.2.0`
- **Version:** `0.2.0` (`VERSION`)
- **CI:** GitHub Actions in diesem Branch zu prüfen; Status vor Merge erneut ermitteln.

## Vorbereitet

- Go-CLI-Grundstruktur mit Hilfe-/Versionsfunktion
- Go-Tests und CI-Workflow
- `justfile` inkl. `next` und `next --run`
- `update-cli.yaml` mit Schema v2
- Bestehende Agent-Skills/MCP-Beispiele beibehalten

## Noch nicht implementiert

- Installation oder Erkennung von AI-Clients
- Skill-/Agent-/MCP-Management
- Generierung von Agent-Konfigurationsdateien

## Nächste Schritte

1. PR-Checks und lokalen `just build` prüfen.
2. Feature-PR reviewen und erst mit Freigabe mergen.
3. Client-Erkennung und deklaratives Konfigurationsformat implementieren.

Der Status ist eine Momentaufnahme, keine Live-GitHub-Abfrage.
