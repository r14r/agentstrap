# Entwicklungs- und Dokumentationskonventionen

## Grundsätze
- Konsistente, nachvollziehbare Änderungen statt großer ungetesteter Umbauten.
- Bestehende Projektstruktur und bewährte Werkzeuge bevorzugen.
- Schnittstellenänderungen mit Migration/Kompatibilitätshinweisen dokumentieren.
- Fehlerfälle und Sicherheitsgrenzen explizit prüfen.

## Benennung und Sprache
- Dokumentation und fachliche Ausgaben: Deutsch, wenn nichts anderes festgelegt wird.
- Code-Symbole, Skript-Namen, Commit-Nachrichten: Englisch.
- Ordner und Dateinamen möglichst aussagekräftig, kleingeschrieben und stabil.

## Testing
- Jede Verhaltensänderung durch passende Tests absichern.
- In Reviews benennen, was geprüft wurde und was nicht.
- Lokale Validierung: `just validate`; CI muss denselben Validator ausführen.

## Git
- Feature-Branches und Pull Requests für nachvollziehbare Änderungen.
- Commit-Scopes bei Bedarf nutzen, z. B. `docs:`, `feat:`, `fix:`, `chore:`.
- SemVer über `VERSION` und `CHANGELOG.md` pflegen.

## KI-spezifisch
- Keine Secrets im Prompt, Log oder Git.
- Befehle mit Seiteneffekten nur bei berechtigtem Auftrag ausführen.
- Bei unklarem Arbeitsstand erst Branch/HEAD/Dateien verifizieren.
- Änderungen an Kontext und Skills wie Produktionscode reviewen.
