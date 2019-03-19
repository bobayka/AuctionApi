package myerr

import (
	"github.com/pkg/errors"
)

type AppError struct {
	Err     error
	Message string
	Code    int
}

func (a *AppError) MyWrap(msg string) *AppError {
	a.Err = errors.Wrap(a.Err, msg)
	return a
}
func NewErr(err error, msg string, code int) *AppError {
	return &AppError{Err: err, Message: msg, Code: code}
}
