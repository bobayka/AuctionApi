package services

import (
	"database/sql"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/internal/domains"
	"gitlab.com/bobayka/courseproject/internal/postgres"
)

func CheckValidToken(token string, store *postgres.UsersStorage) (*domains.Session, error) {
	ses, err := store.FindSessionByToken(token)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, errors.Wrap(myerr.ErrUnauthorized, "$invalid token$")
		}
		return nil, errors.Wrap(err, "cant check user")
	}
	if !ses.CheckTokenTime() {
		return nil, errors.Wrap(myerr.ErrUnauthorized, "$token valid timeout$")

	}
	return ses, nil

}
