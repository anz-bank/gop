package gop

/* Retriever is an interface that returns a Object and if the object should be cached in later steps */
type Retriever interface {
	Retrieve(resource string) (content []byte, cached bool, err error)
}

/* Cacher is an interface that saves the res object to a data source */
type Cacher interface {
	Cache(resource string, content []byte) (err error)
}

/* Resolver is an interface that returns the resolved version of the original resource */
type Resolver interface {
	Resolve(resource string) (resolved string)
	Update(old, new string) error
}

/* Gopper is the composition of both Retriever and Cacher */
type Gopper interface {
	Retriever
	Cacher
}

// Object ...
type Object struct {
	Content  []byte `json:"content"`
	Resource string `json:"resource"`
}
