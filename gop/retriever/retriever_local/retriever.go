package retriever_local

import (
	"io/ioutil"

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

func (r Retriever) Retrieve(resource string) ([]byte, bool, error) {
	file, err := r.fs.Open(resource)
	if file == nil {
		return nil, false, gop.CreateError(gop.CacheAccessError, "Error opening file", err)
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, false, gop.CreateError(gop.CacheAccessError, "Error opening file", err)
	}
	return b, true, nil
}
