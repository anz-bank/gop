package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

type saver func(res *pbmod.KeyValue) (err error)

func (a server) saveToFile(res *pbmod.KeyValue) (err error) {
	location := path.Join(a.saveLocation, fmt.Sprintf("%s/%s@%s", res.Repo, res.Resource, res.Version))
	if err := os.MkdirAll(path.Dir(location), os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(location, []byte(res.Value), os.ModePerm)
}

func (a server) saveToPbJsonFile(res *pbmod.KeyValue) (err error) {
	location := path.Join(a.saveLocation, fmt.Sprintf("%s/%s.pb.json@%s", res.Repo, res.Resource, res.Version))
	if err := os.MkdirAll(path.Dir(location), os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(location, []byte(*res.Extra), os.ModePerm)
}

func save(repo, resource, version string, contents []byte) (err error) {
	files[fmt.Sprintf("%s/%s@%s", repo, resource, version)] = string(contents)
	return nil
}
