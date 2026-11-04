---
name: implement-feature
description: Use when a requested feature requires changes to code or behavior, including implementation, tests and documentation.
---

# Implement Feature

## Ablauf
1. `AGENTS.md` und passende Abschnitte aus `doc/ai/CONTEXT.md` lesen.
2. Branch, HEAD, betroffene Dateien, Tests und Abhängigkeiten prüfen.
3. Ziel, Akzeptanzkriterien und Risiken kurz festhalten.
4. Einen kleinen, überprüfbaren Implementierungsplan erstellen.
5. Änderungen inklusive aussagekräftiger Tests durchführen.
6. `just validate` und alle projektspezifischen Tests ausführen.
7. Bei Änderungen von Verhalten, Struktur oder Scope relevante Kontextdateien aktualisieren.
8. Diff und offene Risiken prüfen; Änderung über PR vorschlagen.

## Grenzen
- Keine Infrastruktur, Bibliotheken oder Frameworks ohne sachlichen Grund hinzufügen.
- Kein Merge oder Deployment ohne ausdrücklichen Auftrag.
- Tests niemals als erfolgreich deklarieren, wenn sie nicht ausgeführt wurden.
