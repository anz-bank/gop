// Code generated by sysl DO NOT EDIT.
package gop

import (
	"time"

	"github.com/anz-bank/sysl-go/validator"

	"github.com/rickb777/date"
)

// Reference imports to suppress unused errors
var _ = time.Parse

// Reference imports to suppress unused errors
var _ = date.Parse

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

// *Object validator
func (s *Object) Validate() error {
	return validator.Validate(s)
}
