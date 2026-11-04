# Projektsteckbrief

## Identität

- **Name:** AgentStrap
- **Kanonisches Repository:** `r14r/agentstrap`
- **CLI-Befehl:** `agentstrap`
- **Leitmotiv:** One CLI. Every AI coding environment.
- **Technologie:** Go (Standardbibliothek als Ausgangspunkt)

## Ziel und Nutzen

Eine konsistente, reproduzierbare Einrichtung von AI-Coding-Umgebungen über mehrere Clients hinweg. Skills, Agents, MCP-Server und projektbezogene Agent-Anweisungen sollen langfristig aus einer zentralen Konfiguration verwaltet werden.

## Zielgruppen

Softwareentwickler, Architekten und Teams, die Codex, Claude Code, Antigravity oder weitere AI-Clients parallel einsetzen.

## Scope

- Erkennung und Vorbereitung unterstützter AI-Clients.
- Verwaltung von Skills, Agents und MCP-Servern.
- Generierung und Synchronisierung von Agent-Dateien (`AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, `SKILL.md`).
- Konfigurationsvalidierung, Vorschau und sichere Änderungen.
- Lokale, reproduzierbare Go-CLI mit Build- und Testworkflow.

## Nicht-Ziele der Bootstrap-Phase

- Noch keine tatsächliche Installation von AI-Clients.
- Noch keine Konfigurationsänderungen am Benutzerverzeichnis.
- Keine Cloud-Anbindung oder zentrale Secret-Speicherung.

## Qualitätsziele

- macOS, Linux und perspektivisch Windows.
- Sichere, explizite Änderungen statt stiller Überschreibungen.
- Minimaler Abhängigkeitsumfang, überprüfbare Tests.
- Stabiler Go-Build, kanonische Versionierung und reproduzierbare CLI-Aufrufe.
