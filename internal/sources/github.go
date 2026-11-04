package sources

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"
)

const maxIndexBytes int64 = 16 << 20
const maxArtifactBytes int64 = 12 << 20
const maxArtifactFiles = 250

var repoSegment = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
var resourceName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`)

func GitHubRepo(raw string) (string, string, error) {
	u, err := url.Parse(raw)
	if err != nil { return "", "", err }
	if u.Scheme != "https" || u.Host != "github.com" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return "", "", fmt.Errorf("only https://github.com/OWNER/REPO sources are supported")
	}
	segments := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(segments) != 2 { return "", "", fmt.Errorf("expected GitHub repository URL, got %q", raw) }
	owner, repo := segments[0], strings.TrimSuffix(segments[1], ".git")
	if owner == "." || owner == ".." || repo == "." || repo == ".." || !repoSegment.MatchString(owner) || !repoSegment.MatchString(repo) {
		return "", "", fmt.Errorf("invalid repository URL %q", raw)
	}
	return owner, repo, nil
}

func SafePath(s string) bool {
	if s == "" || strings.HasPrefix(s, "/") || strings.Contains(s, "\\") || strings.ContainsRune(s, 0) { return false }
	for _, part := range strings.Split(s, "/") {
		if part == "." || part == ".." || part == "" { return false }
	}
	return true
}

type TreeEntry struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	SHA  string `json:"sha"`
	Size int64  `json:"size"`
}

type githubCommit struct { SHA string `json:"sha"` }
type githubTree struct {
	SHA string `json:"sha"`
	Truncated bool `json:"truncated"`
	Tree []TreeEntry `json:"tree"`
}

type GitHubClient struct { HTTP *http.Client }

func NewGitHubClient() *GitHubClient {
	return &GitHubClient{HTTP: &http.Client{Timeout: 35 * time.Second}}
}

func (c *GitHubClient) readURL(ctx context.Context, raw string, max int64) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil { return nil, err }
	req.Header.Set("User-Agent", "agentstrap/0.4")
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := c.HTTP.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub HTTP %d for %s", resp.StatusCode, req.URL.Path)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, max+1))
	if err != nil { return nil, err }
	if int64(len(data)) > max { return nil, fmt.Errorf("response exceeds %d bytes", max) }
	return data, nil
}

func (c *GitHubClient) Sync(ctx context.Context, s *Source) error {
	owner, repo, err := GitHubRepo(s.URL)
	if err != nil { return err }
	// Resolve the ref to an immutable commit to prevent branch movement during install.
	commitURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", owner, repo, url.PathEscape(s.Ref))
	b, err := c.readURL(ctx, commitURL, 1<<20)
	if err != nil { return err }
	var commit githubCommit
	if err = json.Unmarshal(b, &commit); err != nil { return err }
	if !validSHA(commit.SHA) { return fmt.Errorf("invalid GitHub commit SHA") }

	treeURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/trees/%s?recursive=1", owner, repo, commit.SHA)
	b, err = c.readURL(ctx, treeURL, maxIndexBytes)
	if err != nil { return err }
	var tree githubTree
	if err = json.Unmarshal(b, &tree); err != nil { return err }
	if tree.Truncated { return fmt.Errorf("GitHub tree is truncated; refusing incomplete index") }

	items, err := Discover(*s, tree.Tree)
	if err != nil { return err }
	s.Items = items
	s.Revision = commit.SHA
	s.SyncedAt = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func validSHA(s string) bool {
	if len(s) != 40 { return false }
	for _, c := range s { if !strings.ContainsRune("0123456789abcdefABCDEF", c) { return false } }
	return true
}

func withinRoots(s Source, target string) bool {
	if len(s.Roots) == 0 { return true }
	for _, root := range s.Roots {
		if target == root || strings.HasPrefix(target, root+"/") { return true }
	}
	return false
}

func Discover(s Source, tree []TreeEntry) ([]Entry, error) {
	// Reject an unsafe tree rather than writing files outside the target project.
	blobs := make([]TreeEntry, 0, len(tree))
	for _, file := range tree {
		if !SafePath(file.Path) { return nil, fmt.Errorf("unsafe source tree path %q", file.Path) }
		if file.Type != "blob" || file.Mode != "100644" && file.Mode != "100755" { continue }
		if !validSHA(file.SHA) { return nil, fmt.Errorf("invalid blob SHA for %s", file.Path) }
		if file.Size < 0 { return nil, fmt.Errorf("invalid blob size for %s", file.Path) }
		blobs = append(blobs, file)
	}
	var items []Entry
	for _, file := range blobs {
		if !withinRoots(s, file.Path) { continue }
		var entry Entry
		switch {
		case s.HasKind("skill") && path.Base(file.Path) == "SKILL.md":
			parent := path.Dir(file.Path)
			if parent == "." { continue }
			name := path.Base(parent)
			if !resourceName.MatchString(name) { continue }
			entry = Entry{Kind: "skill", Name: name, Path: parent}
			var sum int64
			for _, asset := range blobs {
				if !strings.HasPrefix(asset.Path, parent+"/") { continue }
				if len(entry.Files) >= maxArtifactFiles || asset.Size > maxArtifactBytes-sum {
					return nil, fmt.Errorf("%s exceeds package limits", parent)
				}
				entry.Files = append(entry.Files, File{Path: strings.TrimPrefix(asset.Path, parent+"/"), SHA: asset.SHA, Size: asset.Size})
				sum += asset.Size
			}
		case s.HasKind("agent") && strings.HasPrefix(file.Path, "agents/") && strings.HasSuffix(file.Path, ".md"):
			name := strings.TrimSuffix(path.Base(file.Path), ".md")
			if !resourceName.MatchString(name) || strings.EqualFold(name, "CLAUDE") || strings.EqualFold(name, "README") { continue }
			entry = Entry{Kind: "agent", Name: name, Path: file.Path, Files: []File{{Path: path.Base(file.Path), SHA: file.SHA, Size: file.Size}}}
		case s.HasKind("mcp") && (path.Base(file.Path) == ".mcp.json" || path.Base(file.Path) == "mcp.json"):
			name := "root-mcp"
			if path.Dir(file.Path) != "." { name = strings.ReplaceAll(path.Dir(file.Path), "/", "-")+"-mcp" }
			entry = Entry{Kind: "mcp", Name: name, Path: file.Path, Files: []File{{Path: path.Base(file.Path), SHA: file.SHA, Size: file.Size}}}
		default:
			continue
		}
		items = append(items, entry)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Kind != items[j].Kind { return items[i].Kind < items[j].Kind }
		return items[i].Path < items[j].Path
	})
	return items, nil
}

func (c *GitHubClient) Download(ctx context.Context, s Source, item Entry) (map[string][]byte, error) {
	owner, repo, err := GitHubRepo(s.URL)
	if err != nil { return nil, err }
	if !validSHA(s.Revision) { return nil, fmt.Errorf("source %s is not synchronized", s.ID) }
	if item.Kind != "skill" && item.Kind != "agent" && item.Kind != "mcp" { return nil, fmt.Errorf("unsupported item kind %q", item.Kind) }
	if len(item.Files) == 0 || len(item.Files) > maxArtifactFiles { return nil, fmt.Errorf("invalid item file count") }
	output := make(map[string][]byte)
	var sum int64
	for _, f := range item.Files {
		if !SafePath(f.Path) || !validSHA(f.SHA) || f.Size < 0 || f.Size > maxArtifactBytes-sum {
			return nil, fmt.Errorf("unsafe file manifest for %s", f.Path)
		}
		fullPath := f.Path
		if item.Kind == "skill" { fullPath = item.Path + "/" + f.Path }
		// Escape each path segment separately; preserve directory separators.
		segments := strings.Split(fullPath, "/")
		for i, segment := range segments { segments[i] = url.PathEscape(segment) }
		rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, s.Revision, strings.Join(segments, "/"))
		data, err := c.readURL(ctx, rawURL, f.Size)
		if err != nil { return nil, err }
		if int64(len(data)) != f.Size { return nil, fmt.Errorf("file size changed for %s", fullPath) }
		header := fmt.Sprintf("blob %d ", len(data))
		hash := sha1.Sum(append([]byte(header), data...))
		if !strings.EqualFold(hex.EncodeToString(hash[:]), f.SHA) {
			return nil, fmt.Errorf("content hash mismatch for %s", fullPath)
		}
		output[f.Path] = data
		sum += int64(len(data))
	}
	return output, nil
}
