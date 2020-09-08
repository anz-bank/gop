package app

type Error struct {
	Message string
	Kind    Kind
	Cause   error
}

type Kind int

const (
	UnknownError Kind = iota
	BadRequestError
	InternalError
	UnauthorizedError
	TimeoutError
	CacheAccessError
	CacheReadError
	CacheWriteError
	FileNotFoundError
)

func (k Error) String() string {
	return k.Message
}

func (k Error) Error() string {
	return k.Message
}

func CreateError(kind Kind, message string, cause error) error {
	return Error{Kind: kind, Message: message, Cause: cause}
}
