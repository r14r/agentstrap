# Changelog

## 0.4.2 — 2026-11-04

- Move Go packages and tests from `internal/` to `src/`, update Go imports to `github.com/r14r/agentstrap/src/...`.
- Move developer scripts from `scripts/` to `bin/` while preserving executable permissions.
- Update Justfile, Python structure validator, architecture docs and CI references; reject legacy folders.

## 0.4.1 — 2026-11-04

- Added `just` shortcuts for source list/show/add/remove/sync/search/install preview/install apply.
- Included no-argument Anthropic frontend-design examples.
- Added offline CI smoke tests for the just-based source lifecycle.

## 0.4.0 — 2026-11-04

- Local, persistent JSON source registry seeded with Anthropic Skills and alirezarezvani/claude-skills.
- CLI `source list|show|add|remove|sync|search|install` manages indexed GitHub resources.
- SHA-pinned installs for skills and agents with preview, blob verification and safe destination handling.
- External MCP definitions indexed for manual review but never auto-imported.
- Registry, URL, GitHub sync, manifest and CLI tests.

## 0.3.0 — 2026-11-04

- Lokaler, themenbasierter Katalog: Skills, Agents, MCPs für Development, Video Editing und Animation.
- Befehle `skill|agent|mcp list|search|install` inklusive Filtern.
- Projektbezogene Installation für Codex, Claude Code und Gemini CLI mit Dry-Run als Standard.
- Sicheres Hinzufügen ausgewählter MCP-Konfigurationen mit Backup bestehender Einstellungen.
- `inspect` und `repo inspect` erkennen Projektmanifeste und schlagen passende Instruction-/Skill-Dateien vor.
- Tests für Katalog, Preview, Schreibschutz, client-spezifische Konfiguration und Inspektion.

## 0.2.1 — 2026-11-04

- `just install` installs the CLI to `/usr/local/bin` by default, escalating to `sudo` only when needed.
- Preserve an optional `INSTALL_DIR` override for CI and custom paths.
- Align and document `bootstrap`, `fmt`, `precommit`, `ci`, `build`, `install`, and `next` workflows with other Go CLI projects.
- Add `help`, `setup`, `lint`, and `version` convenience recipes.
- CI validates the canonical `just install` dependency chain with a temporary destination.

## 0.2.0 — 2026-11-04

- Go-Modul `github.com/r14r/agentstrap` und ausführbarer CLI-Einstiegspunkt.
- CLI-Befehle `help` und `version` mit Unit-Tests.
- Vollständiges `justfile` für Formatierung, Tests, CI, Build, Installation, Update und PR-Übersicht.
- `update-cli.yaml` mit Schema v2 und kanonischen `just`-Workflows.
- GitHub Actions für Go-Checks, Vorlagevalidierung, Build und Versions-Smoke-Test.
- Projektdokumentation von der AI-Projektvorlage auf AgentStrap umgestellt.

## 0.1.0 — 2026-10-23

- Initiale, frameworkneutrale KI-Projektvorlage.
- AGENTS.md, Kontext-Routing und ChatGPT-Projektanweisungen.
- Wiederverwendbare Agent Skills und optionale MCP-Beispiele.
- Template-Validierung mit Python und GitHub Actions.
