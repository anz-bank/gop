package gop

import "github.com/joshcarp/gop/gen/pkg/servers/gop"

type Retriever interface {
	Retrieve(repo, resource, version string) (res gop.Object, cache bool, err error)
}

type Cacher interface {
	Cache(res gop.Object) (err error)
}

type Gopper interface {
	Retriever
	Cacher
}
