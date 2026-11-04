#!/usr/bin/env python3
"""Validate the framework-neutral AI project template (stdlib only)."""

from __future__ import annotations

import json
from pathlib import Path
import re
import sys
import tomllib

ROOT = Path(__file__).resolve().parents[1]
REQUIRED = (
    "README.md",
    "AGENTS.md",
    "VERSION",
    "CHANGELOG.md",
    "justfile",
    "doc/ai/PROJECT.md",
    "doc/ai/CONTEXT.md",
    "doc/ai/ARCHITECTURE.md",
    "doc/ai/CONVENTIONS.md",
    "doc/ai/STATUS.md",
    "doc/ai/ROADMAP.md",
    "doc/ai/CHATGPT_PROJECT.md",
    "doc/ai/TOOLS.md",
    "doc/ai/decisions/0000-template.md",
    "examples/mcp/codex.config.toml",
    "examples/mcp/vscode.mcp.json",
    ".agents/skills/implement-feature/SKILL.md",
    ".agents/skills/code-review/SKILL.md",
    ".agents/skills/release/SKILL.md",
    ".agents/skills/documentation/SKILL.md",
    ".github/workflows/ci.yml",
)
ERRORS: list[str] = []


def check(condition: bool, message: str) -> None:
    if not condition:
        ERRORS.append(message)


for name in REQUIRED:
    p = ROOT / name
    check(p.is_file(), f"Missing: {name}")
    if p.is_file():
        check(bool(p.read_text(encoding="utf-8").strip()), f"Empty: {name}")

version_file = ROOT / "VERSION"
if version_file.exists():
    version = version_file.read_text(encoding="utf-8").strip()
    check(bool(re.fullmatch(r"(?:0|[1-9]\d*)\.(?:0|[1-9]\d*)\.(?:0|[1-9]\d*)(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?", version)), f"Invalid SemVer: {version}")

for skill in sorted((ROOT / ".agents" / "skills").glob("*/SKILL.md")):
    content = skill.read_text(encoding="utf-8")
    parts = content.split("---", 2)
    check(len(parts) == 3 and not parts[0].strip(), f"Missing YAML frontmatter: {skill}")
    if len(parts) == 3:
        metadata = parts[1]
        check(bool(re.search(r"(?m)^name:\s*[a-z0-9-]+\s*$", metadata)), f"Missing skill name: {skill}")
        check(bool(re.search(r"(?m)^description:\s*\S", metadata)), f"Missing description: {skill}")
        match = re.search(r"(?m)^name:\s*([a-z0-9-]+)\s*$", metadata)
        if match:
            check(match.group(1) == skill.parent.name, f"Skill name mismatch: {skill}")

mcp_json = ROOT / "examples/mcp/vscode.mcp.json"
if mcp_json.exists():
    with mcp_json.open(encoding="utf-8") as stream:
        check(isinstance(json.load(stream).get("servers"), dict), "Missing MCP server mapping")

mcp_toml = ROOT / "examples/mcp/codex.config.toml"
if mcp_toml.exists():
    with mcp_toml.open("rb") as stream:
        check(isinstance(tomllib.load(stream).get("mcp_servers"), dict), "Missing Codex MCP server mapping")

if ERRORS:
    print("Template validation FAILED")
    for error in ERRORS:
        print(f"  - {error}")
    sys.exit(1)
print(f"Template validation OK: {len(REQUIRED)} required files checked")
