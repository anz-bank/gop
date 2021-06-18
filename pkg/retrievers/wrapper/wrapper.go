package wrapper

import (
	"github.com/anz-bank/gop/pkg/gop"
	"github.com/pkg/errors"
)

type Retriever struct {
	retrievers []gop.Retriever
}

func New(retrievers ...gop.Retriever) Retriever {
	return Retriever{retrievers: retrievers}
}

func (a Retriever) Retrieve(resource string) ([]byte, bool, error) {
	var err error
	for _, retr := range a.retrievers {
		x, y, z := retr.Retrieve(resource)
		if z != nil {
			err = errors.Wrap(err, z.Error())
			continue
		}
		return x, y, z
	}
	return nil, false, err
}
