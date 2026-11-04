package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/r14r/agentstrap/src/sources"
)

func TestSourceLifecycle(t *testing.T) {
	t.Setenv("AGENTSTRAP_HOME", t.TempDir())
	code, out, stderr := execute("source", "list")
	if code != 0 || !strings.Contains(out, "anthropic-skills") || !strings.Contains(out, "claude-skills") {
		t.Fatalf("%d %s %s", code, out, stderr)
	}
	code, out, stderr = execute("source", "add", "demo", "https://github.com/example/repo", "--kind", "skill,agent", "--root", "skills")
	if code != 0 {
		t.Fatalf("add: %s %s", out, stderr)
	}
	code, _, _ = execute("source", "add", "demo", "https://github.com/example/repo")
	if code == 0 {
		t.Fatal("duplicate source accepted")
	}
	code, out, _ = execute("source", "show", "demo")
	if code != 0 || !strings.Contains(out, "copy-directory") {
		t.Fatalf("show: %s", out)
	}
	code, _, stderr = execute("source", "add", "bad", "https://attacker.example/repo")
	if code == 0 || !strings.Contains(stderr, "github.com") {
		t.Fatalf("unsafe source accepted: %s", stderr)
	}
	code, out, stderr = execute("source", "remove", "demo")
	if code != 0 {
		t.Fatalf("remove: %s %s", out, stderr)
	}
	code, _, _ = execute("source", "show", "demo")
	if code == 0 {
		t.Fatal("removed source still registered")
	}
}

func TestSourceSearchAndPreview(t *testing.T) {
	t.Setenv("AGENTSTRAP_HOME", t.TempDir())
	path, err := sources.DatabasePath()
	if err != nil {
		t.Fatal(err)
	}
	db, err := sources.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	s, err := db.Find("anthropic-skills")
	if err != nil {
		t.Fatal(err)
	}
	s.Revision = strings.Repeat("a", 40)
	s.Items = []sources.Entry{{Kind: "skill", Name: "frontend-design", Path: "skills/frontend-design", Files: []sources.File{{Path: "SKILL.md", SHA: strings.Repeat("b", 40), Size: 20}}}}
	if err := sources.Save(path, db); err != nil {
		t.Fatal(err)
	}
	code, out, stderr := execute("source", "search", "frontend", "--kind", "skill")
	if code != 0 || !strings.Contains(out, "frontend-design") {
		t.Fatalf("search: %s %s", out, stderr)
	}
	project := t.TempDir()
	code, out, stderr = execute("skill", "install", "anthropic-skills:frontend-design", "--client", "claude", "--path", project)
	if code != 0 || !strings.Contains(out, "PREVIEW") || !strings.Contains(out, ".claude/skills/frontend-design") {
		t.Fatalf("preview: %s %s", out, stderr)
	}
	if _, err := os.Stat(filepath.Join(project, ".claude")); !os.IsNotExist(err) {
		t.Fatalf("preview created files: %v", err)
	}
	code, _, stderr = execute("source", "install", "anthropic-skills:frontend-design", "--client", "bad", "--path", project)
	if code == 0 || !strings.Contains(stderr, "unsupported client") {
		t.Fatalf("unsupported client accepted: %s", stderr)
	}
}

func TestSourceMCPRejectsUnreviewedImports(t *testing.T) {
	t.Setenv("AGENTSTRAP_HOME", t.TempDir())
	path, _ := sources.DatabasePath()
	db, err := sources.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	s, _ := db.Find("claude-skills")
	s.Revision = strings.Repeat("c", 40)
	s.Items = []sources.Entry{{Kind: "mcp", Name: "root-mcp", Path: ".mcp.json", Files: []sources.File{{Path: ".mcp.json", SHA: strings.Repeat("d", 40), Size: 100}}}}
	if err := sources.Save(path, db); err != nil {
		t.Fatal(err)
	}
	code, _, stderr := execute("source", "install", "claude-skills:root-mcp", "--client", "claude", "--path", t.TempDir(), "--apply")
	if code == 0 || !strings.Contains(stderr, "manual review") {
		t.Fatalf("dangerous MCP import accepted: %s", stderr)
	}
}
