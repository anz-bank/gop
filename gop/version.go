package gop

import (
	"bufio"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

/* Modules is the representation of the module file: */
type Modules struct {
	Direct []Direct `yaml:"direct"`
	Sum    []Direct `yaml:"sum"`
}

type Direct struct {
	Pattern string `yaml:"pattern"`
	Resolve string `yaml:"resolve"`
}

func ReplaceImports(retriever Retriever, resource string, content []byte) ([]byte, error) {
	var mod Modules
	content1, _, err := retriever.Retrieve(resource)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(content1, &mod); err != nil {
		return nil, err
	}
	for _, e := range mod.Direct {
		repoFrom, _, verFrom := ProcessRepo(e.Pattern)
		repoTo, _, verTo := ProcessRepo(e.Resolve)
		content = []byte(ReplaceSpecificImport(string(content), repoFrom, verFrom, repoTo, verTo))
	}
	return content, nil
}

func ReplaceSpecificImport(content string, oldimp, oldver, newimp, newver string) string {
	var pth string
	if oldver != "" {
		oldver = "(?P<version>" + regexp.QuoteMeta(oldver) + ")"
	}
	//oldver = `[a-zA-Z0-9/._]`
	//else {
	//oldver = regexp.QuoteMeta(oldver)
	//}
	re := fmt.Sprintf(`(?:%s)(?P<path>[a-zA-Z0-9/._\-]*)@*%s(?:\S)?`,
		regexp.QuoteMeta(oldimp), oldver)

	impRe := regexp.MustCompile(re)
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		txt := scanner.Text()
		for _, match := range impRe.FindAllStringSubmatch(scanner.Text(), -1) {
			if match == nil {
				continue
			}
			for i, name := range impRe.SubexpNames() {
				if match[i] != "" {
					switch name {
					case "path":
						pth = match[i]
					}
				}
			}
			for _, match := range impRe.FindAllString(txt, -1) {
				newImport := fmt.Sprintf("%s@%s", path.Join(newimp, pth), newver)
				content = strings.ReplaceAll(content, match, newImport)
			}
		}
	}
	return content

}

/* LoadVersion returns the version from a version */
func LoadVersion(retriever Retriever, cacher Cacher, resolver func(string) (string, error), cacheFile, resource string) (string, error) {
	var content []byte
	if cacheFile != "" {
		content, _, _ = retriever.Retrieve(cacheFile)
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
	if cacher == nil {
		return resource, nil
	}
	hash, err := resolver(repoVer)
	if err != nil {
		return AddPath(repoVer, path), nil
	}
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
