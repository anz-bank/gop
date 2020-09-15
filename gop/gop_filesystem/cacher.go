package gop_filesystem

import (
	"os"
	"path"

	"github.com/spf13/afero"
)

type GOP struct {
	fs  afero.Fs
	dir string
}

func (a GOP) Cache(resource string, content []byte) (err error) {
	location := path.Join(a.dir, resource)
	if err := a.fs.MkdirAll(path.Dir(location), os.ModePerm); err != nil {
		return err
	}
	return afero.WriteFile(a.fs, location, content, os.ModePerm)
}
