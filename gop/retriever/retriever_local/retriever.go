package retriever_local

import (
	"github.com/joshcarp/gop/app"
	"github.com/joshcarp/gop/gop"
	"github.com/spf13/afero"
)

type Retriever struct {
	fs afero.Fs
}

func New(fs afero.Fs) Retriever {
	return Retriever{
		fs: fs,
	}
}

func (r Retriever) Retrieve(repo, resource, version string) (gop.Object, bool, error) {
	res := app.New(repo, resource, version)
	file, err := r.fs.Open(res.Resource)
	if file == nil {
		return res, false, app.CreateError(app.CacheAccessError, "Error opening file", err)
	}
	return res, true, app.ScanIntoString(&res.Content, file)
}
