package gop_filesystem

import (
	"fmt"
	"path"

	"github.com/joshcarp/gop/app"
	"github.com/joshcarp/gop/gen/pkg/servers/gop"
)

func (a GOP) Retrieve(res *gop.Object) error {
	file, err := a.fs.Open(path.Join(a.AppConfig.CacheLocation, fmt.Sprintf("%s/%s.pb.json@%s", res.Repo, res.Resource, res.Version)))
	if file == nil {
		return err
	}
	if err := app.ScanIntoString(res.Processed, file); err != nil {
		return err
	}
	if path.Ext(res.Resource) == "sysl" {
		return a.RetrieverFile(res)
	}
	return nil
}

func (a GOP) RetrieverFile(res *gop.Object) error {
	file, err := a.fs.Open(path.Join(a.AppConfig.CacheLocation, fmt.Sprintf("%s/%s@%s", res.Repo, res.Resource, res.Version)))
	if file == nil {
		return err
	}
	res.Imported = true
	return app.ScanIntoString(&res.Content, file)
}
