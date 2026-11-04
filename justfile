set shell := ["bash", "-cu"]

default:
    @just --list

# Vorlage ohne Drittanbieter-Abhängigkeiten prüfen
validate:
    python3 scripts/validate_template.py

# CI-kompatibler Alias
ci: validate

test: validate
