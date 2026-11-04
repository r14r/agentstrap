// Package catalog contains built-in templates, not downloaded third-party code.
package catalog

import "strings"

type Item struct {
	Kind, ID, Description, Guide, Command string
	Topics, Args []string
}

var Entries = []Item{
	{"skill", "go-development", "Go development, tests and reproducible builds", "Inspect go.mod first. Prefer standard-library solutions and table-driven tests. Run gofmt, go vet ./..., and go test ./....", "", []string{"development"}, nil},
	{"skill", "python-development", "Python packaging, typing and automated tests", "Inspect pyproject.toml and the dependency manager. Prefer typed APIs, focused tests and documented commands.", "", []string{"development"}, nil},
	{"skill", "code-review", "Security, correctness and regression review", "Review the actual diff and surrounding callers. Prioritize correctness, security and missing tests; report concrete findings.", "", []string{"development"}, nil},
	{"skill", "video-editing", "Plan and validate reproducible video editing", "Inspect source resolution, frame rate and audio. Preserve source media. Document cutting decisions and export settings.", "", []string{"video-editing"}, nil},
	{"skill", "ffmpeg-workflow", "Design safe FFmpeg media processing pipelines", "Probe media streams first. Quote filenames. Write outputs to new files. Verify audio sync, streams and codec.", "", []string{"video-editing", "animation"}, nil},
	{"skill", "animation-workflow", "Storyboard, keyframes and animation validation", "Define frame rate, timeline and shot list. Plan key poses, easing and iteration checkpoints. Validate rendered motion.", "", []string{"animation"}, nil},
	{"agent", "software-architect", "Architecture review specialist", "Analyze boundaries, interfaces, tradeoffs, tests and migrations. Do not modify files without permission.", "", []string{"development"}, nil},
	{"agent", "qa-reviewer", "Testing and regression specialist", "Find test gaps and reproducible failure cases. Distinguish verified results from untested assumptions.", "", []string{"development"}, nil},
	{"agent", "video-editor", "Video editing and delivery specialist", "Build a shot-by-shot plan, considering pacing, captions, transitions and audio. Preserve media originals.", "", []string{"video-editing"}, nil},
	{"agent", "animation-planner", "Storyboard and animation planning specialist", "Translate the brief into shots, keyframes, timing and assets. Check continuity and motion arcs.", "", []string{"animation"}, nil},
	{"mcp", "filesystem", "Local filesystem MCP; grants the target project directory access", "", "npx", []string{"development", "video-editing", "animation"}, []string{"-y", "@modelcontextprotocol/server-filesystem"}},
	{"mcp", "playwright", "Playwright MCP for browser inspection and automation", "", "npx", []string{"development"}, []string{"-y", "@playwright/mcp@latest"}},
}

func NormalizeTopic(topic string) string {
	return strings.ReplaceAll(strings.Join(strings.Fields(strings.ToLower(strings.ReplaceAll(topic, "_", "-"))), "-"), " ", "-")
}

func Find(kind, id string) (Item, bool) {
	for _, item := range Entries {
		if item.Kind == kind && item.ID == id {
			return item, true
		}
	}
	return Item{}, false
}

func Search(kind, topic, query string) []Item {
	topic = NormalizeTopic(topic)
	query = strings.ToLower(strings.TrimSpace(query))
	var results []Item
	for _, item := range Entries {
		if kind != "" && kind != item.Kind {
			continue
		}
		if topic != "" {
			matched := false
			for _, tag := range item.Topics {
				if NormalizeTopic(tag) == topic {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		if query != "" && !strings.Contains(strings.ToLower(item.ID+" "+item.Description+" "+strings.Join(item.Topics, " ")), query) {
			continue
		}
		results = append(results, item)
	}
	return results
}
