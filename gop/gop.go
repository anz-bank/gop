package gop

/* Retriever is an interface that returns a Object and if the object should be cached in later steps */
type Retriever interface {
	Retrieve(repo, resource, version string) (res Object, cached bool, err error)
}

/* Cacher is an interface that saves the res object to a data source */
type Cacher interface {
	Cache(res Object) (err error)
}

/* Gopper is the composition of both Retriever and Cacher */
type Gopper interface {
	Retriever
	Cacher
}

// Object ...
type Object struct {
	Content  []byte `json:"content"`
	Repo     string `json:"repo"`
	Resource string `json:"resource"`
	Version  string `json:"version"`
}

// GetRequest ...
type GetRequest struct {
	Resource string
}
