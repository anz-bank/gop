package gop

import "github.com/joshcarp/gop/gen/pkg/servers/gop"

/* Retriever is an interface that returns a gop.Object and if the object should be cached in later steps */
type Retriever interface {
	Retrieve(repo, resource, version string) (res gop.Object, cached bool, err error)
}

/* Cacher is an interface that saves the res object to a data source */
type Cacher interface {
	Cache(res gop.Object) (err error)
}

/* Gopper is the composition of both Retriever and Cacher */
type Gopper interface {
	Retriever
	Cacher
}
