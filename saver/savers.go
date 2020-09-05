package saver

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	retrieve "github.com/joshcarp/pb-mod/config"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

type Saver interface {
	Saver(res *pbmod.Module) (err error)
}

type SaveToFile struct {
	AppConfig retrieve.AppConfig
}

func (a SaveToFile) Saver(res *pbmod.Module) (err error) {
	location := path.Join(a.AppConfig.SaveLocation, fmt.Sprintf("%s/%s@%s", res.Repo, res.Resource, res.Version))
	if err := os.MkdirAll(path.Dir(location), os.ModePerm); err != nil {
		return err
	}
	if err := a.SaveToPbJsonFile(res); err != nil {
		return err
	}
	return ioutil.WriteFile(location, []byte(res.Value), os.ModePerm)
}

func (a SaveToFile) SaveToPbJsonFile(res *pbmod.Module) (err error) {
	location := path.Join(a.AppConfig.SaveLocation, fmt.Sprintf("%s/%s.pb.json@%s", res.Repo, res.Resource, res.Version))
	if err := os.MkdirAll(path.Dir(location), os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(location, []byte(*res.Extra), os.ModePerm)
}
