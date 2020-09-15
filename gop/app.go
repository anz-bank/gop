package gop

import (
	"path"
	"strings"
)

type Config struct {
	Username      string `yaml:"username"`
	Token         string `yaml:"token"`
	CacheLocation string `yaml:"cachelocation"`
	FsType        string `yaml:"fstype"` // one of "os", "mem"
	ImportRegex   string `yaml:"importregex"`
	Proxy         string `yaml:"proxy"`
}

func ProcessRequest(resource string) (string, string, string, error) {
	location_version := strings.Split(resource, "@")
	if len(location_version) != 2 {
		return "", resource, "", nil
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
