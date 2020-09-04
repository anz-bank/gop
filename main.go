package main

import (
	"bufio"
	"context"
	"io"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

// AppConfig ...
type AppConfig struct {
	StartingBalance int64 `yaml:"startingBalance"`
	username        string
	token           string
}

func main() {
	log.Fatal(pbmod.Serve(context.Background(), LoadService))
}

func LoadService(ctx context.Context, a AppConfig) (*pbmod.ServiceInterface, error) {
	return &pbmod.ServiceInterface{
		GetResource: server{a.getFromGit}.GetResource,
	}, nil
}

type retrieveFile func(repo, resource, version string) (contents []byte, err error)

func (a AppConfig) getFromGit(repo, resource, version string) ([]byte, error) {
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
		URL:           repo,
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
	return ioutil.ReadAll(reader)
}

type server struct {
	retrieveFile
}

func (s server) GetResource(ctx context.Context, req *pbmod.GetResourceRequest, client pbmod.GetResourceClient) (*pbmod.RetrieveResponse, error) {
	repo, resource := processRequest(req.Resource)
	contents, err := s.retrieveFile(repo, resource, req.Version)
	if err != nil {
		return nil, err
	}
	return &pbmod.RetrieveResponse{
		Content: contents,
	}, nil
}

/*
github.com/a/b/file.ext -> (https://github.com/a/b/file.ext, file.ext)
*/
func processRequest(resource string) (string, string) {
	parts := strings.Split(resource, "/")
	if len(parts) < 3 {
		return "", ""
	}
	repo := "https://" + path.Join(parts[0], parts[1], parts[2]) + ".git"
	relresource := path.Join(parts[3:]...)
	return repo, relresource
}

/* re is the regex that matches the import statements, and  */
func findImports(importRegex string, file io.Reader) []string {
	var re = regexp.MustCompile(importRegex)
	scanner := bufio.NewScanner(file)
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
