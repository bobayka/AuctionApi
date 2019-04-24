package services

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/Protobuf"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/internal/domains"
	"gitlab.com/bobayka/courseproject/internal/postgres/storage"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/responce"
	"gitlab.com/bobayka/courseproject/pkg/floatOperations"
	"google.golang.org/grpc/status"
	"net/http"
)

type LotService struct {
	LotStmtsStorage *storage.LotsStorage
}

func (ls *LotService) ConvertLotToRespLot(dbLot *domains.Lot) (*lotspb.Lot, error) {

	creator, err := ls.LotStmtsStorage.FindShortUserByID(dbLot.CreatorID)
	if err != nil {
		return nil, err
	}
	var buyer *responce.ShortUser
	if dbLot.BuyerID != nil {
		buyer, err = ls.LotStmtsStorage.FindShortUserByID(*dbLot.BuyerID)
		if err != nil {
			return nil, err
		}
	}
	respLot := &responce.RespLot{LotGeneral: dbLot.LotGeneral, Creator: *creator, Buyer: buyer}
	return responce.ConvertRespLotToGRPC(respLot)
}

func (ls *LotService) BackgroundUpdateLots(ctx context.Context, in *lotspb.Empty) (*lotspb.Lots, error) {

	lotsID, err := ls.LotStmtsStorage.BackgroundUpdateLotsBD()
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, "cant update lots bd")
	}
	var respLotsGRPC []*lotspb.Lot

	for _, v := range lotsID {
		dbLot, err := ls.LotStmtsStorage.FindLotByID(*v)
		if err != nil {
			return nil, err
		}
		respLot, err := ls.ConvertLotToRespLot(dbLot)
		if err != nil {
			return nil, err
		}
		respLotsGRPC = append(respLotsGRPC, respLot)
	}
	return &lotspb.Lots{Lots: respLotsGRPC}, nil
}

func CheckUpdateCreateLotConditions(lot *request.LotCreateUpdate) error {
	if lot.PriceStep != nil {
		if *lot.PriceStep < 1.0 {
			return errors.Wrap(myerr.ErrConflict, "$Price step can be not less than one$")
		}
	}
	if lot.MinPrice < 1 {
		return errors.Wrap(myerr.ErrConflict, "$Min_Price can be not less  than one$")
	}
	if lot.Status != nil {
		if *lot.Status == "finished" {
			return errors.Wrap(myerr.ErrConflict, "$status 'finished' isn't allowed$")
		}
	} else {
		s := "created"
		lot.Status = &s
	}
	return nil
}

func (ls *LotService) CreateLot(ctx context.Context, lot *lotspb.LotCreateUpdate) (*lotspb.Lot, error) {
	lotToCrUp, userID, _, err := request.ConvGRPCToLotCreateUpdate(lot)
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, "can't conv grpc to lot create update")
	}
	var price float64
	if lotToCrUp.PriceStep == nil {
		price = 1
		lotToCrUp.PriceStep = &price
	}
	if err = CheckUpdateCreateLotConditions(lotToCrUp); err != nil {
		return nil, status.Errorf(http.StatusConflict, err.Error())
	}
	LotID, err := ls.LotStmtsStorage.InsertLot(*userID, lotToCrUp)
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, err.Error())
	}
	dbLot, err := ls.LotStmtsStorage.FindLotByID(LotID)
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, err.Error())
	}
	lotGRPS, err := ls.ConvertLotToRespLot(dbLot)
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError,
			errors.Wrap(err, "can't convert lot to respond").Error())
	}
	return lotGRPS, nil
}

func (ls *LotService) UpdateLot(ctx context.Context, lot *lotspb.LotCreateUpdate) (*lotspb.Lot, error) {
	lotToCrUp, UserID, lotID, err := request.ConvGRPCToLotCreateUpdate(lot)
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, "can't conv grpc to lot create update")
	}
	dbLot, err := ls.LotStmtsStorage.FindLotByID(*lotID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, status.Errorf(http.StatusNotFound, "")
		}
		return nil, status.Errorf(http.StatusInternalServerError, err.Error())
	}
	if err = deleteUpdateLotConditions(dbLot, *UserID); err != nil {
		return nil, status.Errorf(http.StatusConflict, err.Error())
	}
	var price float64
	if lotToCrUp.PriceStep == nil {
		price = 1
		lotToCrUp.PriceStep = &price
	}
	if err = CheckUpdateCreateLotConditions(lotToCrUp); err != nil {
		return nil, status.Errorf(http.StatusBadRequest, err.Error())
	}
	if err = ls.LotStmtsStorage.UpdateLotBD(*lotID, lotToCrUp); err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, err.Error())
	}
	dbLot, err = ls.LotStmtsStorage.FindLotByID(*lotID)
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, err.Error())
	}
	lotGRPS, err := ls.ConvertLotToRespLot(dbLot)
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError,
			errors.Wrap(err, "can't convert lot to respond").Error())
	}
	return lotGRPS, nil
}

