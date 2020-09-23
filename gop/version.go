package gop

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ghodss/yaml"
)

/* Modules is the representation of the module file: */
type Modules struct {
	Direct []Direct `yaml:"direct"`
}

type Direct struct {
	Repo string `yaml:"repo"`
	Hash string `yaml:"hash"`
}

/* LoadVersion returns the version from a version */
func LoadVersion(cacher Gopper, resolver func(string) (string, error), cacheFile, resource string) (string, error) {
	content, _, _ := cacher.Retrieve(cacheFile)
	repo, _, ver, _ := ProcessRequest(resource)
	var repoVer = repo
	if ver != "" {
		repoVer += "@" + ver
	}

	mod := Modules{}
	if err := yaml.Unmarshal(content, &mod); err != nil {
		return "", err
	}
	for _, e := range mod.Direct {
		if e.Repo == repoVer {
			return e.Hash, nil
		}
	}
	hash, _ := resolver(repoVer)
	entry := Direct{Repo: repoVer, Hash: hash}
	mod.Direct = append(mod.Direct, entry)
	newfile, err := yaml.Marshal(mod)
	if err != nil {
		return "", err
	}
	if err := cacher.Cache(cacheFile, newfile); err != nil {
		return "", err
	}
	return CreateResource(repo, "", hash), err
}

/* ResolveHash Resolves a github resource to its hash */
func ResolveHash(resource string) (string, error) {
	base := GetApiURL(resource)
	heder := http.Header{}
	repo, _, ref, _ := ProcessRequest(resource)
	if ref == "" {
		ref = "HEAD"
	}
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
