package services

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	StmtsStorage *postgres.UsersStorage
}

func (a *AuthService) RegisterUser(u *request.RegUser) error {
	err := a.StmtsStorage.AddUser(u)
	if err != nil {
		if pqerr, ok := errors.Cause(err).(*pq.Error); ok && pqerr.Code == postgres.UniqueViolation {
			return errors.Wrap(myerr.ErrConflict, "$email already exists$")
		}
		return errors.Wrap(err, "user can't be add")
	}
	return nil
}

func (a *AuthService) AuthorizeUser(u *request.AuthUser) (string, error) {
	dbUser, err := a.StmtsStorage.FindUserByEmail(u.Email)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return "", errors.Wrap(myerr.ErrUnauthorized, "$invalid email$")
		}
		return "", errors.Wrap(err, "cant check user by email")
	}
	if errС := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(u.Password)); errС != nil {
		if errС == bcrypt.ErrMismatchedHashAndPassword || err == bcrypt.ErrHashTooShort {
			return "", errors.Wrap(myerr.ErrUnauthorized, "$invalid password$")
		}
		return "", errors.Wrapf(myerr.ErrBadRequest, "$password cant be compared$: %s", errС)
	}
	token, err := a.StmtsStorage.AddSession(dbUser)
	if err != nil {
		return "", errors.Wrap(err, "session can't be add")
	}
	return token, nil

}
