# Roadmap

Ein Roadmap-Punkt ist eine Planung, **keine implementierte Funktion**.

## v0.2.0 — Go-Grundgerüst

- [x] Go-Modul und ausführbare CLI
- [x] `help` und `version` samt Tests
- [x] `VERSION`, `justfile`, `update-cli.yaml`
- [x] Go-CI, Versions-Smoke-Test und Agent-Template-Validator
- [x] Dokumentation von der Vorlage auf das konkrete Projekt übertragen

## v0.3.0 — Katalog, Installation und Projektinspektion

- [x] Lokaler Starterkatalog mit Skills, Agents und MCPs nach Thema
- [x] `skill|agent|mcp list|search` mit thematischen Filtern
- [x] Projektinstallation von Skills für Codex, Claude und Gemini
- [x] Projektinstallation von Agents für Claude und Gemini
- [x] Projektbezogene MCP-Konfiguration für Codex, Claude und Gemini
- [x] `inspect`/`repo inspect`: vorhandene Dateien und empfohlene Ablageorte
- [x] Vorschau als Standard; nur `--apply` schreibt Dateien
- [x] Unit-Tests für Katalog, Installation, Konflikte und Inspektion

## Nächster Meilenstein — Client-Discovery und Konfiguration

- [ ] `agentstrap detect` für Codex, Claude Code und Antigravity
- [ ] `agentstrap doctor` für lokale Voraussetzungen
- [ ] `agentstrap init` und versioniertes `agentstrap.yaml`
- [ ] Adapter-Schnittstelle für plattformspezifische Client-Konfiguration
- [ ] Testfixtures für macOS, Linux und Windows

## Weitere Meilensteine

- [ ] Externe Skill-Registries, lokale Quellen, Quellenvalidierung, Updates und Remove
- [ ] Erweiterte Agent-Definitionen und client-spezifische Konvertierung
- [ ] MCP-Quellen, Remote-Server, Version-Pinning, Secrets und Remove
- [ ] `generate` für `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, `SKILL.md`
- [ ] `sync --dry-run`, Diff, Backup, Rollback und Konfliktauflösung
- [ ] Distribution per GitHub Releases und optional Homebrew