func (ls *LotService) GetLotByID(ctx context.Context, lotID *lotspb.LotID) (*lotspb.Lot, error) {
	dbLot, err := ls.LotStmtsStorage.FindLotByID(lotID.LotID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, status.Errorf(http.StatusNotFound, "")
		}
		return nil, status.Errorf(http.StatusInternalServerError, errors.Wrap(err, "can't find lot bd").Error())
	}
	lotGRPS, err := ls.ConvertLotToRespLot(dbLot)
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError,
			errors.Wrap(err, "can't convert lot to respond").Error())
	}
	return lotGRPS, nil
}
func deleteUpdateLotConditions(dbLot *domains.Lot, userID int64) error {
	if dbLot.Status != "created" {
		return errors.New("$active or finished lot can't be update/delete$")
	}
	if dbLot.CreatorID != userID {
		return errors.New("$only creator can update lot$")
	}
	if dbLot.IsDeleted() {
		return errors.New("$lot deleted$")
	}
	return nil
}
func (ls *LotService) DeleteLotByID(ctx context.Context, lotID *lotspb.UserLotID) (*lotspb.Empty, error) {
	dbLot, err := ls.LotStmtsStorage.FindLotByID(lotID.LotID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, status.Errorf(http.StatusNotFound, "")
		}
		return nil, status.Errorf(http.StatusInternalServerError, errors.Wrap(err, "can't find lot bd").Error())
	}

	if err := deleteUpdateLotConditions(dbLot, lotID.UserID); err != nil {
		return nil, status.Errorf(http.StatusConflict, err.Error())
	}
	if err := ls.LotStmtsStorage.DeleteLotBD(lotID.LotID); err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, err.Error())
	}
	fmt.Println("fdlfldmf")
	return &lotspb.Empty{}, nil
}

func (ls *LotService) GetAllLots(ctx context.Context, req *lotspb.Status) (*lotspb.Lots, error) {
	dbLots, err := ls.LotStmtsStorage.FindAllLotsBD(req.Status)
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError,
			errors.Wrap(err, "can't find all lots bd").Error())
	}
	if len(dbLots) == 0 {
		return nil, status.Error(http.StatusNotFound, "")
	}
	var respLotsGRPC []*lotspb.Lot

	for _, v := range dbLots {
		lotGRPS, err := ls.ConvertLotToRespLot(v)
		if err != nil {
			return nil, status.Errorf(http.StatusInternalServerError,
				errors.Wrap(err, "can't convert lot to respond").Error())
		}
		respLotsGRPC = append(respLotsGRPC, lotGRPS)
	}
	return &lotspb.Lots{Lots: respLotsGRPC}, nil
}

func (ls *LotService) GetLotsByUserID(ctx context.Context, req *lotspb.UserLots) (*lotspb.Lots, error) {
	dbLots, err := ls.LotStmtsStorage.FindUserLotsBD(req.Id, req.Type)
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError,
			errors.Wrap(err, "can't find user lots bd").Error())
	}
	if len(dbLots) == 0 {
		return nil, status.Error(http.StatusNotFound, "")
	}
	var respLotsGRPC []*lotspb.Lot
	for _, v := range dbLots {
		lotGRPS, err := ls.ConvertLotToRespLot(v)
		if err != nil {
			return nil, status.Errorf(http.StatusInternalServerError,
				errors.Wrap(err, "can't convert lot to respond").Error())
		}
		respLotsGRPC = append(respLotsGRPC, lotGRPS)
	}
	return &lotspb.Lots{Lots: respLotsGRPC}, nil
}

func CheckUpdatePriceConditions(userID int64, dbLot *domains.Lot, price float64) error {
	if dbLot.BuyerID != nil {
		if userID == *dbLot.BuyerID {
			return errors.Wrap(myerr.ErrConflict, "$A lot cannot be purchased by the current buyer$")
		}
	}
	if dbLot.IsDeleted() {
		return errors.Wrap(myerr.ErrConflict, "$lot deleted$")
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
	mult := price / dbLot.PriceStep
	if !floatlib.FloatIsWhole(mult) {
		return errors.Wrapf(myerr.ErrConflict, "$Lot price should be a multiple of price_step: %.2f$", dbLot.PriceStep)
	}
	if dbLot.BuyPrice != nil {
		if !(price > *dbLot.BuyPrice) {
			return errors.Wrapf(myerr.ErrConflict, "$Price must be greater than current buy_price: %.2f$", *dbLot.BuyPrice)
		}
	}
	return nil
}
func (ls *LotService) UpdateLotPrice(ctx context.Context, req *lotspb.BuyLot) (*lotspb.Lot, error) {
	dbLot, err := ls.LotStmtsStorage.FindLotByID(req.LotID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, status.Errorf(http.StatusNotFound, "")
		}
		return nil, status.Errorf(http.StatusInternalServerError, errors.Wrap(err, "can't find lot bd").Error())
	}
	price := req.Price
	if req.IsWS {
		price = setWSPrice(dbLot.BuyPrice, dbLot.MinPrice, req.Price)
	}
	err = CheckUpdatePriceConditions(req.UserID, dbLot, price)
	if err != nil {
		return nil, status.Errorf(http.StatusBadRequest, err.Error())
	}
	err = ls.LotStmtsStorage.UpdateLotPriceBD(req.UserID, req.LotID, price)
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, err.Error())
	}
	dbLot, err = ls.LotStmtsStorage.FindLotByID(req.LotID)
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, err.Error())
	}
	lotGRPS, err := ls.ConvertLotToRespLot(dbLot)
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError,
			errors.Wrap(err, "can't convert lot to respond").Error())
	}
	return lotGRPS, nil
}

func setWSPrice(buyPrice *float64, minPrice float64, priceStep float64) float64 {
	var price float64
	if buyPrice == nil {
		price = minPrice
	} else {
		price = *buyPrice + priceStep
	}
	return price
}
