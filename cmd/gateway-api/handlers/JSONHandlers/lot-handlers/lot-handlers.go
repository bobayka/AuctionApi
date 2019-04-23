package lothandlers

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/Protobuf"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/cmd/utilities"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/responce"
	"google.golang.org/grpc/status"
	"net/http"
)

// nolint: gochecknoglobals
var lotsStatus = map[string]bool{
	"created":  true,
	"active":   true,
	"finished": true,
	"":         true,
}
var lotsTypes = map[string]bool{
	"own":   true,
	"buyed": true,
	"":      true,
}

type LotServiceHandler struct {
	client lotspb.LotsServiceClient
}

func NewLotHandler(client lotspb.LotsServiceClient) *LotServiceHandler {

	return &LotServiceHandler{client: client}
}

func (l *LotServiceHandler) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/", utility.MakeHandler(l.CreateHandler))
	r.Put("/{id:[0-9]*}", utility.MakeHandler(l.UpdateHandler))
	r.Get("/{id:[0-9]*}", utility.MakeHandler(l.GetHandler))
	r.Delete("/{id:[0-9]*}", utility.MakeHandler(l.DeleteHandler))
	r.Get("/", utility.MakeHandler(l.GetAllHandler))
	r.Put("/{id:[0-9]*}/buy", utility.MakeHandler(l.UpdatePriceHandler))

	return r
}

func (l *LotServiceHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	var lot request.LotCreateUpdate
	if err := utility.ReadReqData(r, &lot); err != nil {
		return errors.Wrap(err, "cant be read req")
	}
	userID, err := utility.GetTokenUserID(r)
	if err != nil {
		return errors.Wrap(err, "cant get token user id")
	}
	reqLot, err := request.ConvLotCreateUpdateToGRPC(&lot, &userID, nil)
	if err != nil {
		return errors.Wrap(err, "cant conv lot to create update to grpc")
	}
	dbLot, err := l.client.CreateLot(context.Background(), reqLot)
	if err != nil {
		errStatus, _ := status.FromError(err)
		return errors.Wrap(myerr.ConvGRPCStatusToMyError(errStatus), "lot cant be create")
	}
	respLot, err := responce.ConvertGRPCToRespLot(dbLot)
	if err != nil {
		return errors.Wrap(err, "cant convert grpc to resp lot")
	}
	err = utility.MarshalAndRespondJSON(w, respLot)
	if err != nil {
		return errors.Wrap(err, "marshal and respondJSON")
	}
	return nil
}

