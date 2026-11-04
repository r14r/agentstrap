# Projektstatus

**Stand:** 2026-11-04

- **Repository:** `r14r/agentstrap`
- **Entwicklungsbranch:** `fix/just-install-usr-local-bin-v0.2.1`
- **Version:** `0.2.1` (`VERSION`)
- **CI:** GitHub Actions und PR-Checks vor Merge überprüfen.

## Vorbereitet

- Go-CLI-Grundstruktur mit Hilfe-/Versionsfunktion
- Go-Tests und CI-Workflow
- Einheitliche `just`-Befehle `bootstrap`, `precommit`, `fmt`, `ci`, `build`, `install`, `next`
- Systeminstallation: `/usr/local/bin/agentstrap` (falls nötig mit `sudo`)
- `update-cli.yaml` mit Schema v2
- Bestehende Agent-Skills/MCP-Beispiele beibehalten

## Noch nicht implementiert

- Installation oder Erkennung von AI-Clients
- Skill-/Agent-/MCP-Management
- Generierung von Agent-Konfigurationsdateien

## Nächste Schritte

1. GitHub Actions nach dem Commit prüfen.
2. Feature-PR reviewen und erst mit Freigabe mergen.
3. Client-Erkennung und deklaratives Konfigurationsformat implementieren.

Der Status ist eine Momentaufnahme, keine Live-GitHub-Abfrage.
