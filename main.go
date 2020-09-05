package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"path"
	"regexp"
	"strings"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

var syslimportRegex = `(?:#import.*)|(?:import )(?:\/\/)?(?P<import>.*)`

type AppConfig struct {
	StartingBalance int64 `yaml:"startingBalance"`
	username        string
	token           string
	saveLocation    string
}

func main() {
	log.Fatal(pbmod.Serve(context.Background(), LoadService))

}

func NewKeyValue() *pbmod.KeyValue {
	var a string
	return &pbmod.KeyValue{
		Extra:    &a,
		Repo:     "",
		Resource: "",
		Value:    "",
		Version:  "",
	}
}

func LoadService(ctx context.Context, a AppConfig) (*pbmod.ServiceInterface, error) {
	s := server{}
	s.retrievers = []retriever{s.retrieveSyslPb, s.retrieveGit}
	s.savers = []saver{s.saveToFile, s.saveToPbJsonFile}
	s.posts = []post{processSysl}
	return &pbmod.ServiceInterface{
		GetResourceList: s.GetResource,
	}, nil
}

func importFile(initialrepo, initialImport, ver string, retrievers []retriever) (*pbmod.KeyValue, error) {
	var file = NewKeyValue()
	file.Repo = initialrepo
	file.Resource = initialImport
	file.Version = ver

	for _, r := range retrievers {
		_ = r(file)
		if file.Value != "" {
			break
		}
	}
	if file.Value == "" {
		return nil, fmt.Errorf("Error loading file")
	}
	return file, nil
}

/*
github.com/a/b/file.ext -> (https://github.com/a/b/file.ext, file.ext)
*/
func processRequest(resource string) (string, string) {
	parts := strings.Split(resource, "/")
	if len(parts) < 3 {
		return "", ""
	}
	repo := path.Join(parts[0], parts[1], parts[2])
	relresource := path.Join(parts[3:]...)
	return repo, relresource
}

var files = map[string]string{}

func retrieveFromMap(repo, resource, version string) (io.Reader, error) {
	contents, ok := files[fmt.Sprintf("%s/%s@%s", repo, resource, version)]
	if !ok {
		return nil, fmt.Errorf("Can't find file %s%s@%s", repo, resource, version)
	}
	return strings.NewReader(contents), nil
}

/* re is the regex that matches the import statements, and  */
func findImports(importRegex string, file []byte) []string {
	var re = regexp.MustCompile(importRegex)
	scanner := bufio.NewScanner(bytes.NewReader(file))
	var imports []string
	for scanner.Scan() {
		for _, match := range re.FindAllStringSubmatch(scanner.Text(), -1) {
			if match == nil {
				continue
			}
			for i, name := range re.SubexpNames() {
				if name == "import" && match[i] != "" {
					if path.Ext(match[i]) != ".sysl" {
						match[i] += ".sysl"
					}
					imports = append(imports, match[i])
				}
			}
		}
	}
	return imports
}
