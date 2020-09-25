package modules

import (
	"bufio"
	"fmt"
	"path"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/joshcarp/gop/gop"
)

/* Modules is the representation of the module file: */
type Modules struct {
	Imports map[string]string `yaml:"imports"`
}

func ReplaceImports(retriever gop.Retriever, resource string, content []byte) ([]byte, error) {
	var mod Modules
	content1, _, err := retriever.Retrieve(resource)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(content1, &mod); err != nil {
		return nil, err
	}
	for pattern, resolve := range mod.Imports {
		repoFrom, _, verFrom := gop.ProcessRepo(pattern)
		repoTo, _, verTo := gop.ProcessRepo(resolve)
		content = []byte(ReplaceSpecificImport(string(content), repoFrom, verFrom, repoTo, verTo))
	}
	return content, nil
}

func ReplaceSpecificImport(content string, oldimp, oldver, newimp, newver string) string {
	var pth string
	if oldver != "" {
		oldver = "(?P<version>" + regexp.QuoteMeta(oldver) + ")"
	}
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
