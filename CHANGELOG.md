# Changelog

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
