package retriever_github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/joshcarp/gop/gop/retrievertests"

	"github.com/joshcarp/gop/gop"
)

type GithubMock struct {
	content map[string]string
}

func NewMock() GithubMock {
	return GithubMock{content: retrievertests.GithubRequestPaths}
}

func NewMockFromMap(m map[string]string) GithubMock {
	return GithubMock{content: m}
}

func (g GithubMock) ResolveHash(resource string) (string, error) {
	repo, resource, ver, _ := gop.ProcessRequest(resource)
	r, _ := url.Parse("https://" + repo)
	repo = strings.ReplaceAll(repo, r.Host, "")
	path := fmt.Sprintf("/repos%s/contents/%s?ref=%s", repo, resource, ver)
	if val, ok := g.content[path]; ok {
		return val, nil
	}
	return ver, nil
}

func (g GithubMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var query string
	if r.URL.RawQuery != "" {
		query = "?" + r.URL.RawQuery
	}
	if val, ok := g.content[r.URL.Path+query]; ok {
		if r.Header.Get("Accept") == "application/vnd.github.VERSION.sha" {
			_, _ = w.Write([]byte(val))
			return
		}
		resp := GithubResponse{Content: base64.StdEncoding.EncodeToString([]byte(val))}
		content, _ := json.Marshal(resp)
		_, _ = w.Write(content)
		return
	}
	w.WriteHeader(404)
}
