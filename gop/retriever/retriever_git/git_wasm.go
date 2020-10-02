// +build wasm,js

package retriever_git

import "github.com/joshcarp/gop/gop"

/* Retriever is a stubbed version of the git retriever that returns errors for
running in wasm, this is because the go-git implementation does not work in wasm */
type Retriever struct {
}

func New(map[string]string) Retriever {
	return Retriever{}
}

func (a Retriever) Retrieve(string) ([]byte, bool, error) {
	return nil, false, gop.InternalError
}
