package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func inspectCommand(args []string, out io.Writer) error {
	opts, positional, err := parseOptions(args)
	if err != nil {
		return err
	}
	if len(positional) != 0 || opts.apply || opts.topic != "" {
		return fmt.Errorf("usage: agentstrap inspect [--path DIR] [--client codex|claude|gemini|all]")
	}
	if opts.client == "" {
		opts.client = "all"
	}
	if opts.client != "all" && opts.client != "codex" && opts.client != "claude" && opts.client != "gemini" {
		return fmt.Errorf("unknown client %q", opts.client)
	}
	root, err := filepath.Abs(opts.path)
	if err != nil {
		return err
	}
	info, err := os.Stat(root)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", root)
	}
	fmt.Fprintf(out, "Inspecting: %s\n", root)
	signals := map[string][]string{
		"go.mod":             {"Go"},
		"pyproject.toml":     {"Python"},
		"requirements.txt":   {"Python"},
		"package.json":       {"Node/JavaScript"},
		"Cargo.toml":         {"Rust"},
		"docker-compose.yml": {"Docker Compose"},
		"compose.yaml":       {"Docker Compose"},
		"blender.blend":      {"Animation/Blender"},
	}
	detected := map[string]bool{}
	for manifest, tags := range signals {
		if isRegular(filepath.Join(root, manifest)) {
			for _, tag := range tags {
				detected[tag] = true
			}
		}
	}
	var labels []string
	for s := range detected {
		labels = append(labels, s)
	}
	sort.Strings(labels)
	if len(labels) == 0 {
		labels = append(labels, "no recognized root manifest")
	}
	fmt.Fprintf(out, "Detected: %s\n\n", strings.Join(labels, ", "))
	if opts.client == "all" || opts.client == "codex" {
		suggest(out, root, "AGENTS.md", "Repository-wide Codex instructions: architecture, build, test and review rules")
		suggest(out, root, ".agents/skills/<name>/SKILL.md", "Reusable, task-specific Codex skill; use a real skill name")
	}
	if opts.client == "all" || opts.client == "claude" {
		suggest(out, root, "CLAUDE.md", "Claude project instructions; do not duplicate large AGENTS.md files")
		suggest(out, root, ".claude/skills/<name>/SKILL.md", "Reusable Claude skill")
		suggest(out, root, ".claude/agents/<name>.md", "Optional Claude specialist agent")
	}
	if opts.client == "all" || opts.client == "gemini" {
		suggest(out, root, "GEMINI.md", "Gemini project instructions and conventions")
		suggest(out, root, ".gemini/skills/<name>/SKILL.md", "Reusable Gemini skill")
		suggest(out, root, ".gemini/agents/<name>.md", "Optional Gemini specialist agent")
	}
	suggest(out, root, "README.md", "Developer setup and project overview")
	for _, dir := range []string{"apps", "packages", "services"} {
		children, err := os.ReadDir(filepath.Join(root, dir))
		if err != nil {
			continue
		}
		count := 0
		for _, child := range children {
			if !child.IsDir() || count >= 8 {
				continue
			}
			sub := filepath.Join(dir, child.Name())
			if opts.client == "all" || opts.client == "codex" {
				suggest(out, root, filepath.Join(sub, "AGENTS.md"), "Optional specialized instructions for this subproject")
			}
			count++
		}
	}
	return nil
}

func isRegular(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func suggest(out io.Writer, root, relative, why string) {
	// Paths containing placeholders are recommendations, not literal filenames.
	state := "MISSING"
	if strings.Contains(relative, "<") {
		state = "OPTIONAL"
	} else if isRegular(filepath.Join(root, relative)) {
		state = "EXISTS"
	}
	fmt.Fprintf(out, "%-8s %-43s %s\n", state, relative, why)
}
