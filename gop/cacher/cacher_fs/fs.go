package cacher_fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/afero"

	"github.com/joshcarp/gop/app"
	"github.com/joshcarp/gop/gen/pkg/servers/gop"
)

type Cacher struct {
	AppConfig app.AppConfig
	fs        afero.Fs
}

func New(appconfig app.AppConfig) Cacher {
	var fs afero.Fs
	switch appconfig.FsType {
	case "os":
		fs = afero.NewOsFs()
	case "memory", "mem":
		fs = afero.NewMemMapFs()
	}
	return Cacher{
		AppConfig: appconfig,
		fs:        fs,
	}
}

func (a Cacher) Cache(res *gop.Object) (err error) {
	location := path.Join(a.AppConfig.CacheLocation, fmt.Sprintf("%s/%s@%s", res.Repo, res.Resource, res.Version))
	if err := os.MkdirAll(path.Dir(location), os.ModePerm); err != nil {
		return err
	}
	if err := a.SaveToPbJsonFile(res); err != nil {
		return err
	}
	return ioutil.WriteFile(location, []byte(res.Content), os.ModePerm)
}

func (a Cacher) SaveToPbJsonFile(res *gop.Object) (err error) {
	location := path.Join(a.AppConfig.CacheLocation, fmt.Sprintf("%s/%s.pb.json@%s", res.Repo, res.Resource, res.Version))
	if err := os.MkdirAll(path.Dir(location), os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(location, []byte(*res.Processed), os.ModePerm)
}
