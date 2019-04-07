package myerr

import (
	"github.com/pkg/errors"
	"net/http"
)

func ErrHandler(w http.ResponseWriter, err error) error {
	switch errors.Cause(err) {
	case ErrConflict:
		JSONErrRespond(w, GiveErrToClient(err.Error()), http.StatusConflict)
	case ErrNotFound:
		JSONErrRespond(w, "Content by the passed ID could not be found", http.StatusNotFound)
	case ErrBadRequest:
		JSONErrRespond(w, GiveErrToClient(err.Error()), http.StatusBadRequest)
	case ErrUnauthorized:
		JSONErrRespond(w, GiveErrToClient(err.Error()), http.StatusUnauthorized)
	default:
		return err
	}
	return nil
}
