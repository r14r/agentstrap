# Architektur

## Ist-Zustand (v0.2.0)

Go-CLI mit minimalem, unabhängig testbarem Einstiegspunkt.

| Komponente | Verantwortung | Pfad |
| --- | --- | --- |
| CLI-Binary | Argumente, Prozess-Exitcode, Build-Version | `cmd/agentstrap/main.go` |
| CLI-Kern | Hilfe, Version, Eingabevalidierung | `internal/cli/` |
| Buildsystem | Format, Vet, Tests, Race-Tests, Build, Installation | `justfile` |
| Update-Orchestrierung | Schema-v2-Workflows auf `just` | `update-cli.yaml` |
| CI | Automatische Prüfungen und Binary-Smoke-Test | `.github/workflows/ci.yml` |
| Agent-Kontext | Instruktionen, Skills, MCP-Beispiele | `AGENTS.md`, `.agents/skills/`, `examples/mcp/` |

`VERSION` ist die Single Source of Truth für Releases; `just build` setzt `main.version` mit `-ldflags`. Direkter `go run`-Aufruf verwendet die Entwicklungskennung `dev`.

## Zielarchitektur (noch nicht implementiert)

- **Client-Adapter:** Codex, Claude Code, Antigravity, optional weitere.
- **Skill-/Agent-Katalog:** deklarative Quellen und lokale Installationsziele.
- **MCP-Registry:** client-spezifische Formate und sichere Secret-Referenzen.
- **Konfigurationsmodell:** versioniertes `agentstrap.yaml` mit Validierung.
- **Plan/Apply:** Vorschau, Diff, Backup, Konflikterkennung und bestätigte Änderungen.
- **Generatoren:** standardisierte `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, `SKILL.md`.

Keine pauschale Synchronisierung aller Client-Verzeichnisse: Jeder Adapter benötigt eine explizite Mapping- und Migrationsstrategie.

## Sicherheitsgrenzen

Eingelesene Konfigurationen und externe Registries sind untrusted input. Keine Secrets ins Git, Log oder Shell-Argumente; Schreiboperationen nur mit explizit ausgewähltem Ziel und späteren Dry-Run/Backup-Mechanismen.

## Entscheidungen

Neue Architekturentscheidungen unter `doc/ai/decisions/` dokumentieren.
