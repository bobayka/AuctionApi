package lotweb

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/Protobuf"
	"gitlab.com/bobayka/courseproject/cmd/gateway-api/handlers/htmlhandlers"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/cmd/utilities"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/responce"
	"google.golang.org/grpc/status"
	"html/template"
	"net/http"
)

const baseHTMLDirectory = "cmd/gateway-api/handlers/htmlhandlers/"

// nolint: gochecknoglobals
var lotsStatus = map[string]bool{
	"created":  true,
	"active":   true,
	"finished": true,
	"":         true,
}

//var lotsTypes = map[string]bool{
//	"own":   true,
//	"buyed": true,
//	"":      true,
//}

type WebLotHandler struct {
	client lotspb.LotsServiceClient
	templs htmlhandlers.Templates
}

func NewWebLotHandler(client lotspb.LotsServiceClient) *WebLotHandler {

	templ := htmlhandlers.Templates{
		"getAllLots": template.Must(template.ParseFiles(
			baseHTMLDirectory+"lot-handlers/html/getAllLots.html",
			baseHTMLDirectory+"base.html")),
		"getOneLot": template.Must(template.ParseFiles(
			baseHTMLDirectory+"lot-handlers/html/getOneLot.html",
			baseHTMLDirectory+"base.html")),
	}
	return &WebLotHandler{client: client, templs: templ}
}

func (wb *WebLotHandler) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", utility.MakeHandler(wb.GetAllHandler))
	r.Get("/{id:[0-9]*}", utility.MakeHandler(wb.GetHandler))
	r.Put("/{id:[0-9]*}/buy", utility.MakeHandler(wb.UpdatePriceHandler))

	return r
}

func (wb *WebLotHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) error {
	lotStat := r.URL.Query().Get("status")

	if !lotsStatus[lotStat] {
		return errors.Wrap(myerr.ErrBadRequest, "$Wrong lot status$")
	}

	dbLots, err := wb.client.GetAllLots(context.Background(),
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
	wb.templs.RenderTemplate(w, "getAllLots", "base", respLots)
	return nil
}

func (wb *WebLotHandler) GetHandler(w http.ResponseWriter, r *http.Request) error {
	lotID, err := utility.GetIDURLParam(r)
	if err != nil {
		return errors.Wrap(err, "Wrong Lot ID")
	}
	dbLot, err := wb.client.GetLotByID(context.Background(), &lotspb.LotID{LotID: lotID})
	if err != nil {
		errStatus, _ := status.FromError(err)
		return errors.Wrap(myerr.ConvGRPCStatusToMyError(errStatus), "lot cant be get")
	}
	respLot, err := responce.ConvertGRPCToRespLot(dbLot)
	if err != nil {
		return errors.Wrap(err, "cant convert grpc to resp lot")
	}
	wb.templs.RenderTemplate(w, "getOneLot", "base", respLot)
	return nil
}

func (wb *WebLotHandler) UpdatePriceHandler(w http.ResponseWriter, r *http.Request) error {
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
	dbLot, err := wb.client.UpdateLotPrice(context.Background(),
		&lotspb.BuyLot{UserID: userID, LotID: lotID, Price: price.Price})
	if err != nil {
		errStatus, _ := status.FromError(err)
		return errors.Wrap(myerr.ConvGRPCStatusToMyError(errStatus), "cant update price")
	}
	respLot, err := responce.ConvertGRPCToRespLot(dbLot)
	if err != nil {
		return errors.Wrap(err, "cant convert grpc to resp lot")
	}
	wb.templs.RenderTemplate(w, "getOneLot", "base", respLot)
	return nil
}

func (wb *WebLotHandler) GetUserLotsHandler(w http.ResponseWriter, r *http.Request) error {
	UserID, err := utility.GetUserIDURL(r)
	if err != nil {
		return errors.Wrap(err, "cant get id url param") //после отладки можно убрать
	}
	dbLots, err := wb.client.GetLotsByUserID(context.Background(), &lotspb.UserLots{Id: UserID, Type: "own"})
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
	wb.templs.RenderTemplate(w, "getAllLots", "base", respLots)

	return nil
}
