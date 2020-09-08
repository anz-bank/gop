package gop_filesystem

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/afero"

	"github.com/joshcarp/gop/app"
	"github.com/joshcarp/gop/gen/pkg/servers/gop"
)

type GOP struct {
	AppConfig app.AppConfig
	fs        afero.Fs
}

func (a GOP) Cache(res gop.Object) (err error) {
	location := path.Join(a.AppConfig.CacheLocation, fmt.Sprintf("%s/%s", res.Repo, res.Resource))
	if err := a.fs.MkdirAll(path.Dir(location), os.ModePerm); err != nil {
		return err
	}
	location += "@" + res.Version
	return afero.WriteFile(a.fs, location, []byte(res.Content), os.ModePerm)
}
