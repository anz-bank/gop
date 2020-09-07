package gop

import "github.com/joshcarp/gop/gen/pkg/servers/gop"

type Processor interface {
	Process(pre *gop.Object) (err error)
}

type Retriever interface {
	Retrieve(res *gop.Object) (err error)
}

type Cacher interface {
	Cache(res *gop.Object) (err error)
}

type Gopper interface {
	Retriever
	Cacher
}
