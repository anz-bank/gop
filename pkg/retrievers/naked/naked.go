package naked

import "github.com/anz-bank/gop/pkg/gop"

type Retriever struct {
	retriever  gop.Retriever
	defaultRef string
}

func New(retriever gop.Retriever, defaultRef string) Retriever {
	return Retriever{retriever: retriever, defaultRef: defaultRef}
}

func (a Retriever) Retrieve(resource string) ([]byte, bool, error) {
	if _, _, ver, _ := gop.ProcessRequest(resource); ver == ""{
		resource += "@"+a.defaultRef
	}
	return a.retriever.Retrieve(resource)
}