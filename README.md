# AgentStrap

**One CLI. Every AI coding environment.**

AgentStrap ist eine in Go geschriebene CLI zur zukünftigen einheitlichen Einrichtung von AI-Coding-Clients (Codex, Claude Code, Antigravity), Skills, Agents, MCP-Servern und projektspezifischen Agent-Konfigurationsdateien.

**Stand v0.4.1:** Das Repository enthält einen integrierten Katalog, lesende Repository-Analyse sowie projektbezogene Installationen für Skills, Agents und ausgewählte MCP-Konfigurationen. Client-Installer, externe Registry-Synchronisierung und automatische Dateigenerierung sind **noch nicht implementiert**; siehe [Roadmap](doc/ai/ROADMAP.md).

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

## Katalog, Installation und Repository-Analyse

AgentStrap v0.3.0 enthält einen **eingebauten Starterkatalog** mit Themen `development`, `video-editing` und `animation`. Die Suche ist offline und katalogbasiert; noch keine externe Skill-/MCP-Registry.

```bash
# Skills, Agents und MCPs nach Thema suchen
agentstrap skill list --topic development
agentstrap skill search "video" --topic "video editing"
agentstrap agent list --topic animation
agentstrap mcp search browser

# Projektinstallation: zunächst immer nur Vorschau
agentstrap skill install video-editing --client claude --path ./my-project
agentstrap skill install go-development --client codex --path ./my-project
agentstrap agent install animation-planner --client gemini --path ./my-project
agentstrap mcp install filesystem --client claude --path ./my-project

# Erst --apply schreibt Dateien
agentstrap skill install video-editing --client claude --path ./my-project --apply

# Repository/Ordner auf sinnvolle Konfigurationsdateien prüfen
agentstrap inspect --path ./my-project
agentstrap repo inspect --client codex --path ./my-project
```

**Installationsziele (Projektumfang):**

| Typ | Codex | Claude Code | Gemini CLI |
| --- | --- | --- | --- |
| Skill | `.agents/skills/<id>/SKILL.md` | `.claude/skills/<id>/SKILL.md` | `.gemini/skills/<id>/SKILL.md` |
| Agent | Nicht unterstützt; stattdessen Skill nutzen | `.claude/agents/<id>.md` | `.gemini/agents/<id>.md` |
| MCP | `.codex/config.toml` | `.mcp.json` | `.gemini/settings.json` |
| Projektanweisungen | `AGENTS.md` | `CLAUDE.md` | `GEMINI.md` |

Der Inspektionsbefehl prüft bestehende Dateien und leitet Hinweise aus Manifestdateien und ersten Monorepo-Unterordnern ab. Er erstellt **keine** Dateien. Skills und Agents werden beim Installieren als lokale Startervorlagen erzeugt; es werden keine Drittquellen installiert. Bestehende Skill- und Agentdateien werden nicht überschrieben.

Bei MCPs ergänzt `--apply` die Projektkonfiguration. Bereits vorhandene JSON-Felder bleiben erhalten. Vor einer Änderung bestehender Einstellungen werden versionierte `.agentstrap.bak`-Sicherheitskopien angelegt. MCP-Server werden **nicht während der Installation ausgeführt**, können aber beim späteren Laden durch den AI-Client starten. Prüfe daher Paketquelle und Berechtigungen. Der Filesystem-MCP erhält Zugriff auf den angegebenen Projektordner; Playwright startet potenziell Browserautomation. Die aktuelle Variante verwendet `npx`, weshalb eine Node-Installation nötig ist.

**Nicht enthalten:** zentrale/verteilte Registry, Remote-Suche, Download beliebiger Skills, automatische Client-Installation, Dateigenerierung oder Codex-Subagent-Unterstützung.

## External Source Registry (v0.4.0)

AgentStrap keeps a persistent, schema-versioned JSON source database at `$AGENTSTRAP_HOME/sources.json` or the OS user-config directory (`agentstrap/sources.json`). The database is seeded with two sources: [Anthropic Skills](https://github.com/anthropics/skills) (`anthropic-skills`) and [alirezarezvani/claude-skills](https://github.com/alirezarezvani/claude-skills) (`claude-skills`). Source records include repository URL, ref, supported resource types, installation method, last synced immutable commit, and indexed files.

```bash
# Manage registered GitHub repositories
agentstrap source list
agentstrap source show anthropic-skills
agentstrap source add my-skills https://github.com/owner/repo --kind skill --root skills
agentstrap source remove my-skills

# Index GitHub resources (explicit network access)
agentstrap source sync anthropic-skills
agentstrap source sync claude-skills
agentstrap source sync all

# Offline search of indexed skills, agents and MCP configuration files
agentstrap source search frontend --kind skill
agentstrap source search architect --kind agent --source claude-skills
agentstrap source search mcp --kind mcp

# Preview targeted skill installation into the selected client
agentstrap source install anthropic-skills:frontend-design --client claude --path .
agentstrap skill install anthropic-skills:frontend-design --client codex --path .

# Apply only with explicit approval
agentstrap source install anthropic-skills:frontend-design --client claude --path . --apply
```

**Install strategies:** Skill packages are copied as directories including associated files; agents as individual Markdown files. The selected GitHub revision is pinned during sync and each download is verified against its Git blob SHA. Installation previews changes by default, rejects existing targets and symlink parents, and does not execute downloaded code.

**Remote MCP configs:** Indexed and searchable, but marked `manual-review`. Importing and automatically activating arbitrary third-party `.mcp.json` configurations is deliberately disabled due to potential command execution, credentials and filesystem access. Built-in reviewed MCP definitions remain installable via `agentstrap mcp install`.

Inspect individual licenses, referenced scripts and client compatibility before applying. GitHub API rate limits apply to the sync operation.

## Just commands for external source registries

These shortcuts execute the current Go CLI directly; they do **not** rebuild or rerun CI on every search.

```bash
just source-list
just source-show anthropic-skills
just source-add my-skills https://github.com/owner/repo skill skills
just source-remove my-skills

just source-sync-all
just source-sync anthropic-skills
just source-search frontend skill
just source-search architect agent claude-skills

# Preview first, then explicitly apply
just source-install anthropic-skills:frontend-design claude .
just source-install-apply anthropic-skills:frontend-design claude .

# No-argument Anthropic frontend-design shortcuts
just source-frontend-preview
just source-frontend-apply
```

Run `source-sync` before trying remote installation; the source index is stored locally. `source-install` performs a preview, and only `source-install-apply` writes project files. `source-add`/`source-remove` change the local registry, while `source-sync` makes GitHub network requests. Arguments are escaped using `quote()` from just.

## Geplante Funktionalität

- Client-Erkennung und Setup für Codex, Claude Code und Antigravity
- Erweiterbare Quellen, Updates und Synchronisierung von Skills und Agents
- Sichere Registrierung und Verwaltung von MCP-Servern
- Generierung und Synchronisierung von `AGENTS.md`, `CLAUDE.md`, `GEMINI.md` und `SKILL.md`
- Deklarative Konfiguration (`agentstrap.yaml`), Dry-Run, Backups und Konflikterkennung

Details: [Projektbeschreibung](doc/ai/PROJECT.md) und [Roadmap](doc/ai/ROADMAP.md). Client-spezifische Dateien werden nicht automatisch angelegt oder verändert.

## Sicherheit

Keine API-Schlüssel oder Tokens ins Repository committen. Spätere Änderungen an Client-Konfigurationen müssen explizit erfolgen und Backups/Preview berücksichtigen. Die vorhandenen Skills und MCP-Beispiele sind Dokumentationsressourcen, **keine** automatisch installierten Komponenten.
