package sources

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRegistrySeedAndPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "private", "sources.json")
	db, err := Load(path)
	if err != nil { t.Fatal(err) }
	if db.SchemaVersion != SchemaVersion || len(db.Sources) != 2 { t.Fatalf("bad seeds: %+v", db) }
	anthropic, err := db.Find("anthropic-skills")
	if err != nil || anthropic.URL != "https://github.com/anthropics/skills" { t.Fatalf("anthropic source: %+v %v", anthropic, err) }
	alireza, err := db.Find("claude-skills")
	if err != nil || !alireza.HasKind("agent") || !alireza.HasKind("mcp") { t.Fatalf("second source: %+v %v", alireza, err) }
	if err := db.Add(Source{ID:"test-registry", Name:"Test", URL:"https://github.com/example/repo", Ref:"main", Kinds:[]string{"skill"}}); err != nil { t.Fatal(err) }
	if err := Save(path, db); err != nil { t.Fatal(err) }
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0600 { t.Fatalf("registry permissions %v %v", info, err) }
	next, err := Load(path)
	if err != nil { t.Fatal(err) }
	if len(next.Sources) != 3 { t.Fatalf("persistence: %d sources", len(next.Sources)) }
	if err := next.Remove("test-registry"); err != nil { t.Fatal(err) }
	if err := next.Remove("test-registry"); err == nil { t.Fatal("missing source removed twice") }
}

func TestURLAndPathValidation(t *testing.T) {
	for _, raw := range []string{"https://evil.test/owner/repo", "http://github.com/owner/repo", "https://github.com/owner/repo/tree/main", "https://github.com/owner/../repo", "https://github.com/owner/repo?token=abc"} {
		if _, _, err := GitHubRepo(raw); err == nil { t.Errorf("accepted unsafe GitHub URL %q", raw) }
	}
	owner, repo, err := GitHubRepo("https://github.com/anthropics/skills.git")
	if err != nil || owner != "anthropics" || repo != "skills" { t.Fatalf("%s %s %v", owner, repo, err) }
	for _, p := range []string{"../SECRET", "foo/../bar", "foo//bar", "/abs/file", "foo\\bar", "."} {
		if SafePath(p) { t.Errorf("unsafe path accepted %q", p) }
	}
	if !SafePath("skills/frontend-design/SKILL.md") { t.Fatal("rejected valid skill path") }
	for _, s := range []Source{
		{ID:"Bad ID", URL:"https://github.com/a/b", Ref:"main", Kinds:[]string{"skill"}},
		{ID:"ok", URL:"https://github.com/a/b", Ref:"../main", Kinds:[]string{"skill"}},
		{ID:"ok", URL:"https://github.com/a/b", Ref:"main", Kinds:[]string{"executable"}},
	} {
		if err := Validate(s); err == nil { t.Errorf("accepted invalid source %+v", s) }
	}
}

func TestDiscoverAndAmbiguousNames(t *testing.T) {
	s := Source{ID:"example", URL:"https://github.com/example/repo", Ref:"main", Kinds:[]string{"skill","agent","mcp"}}
	sha := strings.Repeat("a",40)
	tree := []TreeEntry{
		{Path:"skills/one/SKILL.md", SHA:sha, Size:50, Type:"blob", Mode:"100644"},
		{Path:"skills/one/scripts/tool.py", SHA:sha, Size:20, Type:"blob", Mode:"100644"},
		{Path:"docs/one/SKILL.md", SHA:sha, Size:5, Type:"blob", Mode:"100644"},
		{Path:"agents/engineering/reviewer.md", SHA:sha, Size:11, Type:"blob", Mode:"100644"},
		{Path:"agents/CLAUDE.md", SHA:sha, Size:4, Type:"blob", Mode:"100644"},
		{Path:".mcp.json", SHA:sha, Size:10, Type:"blob", Mode:"100644"},
		{Path:"skills/one/unsafe-link", SHA:sha, Size:3, Type:"blob", Mode:"120000"},
	}
	items, err := Discover(s, tree)
	if err != nil || len(items) != 4 { t.Fatalf("items %d err %v",len(items),err) }
	s.Items=items
	if _, err := s.Resolve("one"); err == nil { t.Fatal("duplicate basename resolved ambiguously") }
	skill, err := s.Resolve("skills/one")
	if err != nil || len(skill.Files)!=2 { t.Fatalf("skill: %+v %v", skill,err) }
	agent, err := s.Resolve("reviewer")
	if err != nil || agent.Kind!="agent" { t.Fatalf("agent: %+v %v", agent,err) }
	if _, err := Discover(s, append(tree, TreeEntry{Path:"skills/one/../../outside", Type:"blob", Mode:"100644", SHA:sha})); err == nil { t.Fatal("unsafe tree accepted") }
}
