package retriever_github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/joshcarp/gop/gop/modules"

	"github.com/joshcarp/gop/gop"
)

type Retriever struct {
	token   map[string]string
	Client  *http.Client
	ApiBase string
}

type GithubResponse struct {
	Message          string `json:"message,omitempty"`
	DocumentationURL string `json:"documentation_url,omitempty"`
	Name             string `json:"name,omitempty"`
	Path             string `json:"path,omitempty"`
	Sha              string `json:"sha,omitempty"`
	Size             int    `json:"size,omitempty"`
	URL              string `json:"url,omitempty"`
	HTMLURL          string `json:"html_url,omitempty"`
	GitURL           string `json:"git_url,omitempty"`
	DownloadURL      string `json:"download_url,omitempty"`
	Type             string `json:"type,omitempty"`
	Content          string `json:"content,omitempty"`
	Encoding         string `json:"encoding,omitempty"`
	Links            struct {
		Self string `json:"self,omitempty"`
		Git  string `json:"git,omitempty"`
		HTML string `json:"html,omitempty"`
	} `json:"_links,omitempty"`
}

/* New returns a retriever with a key/value pairs of <host>, <token> eg: New("github.com", "abcdef") */
func New(tokens map[string]string) Retriever {
	if tokens == nil {
		tokens = map[string]string{}
	}
	return Retriever{token: tokens, Client: &http.Client{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment}}}
}

func getToken(token map[string]string, resource string) string {
	u, _ := url.Parse("https://" + resource)
	return token[u.Host]
}

/* Resolve Resolves a github resource to its hash */
func (a Retriever) Resolve(resource string) (string, error) {
	if a.ApiBase == "" {
		a.ApiBase = modules.GetApiURL(resource)
		if a.ApiBase == "" {
			return "", gop.BadRequestError
		}
		a.ApiBase = "https://" + a.ApiBase
	}

	heder := http.Header{}
	repo, _, ref := gop.ProcessRepo(resource)
	if ref == "" {
		ref = "HEAD"
	}
	repoURL, _ := url.Parse("httpps://" + repo)
	heder.Add("accept", "application/vnd.github.VERSION.sha")
	heder.Add("authorization", "Bearer "+a.token[repoURL.Host])
	u, err := url.Parse(fmt.Sprintf("%s/repos%s/commits/%s", a.ApiBase, repoURL.Path, ref))
	if err != nil {
		return "", gop.BadRequestError
	}

	r := &http.Request{
		Method: "GET",
		URL:    u,
		Header: heder,
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil || resp == nil || resp.Body == nil {
		return "", gop.GitCloneError
	}
	if err := gop.HandleHTTPStatus(resp.StatusCode); err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if strings.Contains(string(b), `"message":"Not Found"`) {
		return "", gop.BadReferenceError
	}
	if strings.Contains(string(b), `"message":"Not Found"`) {
		return "", gop.BadReferenceError
	}
	return string(b), nil
}
func (a Retriever) Retrieve(resource string) ([]byte, bool, error) {
	var resp *http.Response
	var repo, path, ver string
	var err error
	var b []byte
	var res GithubResponse

	repo, path, ver, err = gop.ProcessRequest(resource)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", gop.BadRequestError, err)
	}
	h, _ := url.Parse("https://" + repo)
	repo = strings.ReplaceAll(repo, h.Host+"/", "")

	if a.ApiBase == "" {
		a.ApiBase = "https://" + modules.GetApiURL(resource)
	}

	req, err := url.Parse(
		fmt.Sprintf(
			"%s/repos/%s/contents/%s?ref=%s",
			a.ApiBase, repo, path, ver))
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", gop.BadRequestError, err)
	}
	heder := http.Header{}
	heder.Add("accept", "application/vnd.github.v3+json")

	if b := getToken(a.token, resource); b != "" {
		heder.Add("authorization", "token "+b)
	}

	r := &http.Request{
		Method: "GET",
		URL:    req,
		Header: heder,
	}

	resp, err = a.Client.Do(r)
	if err != nil {
		return b, false, fmt.Errorf("%s: %w", gop.GithubFetchError, err)
	}
	if err := gop.HandleHTTPStatus(resp.StatusCode); err != nil {
		return b, false, err
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return b, false, fmt.Errorf("%s: %w", gop.FileReadError, err)
	}
	if err = json.Unmarshal(b, &res); err != nil {
		return nil, false, fmt.Errorf("%s: %w", gop.FileReadError, err)
	}
	if resp.StatusCode == 404 {
		return nil, false, fmt.Errorf("%s", res.Message)
	}
	b, err = base64.StdEncoding.DecodeString(res.Content)
	return b, false, err
}
