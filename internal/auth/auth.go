package auth

import (
	"database/sql"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/session"
	"gitlab.com/bobayka/courseproject/internal/user"
	"gitlab.com/bobayka/courseproject/pkg/myerr"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	StmtsStorage *postgres.UsersStorage
}

func (a *AuthService) RegisterUser(u *request.RegUser) error {
	_, err := a.StmtsStorage.FindUserByEmail(u.Email)
	if err == nil {
		return errors.Wrap(myerr.Conflict, "email already exists")
	}
	if errors.Cause(err) != sql.ErrNoRows {
		return errors.Wrap(err, "user can't be checked")
	}

	err = a.StmtsStorage.AddUser(u)
	if err != nil {
		return errors.Wrap(err, "user can't be add")
	}
	return myerr.Created
}

func (a *AuthService) AuthorizeUser(u *request.AuthUser) (string, error) {
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

func (a *AuthService) getValidToken(u request.TokenGetter) (*session.Session, error) {
	ses, err := a.StmtsStorage.FindSessionByToken(u.GetToken())
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, errors.Wrap(myerr.Unauthorized, "invalid token")
		}
		return nil, errors.Wrap(err, "cant check user by email query")
	}
	if !ses.CheckTokenTime() {
		return nil, errors.Wrap(myerr.Unauthorized, "token valid timeout")

	}
	return ses, errors.Wrap(myerr.Success, "Accepted")

}

func (a *AuthService) UpdateUser(u *request.UpdateUser) (*user.User, error) {
	s, err := a.getValidToken(u)
	if errors.Cause(err) != myerr.Success {
		return nil, errors.Wrap(err, "cant get valid token")
	}
	if err = a.StmtsStorage.UpdateUserBD(s.UserID, u); err != nil {
		return nil, err
	}

	db, err := a.StmtsStorage.FindUserByID(s.UserID)
	if err != nil {
		return nil, err
	}
	return db, myerr.Success
}
