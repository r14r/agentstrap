# Roadmap

Ein Roadmap-Punkt ist eine Planung, **keine implementierte Funktion**.

## v0.2.0 — Go-Grundgerüst

- [x] Go-Modul und ausführbare CLI
- [x] `help` und `version` samt Tests
- [x] `VERSION`, `justfile`, `update-cli.yaml`
- [x] Go-CI, Versions-Smoke-Test und Agent-Template-Validator
- [x] Dokumentation von der Vorlage auf das konkrete Projekt übertragen

## Nächster Meilenstein — Client-Discovery und Konfiguration

- [ ] `agentstrap detect` für Codex, Claude Code und Antigravity
- [ ] `agentstrap doctor` für lokale Voraussetzungen
- [ ] `agentstrap init` und versioniertes `agentstrap.yaml`
- [ ] Adapter-Schnittstelle für plattformspezifische Client-Konfiguration
- [ ] Testfixtures für macOS, Linux und Windows

## Weitere Meilensteine

- [ ] `skill list/search/install/remove` mit Quellen- und Sicherheitsmodell
- [ ] `agent list/create/install` und gemeinsame Agent-Definitionen
- [ ] `mcp list/add/remove` mit Secret-Referenzen
- [ ] `generate` für `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, `SKILL.md`
- [ ] `sync --dry-run`, Diff, Backup, Rollback und Konfliktauflösung
- [ ] Distribution per GitHub Releases und optional Homebrew
