package services

import (
	"database/sql"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/internal/domains"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	request "gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/responce"
)

type UserService struct {
	StmtsStorage *postgres.UsersStorage
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

func (a *UserService) GetUserLotsByID(userID int64, lotsType string) ([]*responce.RespLot, error) {
	dbLots, err := a.StmtsStorage.FindUserLotsBD(userID, lotsType)
	if err != nil {
		return nil, err
	}
	if len(dbLots) == 0 {
		return nil, myerr.ErrNotFound
	}
	var respLots []*responce.RespLot
	for _, v := range dbLots {
		respLot, err := a.StmtsStorage.ConvertLotToRespLot(v)
		if err != nil {
			return nil, err
		}
		respLots = append(respLots, respLot)
	}
	return respLots, nil
}
