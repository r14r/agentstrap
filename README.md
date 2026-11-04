# AgentStrap

**One CLI. Every AI coding environment.**

AgentStrap ist eine in Go geschriebene CLI zur zukünftigen einheitlichen Einrichtung von AI-Coding-Clients (Codex, Claude Code, Antigravity), Skills, Agents, MCP-Servern und projektspezifischen Agent-Konfigurationsdateien.

**Stand v0.2.0:** Das Repository enthält eine lauffähige, getestete CLI-Basis mit `help` und `version`. Client-/Skill-/Agent-/MCP-Verwaltung und Dateigenerierung sind **noch nicht implementiert**; siehe [Roadmap](doc/ai/ROADMAP.md).

## Voraussetzungen

- Go 1.23 oder neuer
- [just](https://github.com/casey/just)
- Python 3.11+ für die aus der Projektvorlage übernommene Strukturvalidierung
- Optional: GitHub CLI `gh` für die PR-Übersicht in `just next` und `update-cli` für standardisierte Workflows

## Quick start

```sh
git clone https://github.com/r14r/agentstrap.git
cd agentstrap

just bootstrap
just build
./dist/agentstrap --help
./dist/agentstrap version
```

Für eine lokale Installation in `~/.local/bin`:

```sh
just install
# Optional ein anderes Installationsverzeichnis:
INSTALL_DIR=/some/bin just install
```

## Entwicklungsbefehle

| Befehl | Wirkung |
| --- | --- |
| `just` | Alle verfügbaren Rezepte anzeigen |
| `just bootstrap` | Go-Module vorbereiten |
| `just fmt` | Go-Code formatieren |
| `just fmt-check` | Formatierung ohne Änderungen prüfen |
| `just test` | Go-Unit-Tests |
| `just validate` | Projektstruktur und Agent-Skills prüfen |
| `just precommit` | Format, `go vet`, Tests und Projektvalidierung |
| `just ci` | `precommit` plus Race-Tests |
| `just build` | `ci`, Binary bauen, Version und Hilfe prüfen |
| `just run --help` | CLI direkt über `go run` ausführen |
| `just next` | Git-Stand, offene PRs und nächste Prüfschritte anzeigen |
| `just next --run` | Zusätzlich `just build` und gegebenenfalls PR-Checks ausführen |
| `just update` | `git pull --ff-only`, danach `just build` |
| `just clean` | `dist/` entfernen |

`just build` führt `just ci` bereits aus; ein separates `just precommit` oder `just test` davor ist nicht erforderlich.

## Update CLI

Die Datei [update-cli.yaml](update-cli.yaml) verwendet **schemaVersion 2**. Sie definiert die Workflows `bootstrap`, `setup`, `precommit`, `test`, `ci`, `build`, `install`, `run` und `clean`. Alle Workflows verweisen auf die kanonischen `just`-Rezepte; so laufen keine abweichenden Prüfketten.

Die Konfiguration verweist auf das GitHub-Repository `r14r/agentstrap` auf `main`. Installation erfolgt ohne root-Rechte standardmäßig nach `~/.local/bin`.

## Struktur

```text
cmd/agentstrap/main.go       CLI-Einstiegspunkt und eingebettete Build-Version
internal/cli/                Kommandoverarbeitung und Tests
go.mod                       Go-Modul
VERSION                      SemVer-Version (Single Source of Truth)
justfile                     Build-, Prüf-, Installations- und Update-Rezepte
update-cli.yaml              Schema-v2-Workflows für Update CLI
.github/workflows/ci.yml     Automatische Prüfungen und Build-Smoke-Test
scripts/next.sh              Lokale/PR-Prüfempfehlungen
scripts/validate_template.py Validierung der übernommenen AI-Projektstruktur
AGENTS.md                    Anweisungen für Coding-Agents
.agents/skills/              Wiederverwendbare Projekt-Skills
doc/ai/                      Projektkontext, Architektur, Roadmap und Status
examples/mcp/                MCP-Konfigurationsbeispiele
```

## Geplante Funktionalität

- Client-Erkennung und Setup für Codex, Claude Code und Antigravity
- Installation und Synchronisierung von Skills und Agents
- Sichere Registrierung und Verwaltung von MCP-Servern
- Generierung und Synchronisierung von `AGENTS.md`, `CLAUDE.md`, `GEMINI.md` und `SKILL.md`
- Deklarative Konfiguration (`agentstrap.yaml`), Dry-Run, Backups und Konflikterkennung

Details: [Projektbeschreibung](doc/ai/PROJECT.md) und [Roadmap](doc/ai/ROADMAP.md). Client-spezifische Dateien werden nicht automatisch angelegt oder verändert.

## Sicherheit

Keine API-Schlüssel oder Tokens ins Repository committen. Spätere Änderungen an Client-Konfigurationen müssen explizit erfolgen und Backups/Preview berücksichtigen. Die vorhandenen Skills und MCP-Beispiele sind Dokumentationsressourcen, **keine** automatisch installierten Komponenten.
