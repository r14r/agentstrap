# AI Project Template

Versionierte Vorlage für die Zusammenarbeit zwischen **ChatGPT-Projekten**, **GitHub-Repositories** und Coding-Agents (z. B. Codex und VS Code).

## Grundprinzip

Das Repository ist die **Single Source of Truth** für Anforderungen, Architektur, Projektregeln und verifizierten Arbeitsstand. Ein ChatGPT-Projekt ist die Arbeitsoberfläche; hochgeladene Projektdateien werden nicht automatisch mit GitHub synchronisiert.

| Ebene | Einstiegspunkt | Zweck |
| --- | --- | --- |
| Regeln | [AGENTS.md](AGENTS.md) | Verbindliche Agent-Anweisungen |
| Kontext | [doc/ai/CONTEXT.md](doc/ai/CONTEXT.md) | Kontext-Routing und Quellen |
| Projektwissen | [doc/ai/](doc/ai/) | Ziele, Architektur, Konventionen, Status, Roadmap |
| Skills | [.agents/skills/](.agents/skills/) | Wiederverwendbare Workflows |
| MCP | [doc/ai/TOOLS.md](doc/ai/TOOLS.md) und [examples/mcp/](examples/mcp/) | Optionale, clientabhängige Konfiguration |
| Prüfung | [scripts/validate_template.py](scripts/validate_template.py) | Template-Struktur validieren |

## Neues Projekt aufsetzen

1. Auf GitHub in den Repository-Einstellungen **Template repository** aktivieren (einmalig für dieses Repository).
2. Mit **Use this template** ein neues Repository erstellen oder die Vorlage in ein bestehendes Projekt übernehmen.
3. In `doc/ai/PROJECT.md` den Projektnamen, Zweck, Umfang und das kanonische Repository eintragen.
4. `doc/ai/CONTEXT.md`, `ARCHITECTURE.md`, `CONVENTIONS.md`, `STATUS.md` und `ROADMAP.md` für das Projekt anpassen.
5. `AGENTS.md` und Skills prüfen; nur relevante Skills behalten.
6. Optional einen MCP-Client aus `examples/mcp/` konfigurieren und geheime Werte ausschließlich über sichere lokale Umgebungsvariablen bzw. den Secret Store bereitstellen.
7. Den Inhalt von `doc/ai/CHATGPT_PROJECT.md` als Vorlage für die **Projektanweisungen in ChatGPT** verwenden. Für aktuelle Dateien GitHub gezielt lesen; ChatGPT lädt `AGENTS.md` nicht automatisch allein wegen seiner Existenz.
8. Mit `just validate` oder `python3 scripts/validate_template.py` die Struktur prüfen.

## Arbeitsweise

- **Beim Einstieg:** `AGENTS.md` und `doc/ai/CONTEXT.md` lesen.
- **Pro Aufgabe:** nur fachlich relevante Dokumente und Quellcode laden.
- **Vor Änderungen:** Branch, Commit und tatsächlichen Zustand prüfen; `STATUS.md` ist ein Snapshot.
- **Bei Umsetzung:** passende Skills verwenden, danach Tests ausführen.
- **Bei Abschluss:** relevante Dokumentation und ggf. Architekturentscheidungen aktualisieren; Änderung über Pull Request einreichen.

## Verzeichnisstruktur

```text
AGENTS.md
.agents/skills/{implement-feature,code-review,release,documentation}/SKILL.md
.github/workflows/ci.yml
doc/ai/{PROJECT,CONTEXT,ARCHITECTURE,CONVENTIONS,STATUS,ROADMAP,CHATGPT_PROJECT,TOOLS}.md
doc/ai/decisions/0000-template.md
examples/mcp/{codex.config.toml,vscode.mcp.json}
scripts/validate_template.py
justfile
VERSION
CHANGELOG.md
```

## Hinweise

- `AGENTS.md` ist die gängige Codex-Datei; `AGENT.md` ist **kein** Ersatz.
- `.agents/skills/` ist ein Agent-Skills-Format. Die automatische Skill-Erkennung unterscheidet sich zwischen Clients.
- MCP-Konfigurationen sind **client-spezifisch**, keine universelle Projektdatei.
- Keine API-Schlüssel, Zugangstokens oder anderen Secrets committen.
- Dieses Repository ist bewusst **frameworkneutral**: kein automatisch ausgewähltes Python-, FastAPI-, Nuxt- oder Docker-Gerüst.
