package services

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/responce"
	"gitlab.com/bobayka/courseproject/pkg/myerr"
)

type LotServ struct {
	StmtsStorage *postgres.UsersStorage
}

func (ls *LotServ) CreateLot(lot *request.LotToCreateUpdate, userID int64) (*responce.LotToResponce, error) {
	//s, err := CheckValidToken(lot.AccessToken, ls.StmtsStorage)
	//if errors.Cause(err) != myerr.Success {
	//	return nil, errors.Wrap(err, "cant get valid token")
	//}
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
	dbUser, err := ls.StmtsStorage.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	shortUser := &responce.ShortUSer{ID: userID, FirstName: dbUser.FirstName, LastName: dbUser.LastName}

	return &responce.LotToResponce{Lot: dbLot, Creator: shortUser, Buyer: shortUser}, myerr.Success
}
