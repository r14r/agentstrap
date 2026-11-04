# AgentStrap

**One CLI. Every AI coding environment.**

AgentStrap ist eine in Go geschriebene CLI zur zukünftigen einheitlichen Einrichtung von AI-Coding-Clients (Codex, Claude Code, Antigravity), Skills, Agents, MCP-Servern und projektspezifischen Agent-Konfigurationsdateien.

**Stand v0.2.1:** Das Repository enthält eine lauffähige, getestete CLI-Basis mit `help` und `version`. Client-/Skill-/Agent-/MCP-Verwaltung und Dateigenerierung sind **noch nicht implementiert**; siehe [Roadmap](doc/ai/ROADMAP.md).

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

### Systemweite Installation

```sh
just install
/usr/local/bin/agentstrap version
```

`just install` baut das Binary, kopiert es nach **`/usr/local/bin/agentstrap`** und überprüft die installierte Versionsnummer. Nur wenn die erforderlichen Schreibrechte fehlen, wird `sudo` für die Installation verwendet.

Für ein alternatives Verzeichnis (z. B. zur Installation ohne `sudo` oder für CI):

```sh
INSTALL_DIR="$HOME/.local/bin" just install
```

## Entwicklungsbefehle

| Befehl | Wirkung |
| --- | --- |
| `just` / `just help` | Alle verfügbaren Rezepte anzeigen |
| `just bootstrap` | Go-Module vorbereiten |
| `just fmt` | Go-Code formatieren |
| `just fmt-check` | Formatierung ohne Änderungen prüfen |
| `just lint` | Formatierung und `go vet` prüfen |
| `just test` | Go-Unit-Tests |
| `just validate` | Projektstruktur und Agent-Skills prüfen |
| `just precommit` | Bootstrap, Lint, Tests und Projektvalidierung |
| `just ci` | Precommit einschließlich Race-Tests |
| `just build` | CI, Binary bauen, Version und Hilfe prüfen |
| `just install` | Build + Installation nach `/usr/local/bin`, Versionsprüfung |
| `just run --help` | CLI direkt über `go run` ausführen |
| `just next` | Git-Stand, offene PRs und nächste Prüfschritte anzeigen |
| `just next --run` | Zusätzlich `just build` und ggf. PR-Checks ausführen |
| `just update` | `git pull --ff-only`, danach `just build` |
| `just version` | Aktuelle Version anzeigen |
| `just clean` | `dist/` entfernen |

**Prüfkette:** `just install` → `just build` → `just ci` → `just precommit`. Ein separates `just test` oder `just precommit` vor `just build` ist nicht nötig. `just fmt` formatiert bewusst nur auf ausdrücklichen Aufruf; CI prüft die Formatierung ohne Änderungen.

## Update CLI

Die Datei [update-cli.yaml](update-cli.yaml) verwendet **schemaVersion 2**. Sie definiert die Workflows `bootstrap`, `setup`, `precommit`, `test`, `ci`, `build`, `install`, `run` und `clean`. Alle Workflows verweisen auf die kanonischen `just`-Rezepte.

Die Konfiguration verweist auf das GitHub-Repository `r14r/agentstrap` auf `main`. Die Standardinstallation erfolgt nach `/usr/local/bin` (bei Bedarf über `sudo`).

## Struktur

```text
cmd/agentstrap/main.go       CLI-Einstiegspunkt und eingebettete Build-Version
internal/cli/                Kommandoverarbeitung und Tests
go.mod                       Go-Modul
VERSION                      SemVer-Version (Single Source of Truth)
justfile                     Build-, Prüf-, Installations- und Update-Rezepte
update-cli.yaml              Schema-v2-Workflows für Update CLI
.github/workflows/ci.yml     Automatische Tests des just-Build- und Installationswegs
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
