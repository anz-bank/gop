package modules

import (
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

func New(retriever gop.Retriever, importFile string) Retriever {
	return Retriever{retriever: retriever, importFile: importFile}
}

type Retriever struct {
	retriever  gop.Retriever
	importFile string
}

/* ReplaceImports replaces the contents of sourcefile with the imports in modfile */
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
func ReplaceSpecificImport(content string, oldrepo, oldver, newrepo, newver string) string {
	var pth string
	if oldver != "" {
		oldver = "(?P<version>" + regexp.QuoteMeta(oldver) + ")"
	}
	re := fmt.Sprintf(`(?:%s)(?P<path>[a-zA-Z0-9/._\-]*)@*%s(?:\S)?`,
		regexp.QuoteMeta(oldrepo), oldver)
	impRe := regexp.MustCompile(re)
	for _, match := range impRe.FindAllStringSubmatch(content, -1) {
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
		for _, match := range impRe.FindAllString(content, -1) {
			newImport := fmt.Sprintf("%s@%s", path.Join(newrepo, pth), newver)
			content = strings.ReplaceAll(content, match, newImport)
		}
	}
	return content
}

/* Retrieve retrieves a resource and replaces its import statements with the patterns described in the import file */
func (a Retriever) Retrieve(resource string) ([]byte, bool, error) {
	content, _, err := a.retriever.Retrieve(resource)
	if !(err != nil || content == nil || len(content) == 0) {
		importFilecontents, _, err := a.retriever.Retrieve(AddPath(resource, a.importFile))
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
