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

var validBearer = regexp.MustCompile("Bearer ([ a-z]*)")

func CheckBearer(bearer string) (string, error) {
	const token = 1
	m := validBearer.FindStringSubmatch(bearer)
	if m == nil {
		return "", errors.Wrap(myerr.ErrBadRequest, `$bearer doesn't match patern: "Bearer {token}"$`)
	}
	return m[token], nil
}
