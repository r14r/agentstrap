package sources

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request)(*http.Response,error)
func (f roundTripFunc) RoundTrip(req *http.Request)(*http.Response,error) { return f(req) }
func response(body string) (*http.Response,error) {
	return &http.Response{StatusCode:200,Body:io.NopCloser(strings.NewReader(body)),Header:make(http.Header)},nil
}

func TestSyncAndPinnedDownload(t *testing.T) {
	revision := strings.Repeat("1",40)
	skillData := []byte("---\nname: demo\n---\n")
	header := fmt.Sprintf("blob %d\x00",len(skillData))
	h := sha1.Sum(append([]byte(header),skillData...))
	blobSHA := hex.EncodeToString(h[:])
	s := Source{ID:"demo",URL:"https://github.com/example/repo",Ref:"main",Kinds:[]string{"skill"},Roots:[]string{"skills"}}
	tree := githubTree{SHA:strings.Repeat("2",40),Tree:[]TreeEntry{{Path:"skills/demo/SKILL.md",Mode:"100644",Type:"blob",SHA:blobSHA,Size:int64(len(skillData))}}}
	treeData,_ := json.Marshal(tree)
	client := &GitHubClient{HTTP:&http.Client{Transport:roundTripFunc(func(r *http.Request)(*http.Response,error){
		if !strings.Contains(r.Host,"githubusercontent.com") && !strings.Contains(r.Host,"api.github.com") { t.Fatalf("untrusted host %s",r.Host) }
		switch {
		case strings.Contains(r.URL.Path,"/commits/main"): return response(fmt.Sprintf(`{"sha":"%s"}`,revision))
		case strings.Contains(r.URL.Path,"/git/trees/"): return response(string(treeData))
		case strings.HasSuffix(r.URL.Path,"skills/demo/SKILL.md"): return response(string(skillData))
		default: return nil,fmt.Errorf("unexpected path %s",r.URL)
		}
	})}}
	if err := client.Sync(context.Background(),&s); err != nil { t.Fatal(err) }
	if s.Revision!=revision || len(s.Items)!=1 { t.Fatalf("index: %+v",s) }
	download,err := client.Download(context.Background(),s,s.Items[0])
	if err != nil || string(download["SKILL.md"])!=string(skillData) { t.Fatalf("download %v %v",download,err) }
	client.HTTP = &http.Client{Transport:roundTripFunc(func(r *http.Request)(*http.Response,error){ return response("tampered content") })}
	if _,err := client.Download(context.Background(),s,s.Items[0]);err==nil { t.Fatal("accepted tampered download") }
}
