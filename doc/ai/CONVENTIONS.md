# Entwicklungs- und Dokumentationskonventionen

## Go

- Go 1.23+; Modul `github.com/r14r/agentstrap`.
- Einstieg in `cmd/agentstrap`, Paketlogik und Tests in `internal/`.
- `gofmt` ist verpflichtend, `go vet ./...` und `go test` sind Teil von CI.
- Öffentliche CLI-Befehle brauchen Tests für Erfolg und Fehlerfälle.
- Standardbibliothek bevorzugen; externe Abhängigkeiten begründen und fixieren.

## Build und Versionierung

- `VERSION` ist die kanonische SemVer-Datei.
- `just ci` enthält alle Precommit-Schritte und Race-Tests.
- `just build` enthält `just ci` und einen Binary-Smoke-Test.
- `update-cli.yaml` ruft dieselben `just`-Rezepte auf, um doppelte Prüfketten zu vermeiden.
- `CHANGELOG.md` zu jeder Versionsänderung aktualisieren.

## Sicherheit

- Keine Secrets in Git, CI-Ausgaben, Logs oder generierten Dateien.
- Agent-Skills und MCP-Definitionen können externe Befehle ausführen und müssen vor der Aktivierung geprüft werden.
- Lokale Client-Dateien später nur nach ausdrücklichem Apply ändern; Dry-Run und Backups priorisieren.

## Zusammenarbeit

- Dokumentation und fachliche Texte auf Deutsch; Quellcode und Commit-Nachrichten auf Englisch.
- Feature-Branch und PR statt Direktänderung an `main`.
- Vor Abschluss Tests und PR-Checks prüfen; ohne Freigabe keinen Merge.
