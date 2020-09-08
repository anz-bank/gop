package app

import (
	"github.com/pkg/errors"
)

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

func CreateError(kind Kind, message string, cause ...error) error {
	var err error
	for _, e := range cause {
		err = errors.Wrap(err, e.Error())
	}
	return Error{Kind: kind, Message: message, Cause: err}
}
