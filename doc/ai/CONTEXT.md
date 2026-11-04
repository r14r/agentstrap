# Kontextindex

## Lesestrategie
1. `AGENTS.md` und diese Datei zur Orientierung lesen.
2. Zweck und Scope nur bei fachlichen Fragen: `PROJECT.md`.
3. Bei Codeänderungen nur die betroffenen Dateien, Tests und relevanten Konventionen öffnen.
4. Den tatsächlichen Stand aus Git und Tests ermitteln; `STATUS.md` ist ein dokumentierter Snapshot.
5. Bei Architekturfragen `ARCHITECTURE.md` und betroffene ADRs hinzuziehen.
6. Bei Planung `ROADMAP.md`; bei externen Werkzeugen `TOOLS.md`.
7. Skills nur bei passender Aufgabe laden.

## Quellenhierarchie
1. Aktuell geprüfter Repository-Stand, Quellcode und Tests
2. Normative Projektdokumente: `AGENTS.md`, `PROJECT.md`, ADRs
3. Aktualisierte Architektur- und Konventionsdokumente
4. `STATUS.md` und `ROADMAP.md` (Snapshots und Planungen)
5. Chatverläufe (Hinweise, keine automatisch verifizierte Wahrheit)

## Kontext pro Aufgabentyp
| Aufgabe | Primäre Quellen |
| --- | --- |
| Neues Feature | PROJECT, ARCHITECTURE, CONVENTIONS, betroffener Code |
| Bugfix | Fehlermeldung, reproduzierbarer Test, betroffener Code |
| Review | Diff, Tests, CONVENTIONS, ADRs |
| Release | VERSION, CHANGELOG, STATUS, CI |
| Planung | PROJECT, STATUS, ROADMAP |
| Tool/MCP | TOOLS und clientabhängige Konfigurationen |

## GitHub und ChatGPT
- ChatGPT-Projektanweisungen: `CHATGPT_PROJECT.md` als Kopiervorlage.
- Die Existenz von `AGENTS.md` im Repo bindet diese Datei nicht automatisch in jeden Chat ein.
- Hochgeladene Kopien müssen nach Änderungen neu bereitgestellt werden; für aktuellen Stand GitHub abrufen.
