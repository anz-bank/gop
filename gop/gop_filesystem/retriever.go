package gop_filesystem

import (
	"io/ioutil"
	"path"

	"github.com/joshcarp/gop/app"
)

func (a GOP) Retrieve(resource string) ([]byte, bool, error) {
	file, err := a.fs.Open(path.Join(a.AppConfig.CacheLocation, resource))
	if file == nil {
		return nil, false, app.CreateError(app.CacheAccessError, "Error opening file", err)
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, false, app.CreateError(app.CacheAccessError, "Error opening file", err)
	}
	return b, true, nil
}
