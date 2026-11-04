package catalog

import "testing"

func TestSearch(t *testing.T) {
	items := Search("skill", "video editing", "")
	if len(items) != 2 {
		t.Fatalf("video skills = %d, want 2", len(items))
	}
	items = Search("agent", "animation", "planner")
	if len(items) != 1 || items[0].ID != "animation-planner" {
		t.Fatalf("unexpected result: %+v", items)
	}
	if _, ok := Find("mcp", "playwright"); !ok {
		t.Fatal("missing catalog MCP")
	}
	if _, ok := Find("skill", "playwright"); ok {
		t.Fatal("incorrect kind matched")
	}
}
