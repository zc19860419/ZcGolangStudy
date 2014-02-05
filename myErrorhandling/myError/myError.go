package myError

import (
	"errors"
)

// type error interface {
// 	Error() string
// }
type errorString struct {
	Op  string
	Err error
}

// New returns an error that formats as the given text.
func New(text, reason string) error {
	return &errorString{text, errors.New(reason)}
}

func (a *errorString) Error() string {
	return a.Op + ": " + a.Err.Error()
}
