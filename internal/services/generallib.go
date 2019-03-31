package services

import (
	"database/sql"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/internal/domains"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/pkg/myerr"
)

const (
	unique_violation      = "23505"
	foreign_key_violation = "23503"
	check_violation       = "23514"
)

func CheckValidToken(token string, store *postgres.UsersStorage) (*domains.Session, error) {
	ses, err := store.FindSessionByToken(token)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, errors.Wrap(myerr.Unauthorized, "invalid token")
		}
		return nil, errors.Wrap(err, "cant check user")
	}
	if !ses.CheckTokenTime() {
		return nil, errors.Wrap(myerr.Unauthorized, "token valid timeout")

	}
	return ses, errors.Wrap(myerr.Success, "Accepted")

}
