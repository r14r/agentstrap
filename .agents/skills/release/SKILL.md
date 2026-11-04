---
name: release
description: Use when preparing, checking or proposing a semantic-versioned project release.
---

# Release

1. Aktuelle `VERSION`, `CHANGELOG.md`, Commits und Release-Regeln prüfen.
2. SemVer (major/minor/patch) aus tatsächlich enthaltenen Änderungen bestimmen.
3. Versionsänderung und Changelog mit nachvollziehbaren Hinweisen vorbereiten.
4. `just validate` und erforderliche Projekttests/CI ausführen.
5. Versionskonsistenz und Dokumentation überprüfen.
6. Release-Änderungen als PR zur Freigabe bereitstellen.
7. Prüfergebnisse und verbleibende Risiken berichten.

Kein Tag, Release, Merge oder Deployment ohne ausdrückliche Berechtigung.
