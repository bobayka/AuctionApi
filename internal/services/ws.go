package services

import (
	"database/sql"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/responce"
)

type WSService struct {
	StmtsStorage *postgres.UsersStorage
}

func (ws *WSService) UpdatePrice(userID int64, lotID int64, priceStep float64) (*responce.RespLot, error) {
	dbLot, err := ws.StmtsStorage.FindLotByID(lotID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, myerr.ErrNotFound
		}
		return nil, err
	}
	var price float64
	if dbLot.BuyPrice == nil {
		price = dbLot.MinPrice
	} else {
		price = *dbLot.BuyPrice + priceStep
	}
	err = CheckUpdatePriceConditions(userID, dbLot, price)
	if err != nil {
		return nil, err
	}
	err = ws.StmtsStorage.UpdateLotPriceBD(userID, lotID, price)
	if err != nil {
		return nil, err
	}
	dbLot.BuyPrice = &price
	respLot, err := ws.StmtsStorage.ConvertLotToRespLot(dbLot)
	if err != nil {
		return nil, err
	}
	return respLot, nil
}
