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

func ReplaceImports(modFile []byte, sourceFile []byte) ([]byte, error) {
	var mod Modules
	if err := yaml.Unmarshal(modFile, &mod); err != nil {
		return nil, err
	}
	for pattern, resolve := range mod.Imports {
		repoFrom, _, verFrom := gop.ProcessRepo(pattern)
		repoTo, _, verTo := gop.ProcessRepo(resolve)
		sourceFile = []byte(ReplaceSpecificImport(string(sourceFile), repoFrom, verFrom, repoTo, verTo))
	}
	return sourceFile, nil
}

/* ReplaceSpecificImport replaces a specific import in content */
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

/* RetrieveAndReplace retrieves a resource and replaces its import statements with the patterns described in the import file */
func RetrieveAndReplace(retriever gop.Retriever, resource string, importFile string) ([]byte, bool, error) {
	content, _, err := retriever.Retrieve(resource)
	if !(err != nil || content == nil || len(content) == 0) {
		importFilecontents, _, err := retriever.Retrieve(AddPath(resource, importFile))
		if err != nil {
			return content, false, nil
		}
		if reindexed, err := ReplaceImports(importFilecontents, content); err == nil {
			content = reindexed
		}
		return content, false, nil
	}
	return content, false, err
}
