package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/r14r/agentstrap/internal/catalog"
)

func installItem(item catalog.Item, opts options, out io.Writer) error {
	if opts.client != "codex" && opts.client != "claude" && opts.client != "gemini" {
		return fmt.Errorf("unsupported client %q (choose codex, claude, or gemini)", opts.client)
	}
	root, err := filepath.Abs(opts.path)
	if err != nil {
		return err
	}
	info, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("project directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("project path is not a directory: %s", root)
	}

	if item.Kind == "mcp" {
		return installMCP(item, opts, root, out)
	}
	var relative, content string
	switch item.Kind {
	case "skill":
		dir := map[string]string{
			"codex":  ".agents/skills",
			"claude": ".claude/skills",
			"gemini": ".gemini/skills",
		}[opts.client]
		relative = filepath.Join(dir, item.ID, "SKILL.md")
		content = fmt.Sprintf("---\nname: %s\ndescription: %s\n---\n\n# %s\n\n%s\n", item.ID, item.Description, item.ID, item.Guide)
	case "agent":
		if opts.client == "codex" {
			return errors.New("Codex does not use the Claude/Gemini .md subagent format; use a skill instead")
		}
		relative = filepath.Join("."+opts.client, "agents", item.ID+".md")
		frontmatter := fmt.Sprintf("---\nname: %s\ndescription: %s\n", item.ID, item.Description)
		if opts.client == "gemini" {
			frontmatter += "kind: local\n"
		}
		content = frontmatter + "---\n\n" + item.Guide + "\n"
	default:
		return fmt.Errorf("unsupported resource type %s", item.Kind)
	}
	target := filepath.Join(root, relative)
	if _, err := os.Lstat(target); err == nil {
		return fmt.Errorf("target already exists (will not overwrite): %s", target)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if !opts.apply {
		fmt.Fprintf(out, "PREVIEW: would create %s\n%s\nApply with --apply.\n", target, content)
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(target, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err = io.WriteString(file, content); err != nil {
		file.Close()
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	fmt.Fprintf(out, "Installed %s\n", target)
	return nil
}

func installMCP(item catalog.Item, opts options, root string, out io.Writer) error {
	if item.Command == "" {
		return fmt.Errorf("%s has no verified launch configuration", item.ID)
	}
	var relative string
	switch opts.client {
	case "claude":
		relative = ".mcp.json"
	case "gemini":
		relative = ".gemini/settings.json"
	case "codex":
		relative = ".codex/config.toml"
	}
	target := filepath.Join(root, relative)
	old, perm, err := readOptionalFile(target)
	if err != nil {
		return err
	}
	args := append([]string{}, item.Args...)
	if item.ID == "filesystem" {
		args = append(args, root)
	}

	var next []byte
	if opts.client == "codex" {
		header := "[mcp_servers." + item.ID + "]"
		if strings.Contains("\n"+string(old), "\n"+header) {
			return fmt.Errorf("MCP %q already configured in %s", item.ID, target)
		}
		section := "\n" + header + "\ncommand = " + strconv.Quote(item.Command) + "\nargs = ["
		for i, arg := range args {
			if i != 0 {
				section += ", "
			}
			section += strconv.Quote(arg)
		}
		next = []byte(strings.TrimRight(string(old), "\n") + section + "]\n")
	} else {
		doc := map[string]any{}
		if len(old) > 0 {
			if err := json.Unmarshal(old, &doc); err != nil {
				return fmt.Errorf("%s: invalid JSON: %w", target, err)
			}
			if doc == nil {
				return fmt.Errorf("%s: expected JSON object", target)
			}
		}
		var servers map[string]any
		if existing, exists := doc["mcpServers"]; exists {
			var ok bool
			servers, ok = existing.(map[string]any)
			if !ok {
				return fmt.Errorf("%s: mcpServers must be an object", target)
			}
		} else {
			servers = map[string]any{}
			doc["mcpServers"] = servers
		}
		if _, exists := servers[item.ID]; exists {
			return fmt.Errorf("MCP %q already configured in %s", item.ID, target)
		}
		servers[item.ID] = map[string]any{"command": item.Command, "args": args}
		next, err = json.MarshalIndent(doc, "", "  ")
		if err != nil {
			return err
		}
		next = append(next, '\n')
	}
	if !opts.apply {
		fmt.Fprintf(out, "PREVIEW: would configure MCP %q in %s\nCommand: %s %s\nThis server may execute when the client loads. Review the package and access permissions.\nApply with --apply.\n", item.ID, target, item.Command, strings.Join(args, " "))
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}
	if err := replaceWithBackup(target, old, next, perm); err != nil {
		return err
	}
	fmt.Fprintf(out, "Configured MCP %q in %s (review server permissions before enabling).\n", item.ID, target)
	return nil
}

func readOptionalFile(path string) ([]byte, os.FileMode, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, 0600, nil
	}
	if err != nil {
		return nil, 0, err
	}
	if !info.Mode().IsRegular() {
		return nil, 0, fmt.Errorf("refusing non-regular config file: %s", path)
	}
	data, err := os.ReadFile(path)
	return data, info.Mode().Perm(), err
}

// Existing configuration files are backed up once and replaced atomically.
func replaceWithBackup(path string, old, next []byte, perm os.FileMode) error {
	if _, err := os.Lstat(path); err == nil {
		backup := path + ".agentstrap.bak"
		f, err := os.OpenFile(backup, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return fmt.Errorf("cannot create backup %s: %w", backup, err)
		}
		if _, err = f.Write(old); err != nil {
			f.Close()
			return err
		}
		if err = f.Close(); err != nil {
			return err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".agentstrap-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if err = f.Chmod(perm); err != nil {
		f.Close()
		return err
	}
	if _, err = f.Write(next); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}
