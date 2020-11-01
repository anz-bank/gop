package retriever_local

import (
	"fmt"
	"io/ioutil"

	"github.com/joshcarp/gop/pkg/gop"

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
		return nil, false, fmt.Errorf("%s: %w", gop.FileNotFoundError, err)
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", gop.FileReadError, err)
	}
	return b, true, nil
}
