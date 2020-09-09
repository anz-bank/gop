package app

import (
	"io"
	"io/ioutil"
	"path"
	"strings"

	"github.com/joshcarp/gop/gop"
)

type AppConfig struct {
	Username      string `yaml:"username"`
	Token         string `yaml:"token"`
	CacheLocation string `yaml:"cachelocation"`
	FsType        string `yaml:"fstype"` // one of "os", "mem"
	ImportRegex   string `yaml:"importregex"`
	Proxy         string `yaml:"proxy"`
}

func ScanIntoString(res *[]byte, file io.Reader) error {
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	*res = contents
	return nil
}

func ProcessRequest(resource string) (string, string, string, error) {
	location_version := strings.Split(resource, "@")
	if len(location_version) != 2 {
		return "", "", "", CreateError(BadRequestError, "resource must be in form gitx.com/repo/resource.ext@hash")
	}
	repo_resource := location_version[0]
	version := location_version[1]
	parts := strings.Split(repo_resource, "/")
	if len(parts) < 3 {
		return "", "", "", CreateError(BadRequestError, "resource must be in form gitx.com/repo/resource.ext@hash")
	}
	repo := path.Join(parts[0], parts[1], parts[2])
	relresource := path.Join(parts[3:]...)
	return repo, relresource, version, nil
}

func NewObject(resource string) *gop.Object {
	repo, resource, version, err := ProcessRequest(resource)
	if err != nil {
		return nil
	}
	return &gop.Object{
		Repo:     repo,
		Resource: resource,
		Version:  version,
	}
}

func New(repo, resource, version string) gop.Object {
	return gop.Object{
		Repo:     repo,
		Resource: resource,
		Version:  version,
	}
}
