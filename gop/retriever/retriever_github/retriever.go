package retriever_github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/joshcarp/gop/gop"
)

type Retriever struct {
	token map[string]string
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
	return Retriever{token: tokens}
}

func getToken(token map[string]string, resource string) string {
	u, _ := url.Parse("https://" + resource)
	return token[u.Host]
}

func (a Retriever) Retrieve(resource string) ([]byte, bool, error) {
	var resp *http.Response
	var apibase string
	var repo, path, ver string
	var err error
	var b []byte
	var res GithubResponse

	repo, path, ver, err = gop.ProcessRequest(resource)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", gop.BadRequestError, err)
	}
	apibase = gop.GetApiURL(resource)

	req, err := url.Parse(
		fmt.Sprintf(
			"https://%s/repos/%s/contents/%s?ref=%s",
			apibase, repo, path, ver))
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

	resp, err = http.DefaultClient.Do(r)
	if err != nil {
		return b, false, fmt.Errorf("%s: %w", gop.GithubFetchError, err)
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
