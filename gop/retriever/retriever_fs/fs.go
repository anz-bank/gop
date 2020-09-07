package retriever_fs

import (
	"fmt"
	"path"

	"github.com/spf13/afero"

	"github.com/joshcarp/gop/app"
	"github.com/joshcarp/gop/gen/pkg/servers/gop"
)

type Retriever struct {
	AppConfig app.AppConfig
	fs        afero.Fs
}

func New(appconfig app.AppConfig) Retriever {
	var fs afero.Fs
	switch appconfig.FsType {
	case "os":
		fs = afero.NewOsFs()
	case "memory", "mem":
		fs = afero.NewMemMapFs()
	}
	return Retriever{
		AppConfig: appconfig,
		fs:        fs,
	}
}

func (a Retriever) Retrieve(res *gop.Object) error {
	file, err := a.fs.Open(path.Join(a.AppConfig.CacheLocation, fmt.Sprintf("%s/%s.pb.json@%s", res.Repo, res.Resource, res.Version)))
	if file == nil {
		return err
	}
	if err := app.ScanIntoString(res.Processed, file); err != nil {
		return err
	}
	return a.RetrieverFile(res)
}

func (a Retriever) RetrieverFile(res *gop.Object) error {
	file, err := a.fs.Open(path.Join(a.AppConfig.CacheLocation, fmt.Sprintf("%s/%s@%s", res.Repo, res.Resource, res.Version)))
	if file == nil {
		return err
	}
	res.Imported = true
	return app.ScanIntoString(&res.Content, file)
}
