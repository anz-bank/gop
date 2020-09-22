package gop

type Error int

const (
	UnknownError Error = iota
	BadRequestError
	InternalError
	UnauthorizedError
	TimeoutError
	CacheAccessError
	CacheReadError
	ProxyReadError
	DownstreamError
	CacheWriteError
	FileNotFoundError
	FileReadError
	GitCloneError
	GitCheckoutError
	GithubFetchError
)

func (k Error) Error() string {
	return [...]string{
		"UnknownError",
		"BadRequestError",
		"InternalError",
		"UnauthorizedError",
		"TimeoutError",
		"CacheAccessError",
		"CacheReadError",
		"ProxyReadError",
		"DownstreamError",
		"CacheWriteError",
		"FileNotFoundError",
		"FileReadError",
		"GitCloneError",
		"GitCheckoutError",
		"GithubFetchError"}[k]
}

func (k Error) String() string {
	return k.Error()
}
