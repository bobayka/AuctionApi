package services

import (
	"database/sql"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/internal/domains"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/responce"
	"gitlab.com/bobayka/courseproject/pkg/floatOperations"
)

type LotService struct {
	StmtsStorage *postgres.UsersStorage
}

func CheckUpdateCreateLotConditions(lot *request.LotToCreateUpdate) error {
	if lot.PriceStep != nil {
		if *lot.PriceStep < 1.0 {
			return errors.Wrap(myerr.ErrBadRequest, "$Price step can be not less than one$")
		}
	}
	if lot.MinPrice < 1 {
		return errors.Wrap(myerr.ErrBadRequest, "$Min_Price can be not less  than one$")
	}
	return nil
}

func (ls *LotService) CreateLot(lot *request.LotToCreateUpdate, userID int64) (*responce.RespLot, error) {
	if err := CheckUpdateCreateLotConditions(lot); err != nil {
		return nil, err
	}
	LotID, err := ls.StmtsStorage.InsertLot(userID, lot)
	if err != nil {
		return nil, err
	}
	dbLot, err := ls.StmtsStorage.FindLotByID(LotID)
	if err != nil {
		return nil, err
	}
	respLot, err := ls.StmtsStorage.ConvertLotToRespLot(dbLot)
	if err != nil {
		return nil, err
	}
	return respLot, nil
}

func (ls *LotService) UpdateLot(lot *request.LotToCreateUpdate, lotID int64) (*responce.RespLot, error) {
	if err := CheckUpdateCreateLotConditions(lot); err != nil {
		return nil, err
	}
	if err := ls.StmtsStorage.UpdateLotBD(lotID, lot); err != nil {
		return nil, err
	}
	dbLot, err := ls.StmtsStorage.FindLotByID(lotID)
	if err != nil {
		return nil, err
	}
	respLot, err := ls.StmtsStorage.ConvertLotToRespLot(dbLot)
	if err != nil {
		return nil, err
	}
	return respLot, nil
}

func (ls *LotService) GetLotByID(lotID int64) (*responce.RespLot, error) {
	dbLot, err := ls.StmtsStorage.FindLotByID(lotID)

	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, myerr.ErrNotFound
		}
		return nil, err
	}
	respLot, err := ls.StmtsStorage.ConvertLotToRespLot(dbLot)
	if err != nil {
		return nil, err
	}
	return respLot, nil
}

func (ls *LotService) DeleteLotByID(lotID int64) error {
	if err := ls.StmtsStorage.DeleteLotBD(lotID); err != nil {
		return err
	}
	return nil
}

func (ls *LotService) GetAllLots(status string) ([]*responce.RespLot, error) {
	dbLots, err := ls.StmtsStorage.FindAllLotsBD(status)
	if err != nil {
		return nil, err
	}
	if len(dbLots) == 0 {
		return nil, myerr.ErrNotFound
	}
	var respLots []*responce.RespLot
	for _, v := range dbLots {
		respLot, err := ls.StmtsStorage.ConvertLotToRespLot(v)
		if err != nil {
			return nil, err
		}
		respLots = append(respLots, respLot)
	}
	return respLots, nil
}

func CheckUpdatePriceConditions(userID int64, dbLot *domains.Lot, price float64) error {
	if dbLot.BuyerID != nil {
		if userID == *dbLot.BuyerID {
			return errors.Wrap(myerr.ErrConflict, "$A lot cannot be purchased by the current buyer$")
		}
	}
	if userID == dbLot.CreatorID {
		return errors.Wrap(myerr.ErrConflict, "$A lot cannot be purchased by the lot creator$")
	}
	if dbLot.Status != "active" {
		return errors.Wrap(myerr.ErrConflict, `$Lot status must be 'active'$`)
	}
	if dbLot.MinPrice > price {
		return errors.Wrapf(myerr.ErrConflict, "$Lot price should be greater than min_price: %.2f$", dbLot.MinPrice)
	}
	mult := price / *dbLot.PriceStep
	if !floatlib.FloatIsWhole(mult) {
		return errors.Wrapf(myerr.ErrConflict, "$Lot price should be a multiple of price_step: %.2f$", *dbLot.PriceStep)
	}
	if dbLot.BuyPrice != nil {
		if !(price > *dbLot.BuyPrice) {
			return errors.Wrapf(myerr.ErrConflict, "$Price must be greater than current buy_price: %.2f$", *dbLot.BuyPrice)
		}
	}
	return nil
}
func (ls *LotService) UpdatePrice(userID int64, lotID int64, price float64) (*responce.RespLot, error) {
	dbLot, err := ls.StmtsStorage.FindLotByID(lotID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, myerr.ErrNotFound
		}
		return nil, err
	}
	err = CheckUpdatePriceConditions(userID, dbLot, price)
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
	respLot, err := ls.StmtsStorage.ConvertLotToRespLot(dbLot)
	if err != nil {
		return nil, err
	}
	return respLot, nil
}
