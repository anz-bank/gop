package gop_filesystem

import (
	"fmt"
	"path"

	"github.com/joshcarp/gop/app"
	"github.com/joshcarp/gop/gop"
)

func (a GOP) Retrieve(repo, resource, version string) (gop.Object, bool, error) {
	res := app.New(repo, resource, version)
	file, err := a.fs.Open(path.Join(a.AppConfig.CacheLocation, fmt.Sprintf("%s/%s@%s", res.Repo, res.Resource, res.Version)))
	if file == nil {
		return res, false, app.CreateError(app.CacheAccessError, "Error opening file", err)
	}
	return res, false, app.ScanIntoString(&res.Content, file)
}
