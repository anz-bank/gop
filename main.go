package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

var syslimportRegex = `(?:#import.*)|(?:import )(?:\/\/)?(?P<import>.*)`

// AppConfig ...
type AppConfig struct {
	StartingBalance int64 `yaml:"startingBalance"`
	username        string
	token           string
	saveLocation    string
}

func main() {
	log.Fatal(pbmod.Serve(context.Background(), LoadService))

}

func LoadService(ctx context.Context, a AppConfig) (*pbmod.ServiceInterface, error) {
	return &pbmod.ServiceInterface{
		GetResource: server{retriever: []retriever{a.retrieveFie, a.retrieveGit}, saver: a.saveToFile}.GetResource,
	}, nil
}

type retriever func(repo, resource, version string) (contents io.Reader, err error)
type saver func(repo, resource, version string, contents []byte) (err error)

func (a AppConfig) saveToFile(repo, resource, version string, contents []byte) (err error) {
	location := path.Join(a.saveLocation, key(repo, resource, version))
	if err := os.MkdirAll(path.Dir(location), os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(location, contents, os.ModePerm)
}

func (a AppConfig) retrieveFie(repo, resource, version string) (io.Reader, error) {
	return os.Open(path.Join(a.saveLocation, key(repo, resource, version)))
}

func (a AppConfig) retrieveGit(repo, resource, version string) (io.Reader, error) {
	var auth *http.BasicAuth
	store := memory.NewStorage()
	if a.username != "" {
		auth = &http.BasicAuth{
			Username: a.username,
			Password: a.token,
		}
	}
	r, err := git.Clone(store, nil, &git.CloneOptions{
		ReferenceName: plumbing.ReferenceName(version),
		URL:           "https://" + repo + ".git",
		Depth:         1,
		Auth:          auth,
	})
	if err != nil {
		return nil, err
	}
	ref, err := r.Head()
	if err != nil {
		return nil, err
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}
	f, err := commit.File(resource)
	if err != nil {
		return nil, err
	}
	reader, err := f.Reader()
	if err != nil {
		return nil, err
	}
	return reader, nil
}

type server struct {
	retriever []retriever
	saver
}

func doImport(initialrepo, initialImport, ver string, saver saver, retrievers ...retriever) (map[string]string, error) {
	var final = make(map[string]string)
	var err error
	var imports = []string{initialImport}
	var file io.Reader
	for {
		var newImports []string
		for _, imp := range imports {
			for _, r := range retrievers {
				file, err = r(initialrepo, imp, ver)
				if file != nil && err == nil {
					break
				}
			}
			contents, err := ioutil.ReadAll(file)
			if err != nil {
				return nil, err
			}
			if err := saver(initialrepo, imp, ver, contents); err != nil {
				return nil, err
			}
			newImports = append(newImports, findImports(syslimportRegex, contents)...)
			final[key(initialrepo, imp, ver)] = string(contents)
		}
		imports = newImports
		if len(imports) == 0 {
			break
		}
	}
	return final, nil
}

func key(repo, resource, version string) string {
	if version == "" {
		return fmt.Sprintf("%s/%s", repo, resource)
	}
	return fmt.Sprintf("%s/%s@%s", repo, resource, version)
}
func (s server) GetResource(ctx context.Context, req *pbmod.GetResourceRequest, client pbmod.GetResourceClient) (*pbmod.RetrieveResponse, error) {
	repo, resource := processRequest(req.Resource)
	files, err := doImport(repo, resource, req.Version, s.saver, s.retriever...)
	contents := make([]pbmod.KeyValue, 0, len(files))
	for imp, file := range files {
		contents = append(contents, pbmod.KeyValue{Key: imp, Value: file})
	}
	return &pbmod.RetrieveResponse{Content: contents}, err
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

func save(repo, resource, version string, contents []byte) (err error) {
	files[key(repo, resource, version)] = string(contents)
	return nil
}

func retrieveFromMap(repo, resource, version string) (io.Reader, error) {
	contents, ok := files[key(repo, resource, version)]
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
					imports = append(imports, match[i])
				}
			}
		}
	}
	return imports
}

/*
CREATE TABLE Modules{
location VARCHAR(),
content BLOB(), // *sysl.Module ->Json || sysl
version varchar,
}
*/

/*
Client -> Proxy (*sysl.Module A)

Proxy:
	Module A -> Module B
	return [Module A, Module B] pb.json


Module A:
	a.sysl:```
import //github.com/b/c/b.sysl
import d.sysl
```
	a.pb.json

1. Clone repo A
2. Open file A
3. Save file a,  parsed file a, path, repo, version
4. Walk dependency graph of A, repeat step 1-4 until no more dependencies
5. return array of modules


	func Parse([]*sysl.Module) *sysl.Module
*/

/*
Local file:
import //a.com/b/c/d.sysl

func Parse([]*sysl.Module) *sysl.Module{

}

*/
