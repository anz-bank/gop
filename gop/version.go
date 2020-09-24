package gop

import (
	"fmt"
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
	var content []byte
	if cacheFile != "" {
		content, _, _ = cacher.Retrieve(cacheFile)
	}
	repo, path, ver, _ := ProcessRequest(resource)
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
			return CreateResource(repo, path, e.Hash), nil
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
	return CreateResource(repo, path, hash), err
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
