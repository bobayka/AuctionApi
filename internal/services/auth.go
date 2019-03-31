package services

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/internal/domains"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/pkg/myerr"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	StmtsStorage *postgres.UsersStorage
}

func (a *Auth) RegisterUser(u *request.RegUser) error {
	err := a.StmtsStorage.AddUser(u)
	if err != nil {
		if pqerr, ok := errors.Cause(err).(*pq.Error); ok && pqerr.Code == unique_violation {
			return errors.Wrap(myerr.Conflict, "email already exists")
		}
		return errors.Wrap(err, "user can't be add")
	}
	return myerr.Created
}

func (a *Auth) AuthorizeUser(u *request.AuthUser) (string, error) {
	dbUser, err := a.StmtsStorage.FindUserByEmail(u.Email)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return "", errors.Wrap(myerr.Unauthorized, "invalid email")
		}
		return "", errors.Wrap(err, "cant check user by email")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(u.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword || err == bcrypt.ErrHashTooShort {
			return "", errors.Wrap(myerr.Unauthorized, "invalid password")
		}
		return "", errors.Wrapf(myerr.BadRequest, "$password cant be compared$: %s", err)
	}
	token, err := a.StmtsStorage.AddSession(dbUser)
	if err != nil {
		return "", errors.Wrap(err, "session can't be add")
	}
	return token, myerr.Success

}

func (a *Auth) UpdateUser(u *request.UpdateUser, userID int64) (*domains.User, error) {
	if err := a.StmtsStorage.UpdateUserBD(userID, u); err != nil {
		return nil, err
	}

	db, err := a.StmtsStorage.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	return db, myerr.Success
}
