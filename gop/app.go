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

func ProcessRepo(resource string) (string, string, string) {
	var version string
	location_version := strings.Split(resource, "@")
	repo_resource := location_version[0]
	if len(location_version) > 1 {
		version = location_version[1]
	}
	parts := strings.Split(repo_resource, "/")
	if len(parts) < 3 {
		return repo_resource, "", version
	}
	repo := path.Join(parts[0], parts[1], parts[2])
	relresource := path.Join(parts[3:]...)
	return repo, relresource, version
}

func ProcessRequest(resource string) (string, string, string, error) {
	var version string
	location_version := strings.Split(resource, "@")
	repo_resource := location_version[0]
	if len(location_version) > 1 {
		version = location_version[1]
	} else {
		return "", location_version[0], "", nil
	}
	parts := strings.Split(repo_resource, "/")
	if len(parts) < 3 {
		return "", "", "", BadRequestError
	}
	repo := path.Join(parts[0], parts[1], parts[2])
	relresource := path.Join(parts[3:]...)
	return repo, relresource, version, nil
}

func CreateResource(repo, resource, version string) string {
	return path.Join(repo, resource) + "@" + version
}
