package myerr

import (
	"github.com/pkg/errors"
)

var (
	UnprocessableEntity = errors.New("unprocessable entity")
	Created             = errors.New("created")
	Accepted            = errors.New("accepted")
	Unauthorized        = errors.New("unauthorized")
	BadRequest          = errors.New("bad request")
	Success             = errors.New("success")
)
