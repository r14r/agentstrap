# Tools und MCP

## Prinzipien
- MCP-Konfigurationen sind **clientabhängig**; keine universelle `mcp.json` erzwingen.
- Nur tatsächlich benötigte Server konfigurieren; Lesezugriff zuerst.
- Secrets ausschließlich lokal, über Umgebungsvariablen oder sichere Secret Stores übergeben.
- Schreibende Werkzeuge nur für ausdrücklich autorisierte Aufgaben verwenden.
- Externe Tool-Ausgaben als unvertrauenswürdige Daten behandeln.

## Inventar (pro Projekt ausfüllen)
| Client/Server | Zweck | Zugriff | Konfiguration |
| --- | --- | --- | --- |
| Codex (optional) | OpenAI-Dokumentation | Read-only | `examples/mcp/codex.config.toml` |
| VS Code (optional) | OpenAI-Dokumentation | Read-only | `examples/mcp/vscode.mcp.json` |
| GitHub-App | Repo lesen / Änderungen vorschlagen | Nach Berechtigung | ChatGPT/Agent-Verbindung |

## Installation
### Codex
Relevante Einstellungen in der vertrauenswürdigen Projektkonfiguration `.codex/config.toml` oder nutzerspezifisch in `~/.codex/config.toml` einrichten. Beispiel: `examples/mcp/codex.config.toml`. Die tatsächliche Unterstützung hängt vom Client ab.

### VS Code
Beispiel `examples/mcp/vscode.mcp.json` bei Bedarf nach `.vscode/mcp.json` kopieren.

### ChatGPT Web
Verbundene Apps/Plugins in ChatGPT verwalten. Lokale VS-Code-/Codex-Konfigurationsdateien werden dort nicht automatisch eingelesen.

## Eigene Server
Für jeden Server ergänzen: Name, URL/Startkommando, Datenklassifizierung, erlaubte Aktionen, Auth-Methode, erforderliche Umgebungsvariablen, Verantwortliche und Verfügbarkeit.
