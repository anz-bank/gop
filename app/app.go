package app

import (
	"io"
	"io/ioutil"
	"path"
	"strings"

	"github.com/joshcarp/gop/gen/pkg/servers/gop"
)

type AppConfig struct {
	Username      string `yaml:"username"`
	Token         string `yaml:"token"`
	CacheLocation string `yaml:"cachelocation"`
	FsType        string `yaml:"fstype"` // one of "os", "mem"
	ImportRegex   string `yaml:"importregex"`
}

func ScanIntoString(res *string, file io.Reader) error {
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fin := string(contents)
	*res = fin
	return nil
}

func ProcessRequest(resource string) (string, string) {
	parts := strings.Split(resource, "/")
	if len(parts) < 3 {
		return "", ""
	}
	repo := path.Join(parts[0], parts[1], parts[2])
	relresource := path.Join(parts[3:]...)
	return repo, relresource
}

func NewObject(resource, version string) *gop.Object {
	var a string
	repo, resource := ProcessRequest(resource)
	return &gop.Object{
		Repo:      repo,
		Resource:  resource,
		Version:   version,
		Processed: &a,
		Content:   "",
	}
}
