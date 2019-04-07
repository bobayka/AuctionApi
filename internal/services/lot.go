package services

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/internal/domains"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"log"
	"math"
)

type LotServ struct {
	StmtsStorage *postgres.UsersStorage
}

func (ls *LotServ) CreateLot(lot *request.LotToCreateUpdate, userID int64) (*domains.Lot, error) {
	LotID, err := ls.StmtsStorage.InsertLot(userID, lot)
	if err != nil {
		if pqerr, ok := errors.Cause(err).(*pq.Error); ok && pqerr.Code == postgres.CheckViolation {
			return nil, errors.Wrap(myerr.ErrBadRequest, "$end auction date less than update lot date$")
		}
		return nil, err
	}
	dbLot, err := ls.StmtsStorage.FindLotByID(LotID)
	if err != nil {
		return nil, err
	}
	return dbLot, nil
}

func (ls *LotServ) UpdateLot(lot *request.LotToCreateUpdate, lotID int64) (*domains.Lot, error) {
	if err := ls.StmtsStorage.UpdateLotBD(lotID, lot); err != nil {
		if pqerr, ok := errors.Cause(err).(*pq.Error); ok && pqerr.Code == postgres.CheckViolation {
			return nil, errors.Wrap(myerr.ErrBadRequest, "$end auction date less than update lot date$")
		}
		return nil, err
	}
	dbLot, err := ls.StmtsStorage.FindLotByID(lotID)
	if err != nil {
		return nil, err
	}
	return dbLot, nil
}

func (ls *LotServ) GetLotByID(lotID int64) (*domains.Lot, error) {
	dbLot, err := ls.StmtsStorage.FindLotByID(lotID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, myerr.ErrNotFound
		}
		return nil, err
	}
	return dbLot, nil
}

func (ls *LotServ) DeleteLotByID(lotID int64) error {
	if err := ls.StmtsStorage.DeleteLotBD(lotID); err != nil {
		return err
	}
	return nil
}

func (ls *LotServ) GetAllLots(status string) ([]domains.Lot, error) {
	dbLots, err := ls.StmtsStorage.FindAllUserLotsBD(status)
	if err != nil {
		return nil, err
	}
	if len(dbLots) == 0 {
		return nil, myerr.ErrNotFound
	}
	return dbLots, nil
}
func floatIsWhole(f float64) bool {
	const epsilon = 1e-9
	if _, frac := math.Modf(f); frac < epsilon || frac > 1.0-epsilon {
		return true
	}
	return false
}
func checkUpdatePriceConditions(userID int64, dbLot *domains.Lot, price float64) error {
	if userID == dbLot.BuyerID.ID {
		return errors.Wrap(myerr.ErrConflict, "$A lot cannot be purchased by the  current buyer$")
	}
	if userID == dbLot.CreatorID.ID {
		return errors.Wrap(myerr.ErrConflict, "$A lot cannot be purchased by the lot creator$")
	}
	log.Printf("error")
	if dbLot.Status != "active" {
		return errors.Wrap(myerr.ErrConflict, `$The lot must be "active"$`)
	}
	log.Printf("error")
	mult := price / dbLot.PriceStep
	if !floatIsWhole(mult) {
		return errors.Wrap(myerr.ErrConflict, "$Lot price should be a multiple of price_step$")
	}
	log.Printf("error")

	if dbLot.BuyPrice != nil {
		if !(price > *dbLot.BuyPrice) {
			return errors.Wrap(myerr.ErrConflict, "$Price must be greater than current buy_price$")
		}
	}
	log.Printf("error")

	return nil
}
func (ls *LotServ) UpdatePrice(userID int64, lotID int64, price float64) (*domains.Lot, error) {
	dbLot, err := ls.StmtsStorage.FindLotByID(lotID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, myerr.ErrNotFound
		}
		return nil, err
	}
	err = checkUpdatePriceConditions(userID, dbLot, price)
	if err != nil {
		return nil, err
	}
	err = ls.StmtsStorage.UpdateLotPriceBD(userID, lotID, price)
	if err != nil {
		return nil, err
	}
	dbLot, err = ls.StmtsStorage.FindLotByID(lotID)
	if err != nil {
		return nil, err
	}
	return dbLot, nil
}
