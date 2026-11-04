package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/r14r/agentstrap/internal/sources"
)

func sourceCommand(args []string, out io.Writer) error {
	if len(args) == 0 { return fmt.Errorf("usage: agentstrap source list|show|add|remove|sync|search|install") }
	dbPath, err := sources.DatabasePath()
	if err != nil { return err }
	db, err := sources.Load(dbPath)
	if err != nil { return err }
	action := args[0]
	rest := args[1:]
	switch action {
	case "list":
		if len(rest) != 0 { return fmt.Errorf("usage: source list") }
		for _, s := range db.Sources {
			fmt.Fprintf(out, "%-19s %-44s types=%s items=%d synced=%t\n", s.ID, s.URL, strings.Join(s.Kinds, ","), len(s.Items), s.Revision != "")
		}
		return nil
	case "show":
		if len(rest) != 1 { return fmt.Errorf("usage: source show ID") }
		s, err := db.Find(rest[0]); if err != nil { return err }
		fmt.Fprintf(out, "Source: %s\nName: %s\nGitHub: %s\nRef: %s\nKinds: %s\nInstall methods: %v\nRevision: %s\nSynced at: %s\nItems: %d\nNotes: %s\n", s.ID, s.Name, s.URL, s.Ref, strings.Join(s.Kinds, ","), s.InstallMethods, s.Revision, s.SyncedAt, len(s.Items), s.Notes)
		return nil
	case "add":
		if len(rest) < 2 { return fmt.Errorf("usage: source add ID URL [--ref main] [--kind skill,agent,mcp] [--root PATH]") }
		s := sources.Source{ID: rest[0], Name: rest[0], URL: rest[1], Ref: "main", Kinds: []string{"skill"}, InstallMethods: map[string]string{}, Items: []sources.Entry{}}
		for i := 2; i < len(rest); i++ {
			if i+1 >= len(rest) { return fmt.Errorf("%s requires a value", rest[i]) }
			flag, value := rest[i], rest[i+1]; i++
			switch flag {
			case "--ref": s.Ref = value
			case "--kind": s.Kinds = strings.Split(value, ",")
			case "--root": s.Roots = append(s.Roots, value)
			default: return fmt.Errorf("unknown source option %s", flag)
			}
		}
		for _, k := range s.Kinds {
			switch k {
			case "skill": s.InstallMethods[k] = "copy-directory"
			case "agent": s.InstallMethods[k] = "copy-file"
			case "mcp": s.InstallMethods[k] = "manual-review"
			}
		}
		if err := db.Add(s); err != nil { return err }
		if err := sources.Save(dbPath, db); err != nil { return err }
		fmt.Fprintf(out, "Registered source %s (not synced yet).\n", s.ID)
		return nil
	case "remove":
		if len(rest) != 1 { return fmt.Errorf("usage: source remove ID") }
		if err := db.Remove(rest[0]); err != nil { return err }
		if err := sources.Save(dbPath, db); err != nil { return err }
		fmt.Fprintf(out, "Removed source %s from local registry. Installed resources are unchanged.\n", rest[0])
		return nil
	case "sync":
		if len(rest) > 1 { return fmt.Errorf("usage: source sync [ID|all]") }
		id := "all"
		if len(rest) == 1 { id = rest[0] }
		if id != "all" {
			if _, err := db.Find(id); err != nil { return err }
		}
		api := sources.NewGitHubClient()
		for i := range db.Sources {
			s := &db.Sources[i]
			if id != "all" && s.ID != id { continue }
			fmt.Fprintf(out, "Syncing %s from %s (%s)...\n", s.ID, s.URL, s.Ref)
			if err := api.Sync(context.Background(), s); err != nil { return fmt.Errorf("sync %s: %w", s.ID, err) }
			if err := sources.Save(dbPath, db); err != nil { return err }
			fmt.Fprintf(out, "Indexed %d entries at %s.\n", len(s.Items), s.Revision)
		}
		return nil
	case "search":
		kind, sourceID := "", ""
		var words []string
		for i := 0; i < len(rest); i++ {
			switch rest[i] {
			case "--kind", "--source":
				if i+1 >= len(rest) { return fmt.Errorf("%s requires a value", rest[i]) }
				i++
				if rest[i-1] == "--kind" { kind = rest[i] } else { sourceID = rest[i] }
			default:
				if strings.HasPrefix(rest[i], "--") { return fmt.Errorf("unknown option %s", rest[i]) }
				words = append(words, rest[i])
			}
		}
		if kind != "" && kind != "skill" && kind != "agent" && kind != "mcp" { return fmt.Errorf("invalid kind %s", kind) }
		query := strings.ToLower(strings.Join(words, " "))
		matches := 0
		for _, s := range db.Sources {
			if sourceID != "" && s.ID != sourceID { continue }
			for _, item := range s.Items {
				if kind != "" && item.Kind != kind { continue }
				if !strings.Contains(strings.ToLower(item.Name+" "+item.Path), query) { continue }
				fmt.Fprintf(out, "%-6s %-19s %-30s %s\n", item.Kind, s.ID, item.Name, item.Path)
				matches++
			}
		}
		if matches == 0 { fmt.Fprintln(out, "No indexed matches. Run 'agentstrap source sync all' to populate sources.") }
		return nil
	case "install":
		opts, positional, err := parseOptions(rest)
		if err != nil { return err }
		if len(positional) != 1 || opts.client == "" || opts.topic != "" {
			return fmt.Errorf("usage: source install SOURCE:ITEM --client claude|codex|gemini [--path DIR] [--apply]")
		}
		return installSourceItem("", positional[0], opts, out)
	default:
		return fmt.Errorf("unknown source action %q", action)
	}
}

