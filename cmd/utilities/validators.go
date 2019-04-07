package utility

import (
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"regexp"
)

// nolint: gochecknoglobals
var validEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func CheckEmail(email string) error {
	if validEmail.FindStringSubmatch(email) == nil {
		return errors.Wrap(myerr.ErrBadRequest, "$email doesnt match pattern$")
	}
	return nil
}
