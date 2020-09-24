package gop

import (
	"fmt"
	"net/url"

	"gopkg.in/yaml.v2"
)

/* Modules is the representation of the module file: */
type Modules struct {
	Direct []Direct `yaml:"direct"`
}

type Direct struct {
	Pattern string `yaml:"pattern"`
	Resolve string `yaml:"resolve"`
}

/* LoadVersion returns the version from a version */
func LoadVersion(cacher Gopper, resolver func(string) (string, error), cacheFile, resource string) (string, error) {
	var content []byte
	if cacheFile != "" {
		content, _, _ = cacher.Retrieve(cacheFile)
	}
	repo, path, ver := ProcessRepo(resource)
	var repoVer = repo
	if ver != "" {
		repoVer += "@" + ver
	}

	mod := Modules{}
	if err := yaml.Unmarshal(content, &mod); err != nil {
		return "", err
	}
	for _, e := range mod.Direct {
		if e.Pattern == repoVer {
			return AddPath(e.Resolve, path), nil
		}
	}
	hash, _ := resolver(repoVer)
	entry := Direct{Pattern: repoVer, Resolve: CreateResource(repo, "", hash)}
	mod.Direct = append(mod.Direct, entry)
	newfile, err := yaml.Marshal(mod)
	if err != nil {
		return "", err
	}
	if err := cacher.Cache(cacheFile, newfile); err != nil {
		return "", err
	}
	return AddPath(entry.Resolve, path), err
}
func AddPath(repoVer string, path string) string {
	a, _, c, _ := ProcessRequest(repoVer)
	return CreateResource(a, path, c)
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