func installSourceItem(kind, identifier string, opts options, out io.Writer) error {
	sourceID, name, ok := strings.Cut(identifier, ":")
	if !ok || sourceID == "" || name == "" { return fmt.Errorf("use SOURCE:ITEM (e.g. anthropic-skills:frontend-design)") }
	path, err := sources.DatabasePath(); if err != nil { return err }
	db, err := sources.Load(path); if err != nil { return err }
	source, err := db.Find(sourceID); if err != nil { return err }
	item, err := source.Resolve(name); if err != nil { return err }
	if kind != "" && item.Kind != kind { return fmt.Errorf("%s is %s, not %s", identifier, item.Kind, kind) }
	return installRemoteItem(*source, item, opts, out)
}

func installRemoteItem(source sources.Source, item sources.Entry, opts options, out io.Writer) error {
	if opts.client != "claude" && opts.client != "codex" && opts.client != "gemini" {
		return fmt.Errorf("unsupported client %q", opts.client)
	}
	if item.Kind == "mcp" {
		return fmt.Errorf("remote MCP configs require manual review; see source show %s and %s", source.ID, item.Path)
	}
	if item.Kind == "agent" && opts.client == "codex" {
		return fmt.Errorf("Codex does not support importing Claude/Gemini agent markdown files directly")
	}
	if !sources.SafePath(item.Name) || strings.Contains(item.Name, "/") {
		return fmt.Errorf("unsafe resource name %q", item.Name)
	}
	if source.Revision == "" { return fmt.Errorf("source %s is not synced; run source sync %s", source.ID, source.ID) }
	root, err := filepath.Abs(opts.path); if err != nil { return err }
	root, err = filepath.EvalSymlinks(root); if err != nil { return err }
	info, err := os.Stat(root); if err != nil { return err }
	if !info.IsDir() { return fmt.Errorf("not a directory: %s", root) }

	var rel string
	if item.Kind == "skill" {
		prefix := map[string]string{"claude": ".claude/skills", "codex": ".agents/skills", "gemini": ".gemini/skills"}[opts.client]
		rel = filepath.Join(prefix, item.Name)
	} else {
		rel = filepath.Join("."+opts.client, "agents", item.Name+".md")
	}
	target := filepath.Join(root, rel)
	if _, err := os.Lstat(target); err == nil { return fmt.Errorf("target already exists; not overwriting: %s", target) } else if !os.IsNotExist(err) { return err }
	if err := checkNoSymlinkParents(root, filepath.Dir(target)); err != nil { return err }
	var size int64
	for _, f := range item.Files { size += f.Size }
	fmt.Fprintf(out, "%s: %s (%s)\nSource: %s\nRevision: %s\nTarget: %s\nFiles: %d, bytes: %d\nInstall method: %s\nReview: %s\n", map[bool]string{true:"INSTALL", false:"PREVIEW"}[opts.apply], item.Name, item.Kind, source.URL, source.Revision, target, len(item.Files), size, source.InstallMethods[item.Kind], source.Notes)
	if !opts.apply { fmt.Fprintln(out, "Run again with --apply to download and install the pinned files."); return nil }

	data, err := sources.NewGitHubClient().Download(context.Background(), source, item)
	if err != nil { return err }
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil { return err }
	if err := checkNoSymlinkParents(root, filepath.Dir(target)); err != nil { return err }
	if item.Kind == "agent" {
		if len(item.Files) != 1 { return fmt.Errorf("agent should contain exactly one file") }
		return createExclusiveFile(target, data[item.Files[0].Path])
	}
	// Mkdir is exclusive: no existing skill folder can be overwritten.
	if err := os.Mkdir(target, 0755); err != nil { return err }
	created := false
	defer func() { if !created { os.RemoveAll(target) } }()
	for _, f := range item.Files {
		if !sources.SafePath(f.Path) { return fmt.Errorf("unsafe installed path %q", f.Path) }
		filename := filepath.Join(target, filepath.FromSlash(f.Path))
		if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil { return err }
		if err := createExclusiveFile(filename, data[f.Path]); err != nil { return err }
	}
	created = true
	fmt.Fprintf(out, "Installed %s\n", target)
	return nil
}

func createExclusiveFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil { return err }
	if _, err := f.Write(data); err != nil { f.Close(); return err }
	return f.Close()
}

func checkNoSymlinkParents(root, dest string) error {
	rel, err := filepath.Rel(root, dest); if err != nil { return err }
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) { return fmt.Errorf("destination escapes project directory") }
	current := root
	for _, part := range strings.Split(rel, string(filepath.Separator)) {
		if part == "." || part == "" { continue }
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if os.IsNotExist(err) { continue }
		if err != nil { return err }
		if info.Mode()&os.ModeSymlink != 0 { return fmt.Errorf("refusing symlink in destination path %s", current) }
		if !info.IsDir() { return fmt.Errorf("destination component is not a directory: %s", current) }
	}
	return nil
}

// Stable list of source IDs for future TUI use.
func SourceIDs(db sources.Database) []string {
	ids := make([]string, 0, len(db.Sources))
	for _, source := range db.Sources { ids = append(ids, source.ID) }
	sort.Strings(ids)
	return ids
}
