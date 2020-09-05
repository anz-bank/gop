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

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/anz-bank/sysl/pkg/parse"

	"github.com/go-git/go-git/v5/plumbing"

	"github.com/go-git/go-git/v5"
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
	return &pbmod.ServiceInterface{
		GetResourceList: server{retrievers: []retriever{a.retrieveSyslPb, a.retrieveGit}, savers: []saver{a.saveToFile, a.saveToPbJsonFile}, posts: []post{processSysl}}.GetResource,
	}, nil
}

type retriever func(res *pbmod.KeyValue) (err error)
type saver func(res *pbmod.KeyValue) (err error)
type post func(pre *pbmod.KeyValue) (err error)

func (a AppConfig) saveToFile(res *pbmod.KeyValue) (err error) {
	location := path.Join(a.saveLocation, fmt.Sprintf("%s/%s@%s", res.Repo, res.Resource, res.Version))
	if err := os.MkdirAll(path.Dir(location), os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(location, []byte(res.Value), os.ModePerm)
}

func (a AppConfig) saveToPbJsonFile(res *pbmod.KeyValue) (err error) {
	location := path.Join(a.saveLocation, fmt.Sprintf("%s/%s.pb.json@%s", res.Repo, res.Resource, res.Version))
	if err := os.MkdirAll(path.Dir(location), os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(location, []byte(*res.Extra), os.ModePerm)
}

func retrieveFie(res *string, file io.Reader) error {
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fin := string(contents)
	*res = fin
	return nil
}

func (a AppConfig) retrieveSyslPB(res *pbmod.KeyValue) error {
	file, err := os.Open(path.Join(a.saveLocation, fmt.Sprintf("%s/%s@%s", res.Repo, res.Resource, res.Version)))
	if file == nil {
		return err
	}
	res.Imported = true
	return retrieveFie(res.Extra, file)
}

func (a AppConfig) retrieveFile(res *pbmod.KeyValue) error {
	file, err := os.Open(path.Join(a.saveLocation, fmt.Sprintf("%s/%s@%s", res.Repo, res.Resource, res.Version)))
	if file == nil {
		return err
	}
	res.Imported = true
	return retrieveFie(&res.Value, file)
}

func (a AppConfig) retrieveSyslPb(res *pbmod.KeyValue) error {
	file, err := os.Open(path.Join(a.saveLocation, fmt.Sprintf("%s/%s.pb.json@%s", res.Repo, res.Resource, res.Version)))
	if file == nil {
		return err
	}
	if err := retrieveFie(&res.Value, file); err != nil {
		return err
	}
	return a.retrieveFile(res)
}

func processSysl(a *pbmod.KeyValue) error {
	if *a.Extra != "" {
		return nil
	}
	m, err := parse.NewParser().ParseString(a.Value)
	if err != nil {
		return err
	}
	ma := protojson.MarshalOptions{Multiline: false, Indent: " ", EmitUnpopulated: false}
	mb, err := ma.Marshal(m)
	if err != nil {
		return err
	}
	extra := string(mb)
	a.Extra = &extra
	return nil
}

func (a AppConfig) retrieveGit(res *pbmod.KeyValue) error {
	var auth *http.BasicAuth
	store := memory.NewStorage()
	if a.username != "" {
		auth = &http.BasicAuth{
			Username: a.username,
			Password: a.token,
		}
	}
	r, err := git.Clone(store, nil, &git.CloneOptions{
		URL:   "https://" + res.Repo + ".git",
		Depth: 1,
		Auth:  auth,
	})
	if err != nil {
		return err
	}
	commit, err := r.CommitObject(plumbing.NewHash(res.Version))
	if err != nil {
		return err
	}
	f, err := commit.File(res.Resource)
	if err != nil {
		return err
	}
	reader, err := f.Reader()
	if err != nil {
		return err
	}
	return retrieveFie(&res.Value, reader)
}

type server struct {
	retrievers []retriever
	savers     []saver
	posts      []post
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

func (s server) GetResource(ctx context.Context, req *pbmod.GetResourceListRequest, client pbmod.GetResourceListClient) (*pbmod.KeyValue, error) {
	repo, resource := processRequest(req.Resource)
	files, err := importFile(repo, resource, req.Version, s.retrievers)
	if !files.Imported {
		for _, p := range s.posts {
			p(files)
		}
		for _, s := range s.savers {
			if err := s(files); err != nil {
				fmt.Println(err)
			}
		}
	}
	return files, err
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
	files[fmt.Sprintf("%s/%s@%s", repo, resource, version)] = string(contents)
	return nil
}

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
