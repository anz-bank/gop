package gop

/* Retriever is an interface that returns a Object and if the object should be cached in later steps */
type Retriever interface {
	Retrieve(resource string) (content []byte, cached bool, err error)
}

/* Cacher is an interface that saves the res object to a data source */
type Cacher interface {
	Cache(resource string, content []byte) (err error)
}

/* Updater is an interface that returns the resolved version of the original resource */
type Updater interface {
	Resolve(resource string) (resolved string)
	Update(version string) error
	UpdateAll() error
	UpdateTo(from, to string) error
}

type Resolver interface {
	Resolve(resource string) (resolved string, err error)
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

type Logger func(format string, args ...interface{})
