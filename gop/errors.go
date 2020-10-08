package gop

type Error int

const (
	UnknownError Error = iota
	BadRequestError
	ApiRateLimitError
	BadReferenceError
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
		"ApiRateLimitError",
		"BadReferenceError",
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

func HandleHTTPStatus(statusCode int) error {
	switch statusCode {
	case 200, 204:
		return nil
	case 400:
		return BadRequestError
	case 401:
		return UnauthorizedError
	case 408:
		return TimeoutError
	case 404:
		return FileNotFoundError
	case 403:
		return ApiRateLimitError
	default:
		return UnknownError
	}
}

func ToStatusCode(err error)(info string, httpCode int){
	switch e := err.(type) {
	case Error:
		info = e.String()
		switch e {
		case BadRequestError:
			httpCode = 400
		case UnauthorizedError:
			httpCode = 401
		case TimeoutError:
			httpCode = 408
		case CacheAccessError, CacheWriteError:
			httpCode = 503
		case CacheReadError, FileNotFoundError:
			httpCode = 404
		default:
			httpCode = 500
		}
	default:
		info, httpCode = UnknownError.String(), 500
	}
	return
}