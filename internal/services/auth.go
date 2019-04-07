package services

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/internal/domains"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	StmtsStorage *postgres.UsersStorage
}

func (a *Auth) RegisterUser(u *request.RegUser) error {
	err := a.StmtsStorage.AddUser(u)
	if err != nil {
		if pqerr, ok := errors.Cause(err).(*pq.Error); ok && pqerr.Code == postgres.UniqueViolation {
			return errors.Wrap(myerr.ErrConflict, "$email already exists$")
		}
		return errors.Wrap(err, "user can't be add")
	}
	return nil
}

func (a *Auth) AuthorizeUser(u *request.AuthUser) (string, error) {
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

func (a *Auth) UpdateUser(u *request.UpdateUser, userID int64) (*domains.User, error) {
	if err := a.StmtsStorage.UpdateUserBD(userID, u); err != nil {
		return nil, err
	}
	db, err := a.StmtsStorage.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (a *Auth) GetUserByID(userID int64) (*domains.User, error) {
	db, err := a.StmtsStorage.FindUserByID(userID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, myerr.ErrNotFound
		}
		return nil, err
	}
	return db, nil
}

func (a *Auth) GetUserLotsByID(userID int64, lotsType string) ([]domains.Lot, error) {
	dbLots, err := a.StmtsStorage.FindUserLotsBD(userID, lotsType)
	if err != nil {
		return nil, err
	}
	if len(dbLots) == 0 {
		return nil, myerr.ErrNotFound
	}
	return dbLots, nil
}
