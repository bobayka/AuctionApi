package services

import (
	"database/sql"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/internal/domains"
	"gitlab.com/bobayka/courseproject/internal/postgres/storage"
	"gitlab.com/bobayka/courseproject/internal/requests"
)

type UserService struct {
	StmtsStorage *storage.UsersStorage
}

func (a *UserService) UpdateUser(u *request.UpdateUser, userID int64) (*domains.User, error) {
	if err := a.StmtsStorage.UpdateUserBD(userID, u); err != nil {
		return nil, err
	}
	db, err := a.StmtsStorage.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (a *UserService) GetUserByID(userID int64) (*domains.User, error) {
	db, err := a.StmtsStorage.FindUserByID(userID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, myerr.ErrNotFound
		}
		return nil, err
	}
	return db, nil
}
