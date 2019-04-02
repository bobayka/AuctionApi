package services

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/internal/domains"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/pkg/myerr"
)

type LotServ struct {
	StmtsStorage *postgres.UsersStorage
}

func (ls *LotServ) CreateLot(lot *request.LotToCreateUpdate, userID int64) (*domains.Lot, error) {
	LotID, err := ls.StmtsStorage.InsertLot(userID, lot)
	if err != nil {
		if pqerr, ok := errors.Cause(err).(*pq.Error); ok && pqerr.Code == check_violation {
			return nil, errors.Wrap(myerr.BadRequest, "$end auction date less than update lot date$")
		}
		return nil, err
	}
	dbLot, err := ls.StmtsStorage.FindLotByID(LotID)
	if err != nil {
		return nil, err
	}
	return dbLot, myerr.Success
}

func (ls *LotServ) UpdateLot(lot *request.LotToCreateUpdate, lotID int64) (*domains.Lot, error) {
	if err := ls.StmtsStorage.UpdateLotBD(lotID, lot); err != nil {
		if pqerr, ok := errors.Cause(err).(*pq.Error); ok && pqerr.Code == check_violation {
			return nil, errors.Wrap(myerr.BadRequest, "$end auction date less than update lot date$")
		}
		return nil, err
	}
	dbLot, err := ls.StmtsStorage.FindLotByID(lotID)
	if err != nil {
		return nil, err
	}
	return dbLot, myerr.Success
}

func (ls *LotServ) GetLotByID(lotID int64) (*domains.Lot, error) {
	dbLot, err := ls.StmtsStorage.FindLotByID(lotID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, myerr.NotFound
		}
		return nil, err
	}
	return dbLot, myerr.Success
}

func (ls *LotServ) DeleteLotByID(lotID int64) error {
	if err := ls.StmtsStorage.DeleteLotBD(lotID); err != nil {
		return err
	}
	return myerr.Success
}
