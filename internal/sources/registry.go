// Package sources manages a persistent catalog of external skill, agent and MCP sources.
package sources

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

//go:embed defaults.json
var seedData []byte

const SchemaVersion = 1

type File struct {
	Path string `json:"path"`
	SHA  string `json:"sha"`
	Size int64  `json:"size"`
}

type Entry struct {
	Kind  string `json:"kind"`
	Name  string `json:"name"`
	Path  string `json:"path"`
	Files []File `json:"files"`
}

type Source struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	URL            string            `json:"url"`
	Ref            string            `json:"ref"`
	Kinds          []string          `json:"kinds"`
	InstallMethods map[string]string `json:"install_methods"`
	Roots          []string          `json:"roots,omitempty"`
	Notes          string            `json:"notes,omitempty"`
	Revision       string            `json:"revision,omitempty"`
	SyncedAt       string            `json:"synced_at,omitempty"`
	Items          []Entry           `json:"items"`
}

type Database struct {
	SchemaVersion int      `json:"schema_version"`
	Sources       []Source `json:"sources"`
}

var validID = regexp.MustCompile(`^[a-z][a-z0-9-]{0,63}$`)
var validRef = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/-]{0,127}$`)

func DatabasePath() (string, error) {
	if dir := os.Getenv("AGENTSTRAP_HOME"); dir != "" {
		return filepath.Join(dir, "sources.json"), nil
	}
	dir, err := os.UserConfigDir()
	if err != nil { return "", err }
	return filepath.Join(dir, "agentstrap", "sources.json"), nil
}

func Load(path string) (Database, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		var seed Database
		if err := json.Unmarshal(seedData, &seed); err != nil { return Database{}, err }
		if err := Save(path, seed); err != nil { return Database{}, err }
		return seed, nil
	}
	if err != nil { return Database{}, err }
	var db Database
	if err := json.Unmarshal(b, &db); err != nil { return db, fmt.Errorf("invalid source database: %w", err) }
	if db.SchemaVersion != SchemaVersion { return db, fmt.Errorf("unsupported source schema version %d", db.SchemaVersion) }
	seen := make(map[string]bool)
	for _, source := range db.Sources {
		if seen[source.ID] { return db, fmt.Errorf("duplicate source %q", source.ID) }
		seen[source.ID] = true
		if err := Validate(source); err != nil { return db, err }
	}
	return db, nil
}

func Save(path string, db Database) error {
	db.SchemaVersion = SchemaVersion
	sort.Slice(db.Sources, func(i, j int) bool { return db.Sources[i].ID < db.Sources[j].ID })
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil { return err }
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil { return err }
	data = append(data, '\n')
	tmp, err := os.CreateTemp(filepath.Dir(path), ".sources-*")
	if err != nil { return err }
	defer os.Remove(tmp.Name())
	if err := tmp.Chmod(0600); err != nil { tmp.Close(); return err }
	if _, err := tmp.Write(data); err != nil { tmp.Close(); return err }
	if err := tmp.Close(); err != nil { return err }
	return os.Rename(tmp.Name(), path)
}

func Validate(s Source) error {
	if !validID.MatchString(s.ID) { return fmt.Errorf("invalid source ID %q", s.ID) }
	if !validRef.MatchString(s.Ref) || strings.Contains(s.Ref, "..") || strings.Contains(s.Ref, "//") { return fmt.Errorf("invalid git ref %q", s.Ref) }
	if _, _, err := GitHubRepo(s.URL); err != nil { return err }
	if len(s.Kinds) == 0 { return fmt.Errorf("source %s: no resource types", s.ID) }
	for _, kind := range s.Kinds { if kind != "skill" && kind != "agent" && kind != "mcp" { return fmt.Errorf("invalid resource type %q", kind) } }
	for _, root := range s.Roots { if root != "" && !SafePath(root) { return fmt.Errorf("unsafe source scan root %q", root) } }
	return nil
}

func (db *Database) Find(id string) (*Source, error) {
	for i := range db.Sources {
		if db.Sources[i].ID == id { return &db.Sources[i], nil }
	}
	return nil, fmt.Errorf("source %q not found", id)
}

func (db *Database) Add(s Source) error {
	if err := Validate(s); err != nil { return err }
	for _, existing := range db.Sources { if existing.ID == s.ID { return fmt.Errorf("source %q already exists", s.ID) } }
	if s.Items == nil { s.Items = []Entry{} }
	db.Sources = append(db.Sources, s)
	return nil
}

func (db *Database) Remove(id string) error {
	for i := range db.Sources {
		if db.Sources[i].ID == id { db.Sources = append(db.Sources[:i], db.Sources[i+1:]...); return nil }
	}
	return fmt.Errorf("source %q not found", id)
}

func (s Source) HasKind(kind string) bool {
	for _, k := range s.Kinds { if k == kind { return true } }
	return false
}

// Resolve permits an unambiguous basename or the full source-relative item path.
func (s Source) Resolve(name string) (Entry, error) {
	var matches []Entry
	for _, item := range s.Items {
		if item.Name == name || item.Path == name { matches = append(matches, item) }
	}
	if len(matches) == 1 { return matches[0], nil }
	if len(matches) > 1 { return Entry{}, fmt.Errorf("ambiguous item %q; specify its full source-relative path", name) }
	return Entry{}, fmt.Errorf("item %q not found (run source sync first)", name)
}
