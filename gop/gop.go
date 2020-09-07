package gop

import "github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"

type Processor interface {
	Process(pre *pbmod.Object) (err error)
}

type Retriever interface {
	Retrieve(res *pbmod.Object) (err error)
}

type Cacher interface {
	Cache(res *pbmod.Object) (err error)
}
