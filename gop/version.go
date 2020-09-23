package gop

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

/* LoadVersion returns the version from a version */
func LoadVersion(content []byte, repo string) string {
	re := regexp.MustCompile(regexp.QuoteMeta(repo) + ".*")
	for _, e := range re.FindAllString(string(content), -1) {
		_, _, ver, _ := ProcessRequest(e)
		return ver
	}
	return ""
}

/* ResolveHash Resolves a github resource to its hash */
func ResolveHash(resource string) (string, error) {
	base := GetApiURL(resource)
	heder := http.Header{}
	repo, _, ref, _ := ProcessRequest(resource)
	repoURL, _ := url.Parse("httpps://" + repo)
	heder.Add("accept", "application/vnd.github.VERSION.sha")
	u, err := url.Parse(fmt.Sprintf("https://%s/repos%s/commits/%s", base, repoURL.Path, ref))
	if err != nil {
		return "", BadRequestError
	}

	r := &http.Request{
		Method: "GET",
		URL:    u,
		Header: heder,
	}
	resp, err := http.DefaultClient.Do(r)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func GetApiURL(resource string) string {
	requestedurl, _ := url.Parse("https://" + resource)
	switch requestedurl.Host {
	case "github.com":
		return "api.github.com"
	default:
		return fmt.Sprintf("%s/api/v3", requestedurl.Host)
	}
	return ""
}
