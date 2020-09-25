package modules

import (
	"fmt"
	"net/url"

	"github.com/joshcarp/gop/gop"

	"gopkg.in/yaml.v2"
)

/* Resolver resolves a git ref to a hash */
type Resolver func(string) (string, error)

type Loader struct {
	gopper    gop.Gopper
	resolver  Resolver
	cacheFile string
}

func NewLoader(gopper gop.Gopper, resolver Resolver, cacheFile string) Loader {
	return Loader{
		gopper:    gopper,
		resolver:  resolver,
		cacheFile: cacheFile,
	}
}

/* Retrieve resolves an import from resource to a full commit hash */
func (a Loader) Retrieve(resource string) ([]byte, bool, error) {
	if _, _, v := gop.ProcessRepo(resource); len(v) == 20 {
		return []byte(v), false, nil
	}
	ver, err := LoadVersion(a.gopper, a.gopper, a.resolver, a.cacheFile, resource)
	if err != nil {
		return nil, false, err
	}
	return []byte(ver), false, nil
}

/* LoadVersion returns the version from a version */
func LoadVersion(retriever gop.Retriever, cacher gop.Cacher, resolver Resolver, cacheFile, resource string) (string, error) {
	var content []byte
	if cacheFile != "" {
		content, _, _ = retriever.Retrieve(cacheFile)
	}
	repo, path, ver := gop.ProcessRepo(resource)
	var repoVer = repo
	if ver != "" {
		repoVer += "@" + ver
	}

	mod := Modules{}
	if err := yaml.Unmarshal(content, &mod); err != nil {
		return "", err
	}
	if val, ok := mod.Imports[repoVer]; ok {
		return AddPath(val, path), nil
	}
	if cacher == nil {
		return resource, nil
	}
	hash, err := resolver(repoVer)
	if err != nil {
		return "", gop.GithubFetchError
	}
	resolve := gop.CreateResource(repo, "", hash)
	if mod.Imports == nil {
		mod.Imports = map[string]string{}
	}
	mod.Imports[repoVer] = resolve
	newfile, err := yaml.Marshal(mod)
	if err != nil {
		return "", err
	}
	if err := cacher.Cache(cacheFile, newfile); err != nil {
		return "", err
	}
	return AddPath(resolve, path), err
}

func AddPath(repoVer string, path string) string {
	a, _, c, _ := gop.ProcessRequest(repoVer)
	return gop.CreateResource(a, path, c)
}
func GetApiURL(resource string) string {
	requestedurl, _ := url.Parse("https://" + resource)
	switch requestedurl.Host {
	case "github.com":
		return "api.github.com"
	}
	return fmt.Sprintf("%s/api/v3", requestedurl.Host)
}