func (l *LotServiceHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error { // сделать вложенную функцию
	var lot request.LotCreateUpdate
	if err := utility.ReadReqData(r, &lot); err != nil {
		return errors.Wrap(err, "cant be read req")
	}
	lotID, err := utility.GetIDURLParam(r)
	if err != nil {
		return errors.Wrap(err, "Wrong Lot ID")
	}
	userID, err := utility.GetTokenUserID(r)
	if err != nil {
		return errors.Wrap(err, "cant get token user id")
	}
	reqLot, err := request.ConvLotCreateUpdateToGRPC(&lot, &userID, &lotID)
	if err != nil {
		return errors.Wrap(err, "cant conv lot to create update to grpc")
	}
	dbLot, err := l.client.UpdateLot(context.Background(), reqLot)
	if err != nil {
		errStatus, _ := status.FromError(err)
		return errors.Wrap(myerr.ConvGRPCStatusToMyError(errStatus), "lot cant be update")
	}
	respLot, err := responce.ConvertGRPCToRespLot(dbLot)
	if err != nil {
		return errors.Wrap(err, "cant convert grpc to resp lot")
	}
	err = utility.MarshalAndRespondJSON(w, respLot)
	if err != nil {
		return errors.Wrap(err, "marshal and respondJSON")
	}
	return nil
}

func (l *LotServiceHandler) UpdatePriceHandler(w http.ResponseWriter, r *http.Request) error {
	var price request.Price
	if err := utility.ReadReqData(r, &price); err != nil {
		return errors.Wrap(err, "cant be read req")
	}
	lotID, err := utility.GetIDURLParam(r)
	if err != nil {
		return errors.Wrap(err, "Wrong Lot ID")
	}
	userID, err := utility.GetTokenUserID(r)
	if err != nil {
		return errors.Wrap(err, "cant get token user id")
	}
	dbLot, err := l.client.UpdateLotPrice(context.Background(),
		&lotspb.BuyLot{UserID: userID, LotID: lotID, Price: price.Price})
	if err != nil {
		errStatus, _ := status.FromError(err)
		return errors.Wrap(myerr.ConvGRPCStatusToMyError(errStatus), "cant update price")
	}
	respLot, err := responce.ConvertGRPCToRespLot(dbLot)
	if err != nil {
		return errors.Wrap(err, "cant convert grpc to resp lot")
	}
	err = utility.MarshalAndRespondJSON(w, respLot)
	if err != nil {
		return errors.Wrap(err, "marshal and respondJSON")
	}
	return nil
}

func (l *LotServiceHandler) GetHandler(w http.ResponseWriter, r *http.Request) error {
	lotID, err := utility.GetIDURLParam(r)
	if err != nil {
		return errors.Wrap(err, "Wrong Lot ID")
	}
	dbLot, err := l.client.GetLotByID(context.Background(), &lotspb.LotID{LotID: lotID})
	if err != nil {
		errStatus, _ := status.FromError(err)
		return errors.Wrap(myerr.ConvGRPCStatusToMyError(errStatus), "lot cant be get")
	}
	respLot, err := responce.ConvertGRPCToRespLot(dbLot)
	if err != nil {
		return errors.Wrap(err, "cant convert grpc to resp lot")
	}
	err = utility.MarshalAndRespondJSON(w, respLot)
	if err != nil {
		return errors.Wrap(err, "marshal and respondJSON")
	}
	return nil
}

func (l *LotServiceHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) error {
	lotStat := r.URL.Query().Get("status")
	if !lotsStatus[lotStat] {
		return errors.Wrap(myerr.ErrBadRequest, "$Wrong lot status$")
	}
	dbLots, err := l.client.GetAllLots(context.Background(),
		&lotspb.Status{Status: lotStat})
	if err != nil {
		errStatus, _ := status.FromError(err)
		return errors.Wrap(myerr.ConvGRPCStatusToMyError(errStatus), "cant get all lots")
	}
	var respLots []*responce.RespLot
	for _, v := range dbLots.Lots {
		respLot, err := responce.ConvertGRPCToRespLot(v)
		if err != nil {
			return errors.Wrap(err, "cant convert grpc to resp lot")
		}
		respLots = append(respLots, respLot)
	}
	err = utility.MarshalAndRespondJSON(w, respLots)
	if err != nil {
		return errors.Wrap(err, "marshal and respondJSON")
	}
	return nil
}

func (l *LotServiceHandler) GetUserLotsHandler(w http.ResponseWriter, r *http.Request) error {
	UserID, err := utility.GetUserIDURL(r)
	if err != nil {
		return errors.Wrap(err, "cant get id url param") //после отладки можно убрать
	}
	lotType := r.URL.Query().Get("type")
	if !lotsTypes[lotType] {
		return errors.Wrap(myerr.ErrBadRequest, "$Wrong lot type$")
	}
	dbLots, err := l.client.GetLotsByUserID(context.Background(), &lotspb.UserLots{Id: UserID, Type: lotType})
	if err != nil {
		errStatus, _ := status.FromError(err)
		return errors.Wrap(myerr.ConvGRPCStatusToMyError(errStatus), "lot cant be update")
	}
	var respLots []*responce.RespLot
	for _, v := range dbLots.Lots {
		respLot, err := responce.ConvertGRPCToRespLot(v)
		if err != nil {
			return errors.Wrap(err, "cant convert grpc to resp lot")
		}
		respLots = append(respLots, respLot)
	}
	err = utility.MarshalAndRespondJSON(w, respLots)
	if err != nil {
		return errors.Wrap(err, "marshal and respondJSON")
	}
	return nil
}

func (l *LotServiceHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	lotID, err := utility.GetIDURLParam(r)
	if err != nil {
		return errors.Wrap(err, "Wrong Lot ID")
	}
	userID, err := utility.GetTokenUserID(r)
	if err != nil {
		return errors.Wrap(err, "cant get token user id")
	}
	_, err = l.client.DeleteLotByID(context.Background(), &lotspb.UserLotID{LotID: lotID, UserID: userID})
	if err != nil {
		errStatus, _ := status.FromError(err)
		return errors.Wrap(myerr.ConvGRPCStatusToMyError(errStatus), "lot cant be get")
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
