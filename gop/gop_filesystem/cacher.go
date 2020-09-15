package gop_filesystem

import (
	"os"
	"path"

	"github.com/spf13/afero"

	"github.com/joshcarp/gop/app"
)

type GOP struct {
	AppConfig app.AppConfig
	fs        afero.Fs
}

func (a GOP) Cache(resource string, content []byte) (err error) {
	location := path.Join(a.AppConfig.CacheLocation, resource)
	if err := a.fs.MkdirAll(path.Dir(location), os.ModePerm); err != nil {
		return err
	}
	return afero.WriteFile(a.fs, location, content, os.ModePerm)
}
