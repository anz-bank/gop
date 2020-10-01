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
	log       gop.Logger
}

func NewLoader(gopper gop.Gopper, resolver Resolver, cacheFile string, log gop.Logger) Loader {
	return Loader{
		gopper:    gopper,
		resolver:  resolver,
		cacheFile: cacheFile,
		log:       log,
	}
}

/* Resolve resolves an import from resource to a full commit hash */
func (a Loader) Resolve(resource string) string {
	if _, _, v := gop.ProcessRepo(resource); gop.IsHash(v) {
		return resource
	}
	ver, err := LoadVersion(a.gopper, a.gopper, a.resolver, a.cacheFile, resource)
	if err != nil {
		return ""
	}
	return ver
}

/* Update updates the base of version to the version */
func (a Loader) Update(version string) error {
	x, _, z := gop.ProcessRepo(version)
	return a.UpdateTo(x, gop.CreateResource(x, "", z))
}

func (a Loader) UpdateAll() error {
	var content []byte
	if a.cacheFile != "" {
		content, _, _ = a.gopper.Retrieve(a.cacheFile)
	}
	mod := Modules{}
	if err := yaml.Unmarshal(content, &mod); err != nil {
		return err
	}
	for base := range mod.Imports {
		if x, _, z := gop.ProcessRepo(base); z == "" {
			new := a.Resolve(x + "@HEAD")
			a.log("%s -> %s", x, new)
			mod.Imports[base] = new
		}
	}
	newfile, err := yaml.Marshal(mod)
	if err != nil {
		return err
	}
	if err := a.gopper.Cache(a.cacheFile, newfile); err != nil {
		return err
	}
	return nil
}

func (a Loader) UpdateTo(old, new string) error {
	var content []byte
	if a.cacheFile != "" {
		content, _, _ = a.gopper.Retrieve(a.cacheFile)
	}
	mod := Modules{}
	if err := yaml.Unmarshal(content, &mod); err != nil {
		return err
	}
	repo, _, _ := gop.ProcessRepo(new)
	hash, err := a.resolver(new)
	if err != nil {
		return err
	}
	a.log("updating %s -> %s : %s", old, new, hash)
	if mod.Imports == nil {
		mod.Imports = make(map[string]string)
	}
	mod.Imports[old] = gop.CreateResource(repo, "", hash)
	newfile, err := yaml.Marshal(mod)
	if err != nil {
		return err
	}
	if err := a.gopper.Cache(a.cacheFile, newfile); err != nil {
		return err
	}
	return nil
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
