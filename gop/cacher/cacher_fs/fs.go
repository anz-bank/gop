package cacher_fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	retrieve "github.com/joshcarp/pb-mod/config"
	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

type Cacher struct {
	AppConfig retrieve.AppConfig
}

func (a Cacher) Cache(res *pbmod.Object) (err error) {
	location := path.Join(a.AppConfig.CacheLocation, fmt.Sprintf("%s/%s@%s", res.Repo, res.Resource, res.Version))
	if err := os.MkdirAll(path.Dir(location), os.ModePerm); err != nil {
		return err
	}
	if err := a.SaveToPbJsonFile(res); err != nil {
		return err
	}
	return ioutil.WriteFile(location, []byte(res.Content), os.ModePerm)
}

func (a Cacher) SaveToPbJsonFile(res *pbmod.Object) (err error) {
	location := path.Join(a.AppConfig.CacheLocation, fmt.Sprintf("%s/%s.pb.json@%s", res.Repo, res.Resource, res.Version))
	if err := os.MkdirAll(path.Dir(location), os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(location, []byte(*res.Processed), os.ModePerm)
}
