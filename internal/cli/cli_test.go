package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func execute(args ...string) (int, string, string) {
	var stdout, stderr bytes.Buffer
	code := Run(args, &stdout, &stderr, "1.2.3")
	return code, stdout.String(), stderr.String()
}

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		exitCode int
		want     string
	}{
		{"default help", nil, 0, "Usage:"},
		{"explicit help", []string{"help"}, 0, "Commands:"},
		{"help flag", []string{"--help"}, 0, "agentstrap [command]"},
		{"version", []string{"version"}, 0, "agentstrap 1.2.3"},
		{"version flag", []string{"--version"}, 0, "agentstrap 1.2.3"},
		{"unknown", []string{"unsupported"}, 2, "Unknown command"},
		{"unexpected", []string{"version", "extra"}, 2, "Unexpected arguments"},
		{"skills list", []string{"skill", "list", "--topic", "development"}, 0, "go-development"},
		{"skills search", []string{"skill", "search", "video", "--topic", "video editing"}, 0, "video-editing"},
		{"agents", []string{"agent", "list", "--topic", "animation"}, 0, "animation-planner"},
		{"mcp", []string{"mcp", "search", "browser"}, 0, "playwright"},
		{"no matches", []string{"skill", "search", "unfindable"}, 0, "No matches"},
		{"invalid flag", []string{"mcp", "list", "--unknown"}, 1, "unknown option"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := execute(tt.args...)
			if code != tt.exitCode {
				t.Fatalf("exit = %d, want %d; stderr=%s", code, tt.exitCode, stderr)
			}
			if !strings.Contains(stdout+stderr, tt.want) {
				t.Fatalf("output = %q/%q, expected %q", stdout, stderr, tt.want)
			}
		})
	}
}

func TestSkillInstallPreviewAndApply(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, ".claude", "skills", "video-editing", "SKILL.md")
	args := []string{"skill", "install", "video-editing", "--path", root, "--client", "claude"}
	if code, stdout, stderr := execute(args...); code != 0 || !strings.Contains(stdout, "PREVIEW") {
		t.Fatalf("preview: %d %s %s", code, stdout, stderr)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("preview created file: %v", err)
	}
	args = append(args, "--apply")
	if code, _, stderr := execute(args...); code != 0 {
		t.Fatalf("apply: %s", stderr)
	}
	content, err := os.ReadFile(target)
	if err != nil || !strings.Contains(string(content), "name: video-editing") {
		t.Fatalf("content = %s; %v", content, err)
	}
	if code, _, _ := execute(args...); code == 0 {
		t.Fatal("overwrite unexpectedly succeeded")
	}
	if content, _ = os.ReadFile(target); !strings.Contains(string(content), "name: video-editing") {
		t.Fatal("existing skill changed")
	}
}

func TestAgentInstallationAndUnsupportedClient(t *testing.T) {
	root := t.TempDir()
	if code, _, err := execute("agent", "install", "qa-reviewer", "--client", "codex", "--path", root, "--apply"); code == 0 || !strings.Contains(err, "does not use") {
		t.Fatalf("expected explicit unsupported format, got %d %s", code, err)
	}
	if code, _, err := execute("agent", "install", "qa-reviewer", "--client", "gemini", "--path", root, "--apply"); code != 0 {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(root, ".gemini", "agents", "qa-reviewer.md"))
	if err != nil || !strings.Contains(string(content), "kind: local") {
		t.Fatalf("agent content = %s: %v", content, err)
	}
}

func TestMCPInstallPreservesConfigAndBackup(t *testing.T) {
	root := t.TempDir()
	settings := filepath.Join(root, ".gemini", "settings.json")
	if err := os.MkdirAll(filepath.Dir(settings), 0755); err != nil {
		t.Fatal(err)
	}
	original := []byte("{\"theme\":\"dark\",\"mcpServers\":{\"custom\":{\"command\":\"safe\"}}}\n")
	if err := os.WriteFile(settings, original, 0600); err != nil {
		t.Fatal(err)
	}
	args := []string{"mcp", "install", "filesystem", "--client", "gemini", "--path", root}
	if code, _, err := execute(args...); code != 0 {
		t.Fatal(err)
	}
	content, _ := os.ReadFile(settings)
	if string(content) != string(original) {
		t.Fatal("preview changed settings")
	}
	if code, _, err := execute(append(args, "--apply")...); code != 0 {
		t.Fatal(err)
	}
	var doc struct {
		Theme   string         `json:"theme"`
		Servers map[string]any `json:"mcpServers"`
	}
	content, _ = os.ReadFile(settings)
	if err := json.Unmarshal(content, &doc); err != nil {
		t.Fatal(err)
	}
	if doc.Theme != "dark" || len(doc.Servers) != 2 {
		t.Fatalf("settings lost: %+v", doc)
	}
	backup, _ := os.ReadFile(settings + ".agentstrap.bak")
	if string(backup) != string(original) {
		t.Fatal("backup does not match original")
	}
	if code, _, _ := execute(append(args, "--apply")...); code == 0 {
		t.Fatal("duplicate MCP unexpectedly succeeded")
	}
}

func TestCodexMCPAndRejectMalformedJSON(t *testing.T) {
	root := t.TempDir()
	if code, _, err := execute("mcp", "install", "playwright", "--client", "codex", "--path", root, "--apply"); code != 0 {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(root, ".codex", "config.toml"))
	if err != nil || !strings.Contains(string(content), "[mcp_servers.playwright]") {
		t.Fatalf("toml: %s %v", content, err)
	}
	if code, _, _ := execute("mcp", "install", "playwright", "--client", "codex", "--path", root, "--apply"); code == 0 {
		t.Fatal("duplicate TOML server allowed")
	}
	path := filepath.Join(root, ".mcp.json")
	if err := os.WriteFile(path, []byte("{broken"), 0600); err != nil {
		t.Fatal(err)
	}
	if code, _, err := execute("mcp", "install", "filesystem", "--client", "claude", "--path", root, "--apply"); code == 0 || !strings.Contains(err, "invalid JSON") {
		t.Fatalf("invalid config not rejected: %d %s", code, err)
	}
}

func TestInspectRepository(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module test\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "services", "api"), 0755); err != nil {
		t.Fatal(err)
	}
	code, output, errout := execute("inspect", "--path", root)
	if code != 0 {
		t.Fatal(errout)
	}
	for _, fragment := range []string{"Go", "AGENTS.md", "CLAUDE.md", "GEMINI.md", "SKILL.md", "services/api/AGENTS.md"} {
		if !strings.Contains(output, fragment) {
			t.Errorf("missing %q in %s", fragment, output)
		}
	}
	if code, output, errout := execute("repo", "inspect", "--path", root, "--client", "codex"); code != 0 || !strings.Contains(output, "AGENTS.md") || strings.Contains(output, "GEMINI.md") {
		t.Fatalf("codex inspection: %d %s %s", code, output, errout)
	}
}
