# Projektstatus

**Stand:** 2026-11-04

- **Repository:** `r14r/agentstrap`
- **Entwicklungsbranch:** `feat/catalog-install-inspect-v0.3.0`
- **Version:** `0.3.0` (`VERSION`)
- **CI:** GitHub Actions und PR-Checks vor Merge überprüfen.

## Vorbereitet

- Go-CLI-Grundstruktur mit Hilfe-/Versionsfunktion
- Go-Tests und CI-Workflow
- Thematischer Katalog, Projektinstallation und Repository-Inspektion
- Einheitliche `just`-Befehle `bootstrap`, `precommit`, `fmt`, `ci`, `build`, `install`, `next`
- Systeminstallation: `/usr/local/bin/agentstrap` (falls nötig mit `sudo`)
- `update-cli.yaml` mit Schema v2
- Bestehende Agent-Skills/MCP-Beispiele beibehalten

## Noch nicht implementiert

- Installation oder Erkennung von AI-Clients
- Externe Kataloge, Remote-Installation und fortgeschrittenes Skill-/Agent-/MCP-Management
- Automatische Generierung der durch die Inspektion vorgeschlagenen Dateien

## Nächste Schritte

1. GitHub Actions nach dem Commit prüfen.
2. Feature-PR reviewen und erst mit Freigabe mergen.
3. Client-Erkennung und deklaratives Konfigurationsformat implementieren.

Der Status ist eine Momentaufnahme, keine Live-GitHub-Abfrage.
