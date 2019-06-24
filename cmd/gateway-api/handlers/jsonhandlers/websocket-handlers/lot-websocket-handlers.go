package websockethandlers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/Protobuf"
	"gitlab.com/bobayka/courseproject/cmd/gateway-api/handlers/htmlhandlers"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/cmd/utilities"
	"gitlab.com/bobayka/courseproject/internal/responce"
	"google.golang.org/grpc/status"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

const baseHTMLDirectory = "cmd/gateway-api/handlers/jsonhandlers/websocket-handlers/html/"

type LotWSHandler struct {
	Hub      *Hub
	client   lotspb.LotsServiceClient
	templs   htmlhandlers.Templates
	upgrader websocket.Upgrader
}

func NewLotWSHandler(client lotspb.LotsServiceClient) *LotWSHandler {
	templ := htmlhandlers.Templates{
		"getActiveLot": template.Must(template.ParseFiles(baseHTMLDirectory+"buyLotScript.html", baseHTMLDirectory+"getLot.html", baseHTMLDirectory+"base.html")),
		"getAllLots":   template.Must(template.ParseFiles(baseHTMLDirectory+"getAllLots.html", baseHTMLDirectory+"scriptGetAllLots.html", baseHTMLDirectory+"base.html")),
	}

	hub := newHub()
	go hub.run()

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return &LotWSHandler{
		Hub:      hub,
		client:   client,
		templs:   templ,
		upgrader: upgrader}
}

func (l *LotWSHandler) BackGroundFinishEndedLots(d time.Duration) {
	go func() {
		for range time.Tick(d) {
			lots, err := l.client.BackgroundUpdateLots(context.Background(), &lotspb.Empty{})
			if err != nil {
				errStatus, _ := status.FromError(err)
				log.Printf("Background: %s", myerr.ConvGRPCStatusToMyError(errStatus))
			}
			if lots == nil {
				continue
			}
			for _, v := range lots.Lots {
				respLot, err := responce.ConvertGRPCToRespLot(v)
				if err != nil {
					log.Printf("cant convert grpc to resp lot: %s", err)
				}
				res, err := json.Marshal(respLot)
				if err != nil {
					log.Printf("Background: can't marshal message: %s", err)
				}
				l.Hub.Broadcast <- BroadcastWithID{LotID: respLot.ID, Data: res}
				l.Hub.Broadcast <- BroadcastWithID{LotID: AllLotsID, Data: res}
			}

		}
	}()

}

func (l *LotWSHandler) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/{id:[0-9]*}", utility.MakeHandler(l.GetActiveLot))
	r.Get("/", utility.MakeHandler(l.GetAllLots))

	r.Put("/lots/{id:[0-9]*}/updprice", utility.MakeHandler(l.UpdatePriceHandler))
	r.HandleFunc("/ws", utility.MakeHandler(l.WSUpdateLot))
	return r
}

func (l *LotWSHandler) GetActiveLot(w http.ResponseWriter, r *http.Request) error {
	lotID, err := utility.GetIDURLParam(r)
	if err != nil {
		return errors.Wrap(err, "Wrong Lot ID")
	}
	dbLot, err := l.client.GetLotByID(context.Background(), &lotspb.LotID{LotID: lotID})
	if err != nil {
		errStatus, _ := status.FromError(err)
		return errors.Wrap(myerr.ConvGRPCStatusToMyError(errStatus), "lot cant be get")
	}
	if dbLot.Status != "active" {
		return myerr.ErrNotFound
	}
	respLot, err := responce.ConvertGRPCToRespLot(dbLot)
	if err != nil {
		return errors.Wrap(err, "cant convert grpc to resp lot")
	}
	l.templs.RenderTemplate(w, "getActiveLot", "base", respLot)
	return nil
}

func (l *LotWSHandler) GetAllLots(w http.ResponseWriter, r *http.Request) error {
	dbLots, err := l.client.GetAllLots(context.Background(),
		&lotspb.Status{Status: ""})
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
	l.templs.RenderTemplate(w, "getAllLots", "base", respLots)
	return nil
}

func (l *LotWSHandler) UpdatePriceHandler(w http.ResponseWriter, r *http.Request) error {
	lotID, err := utility.GetIDURLParam(r)
	if err != nil {
		return errors.Wrap(err, "Wrong Lot ID")
	}
	userID, err := utility.GetTokenUserID(r)
	if err != nil {
		return errors.Wrap(err, "cant get token user id")
	}
	// добавить здесь условия , а в шаблон не добавлять
	priceStep, err := strconv.ParseFloat(r.PostFormValue("priceStep"), 64)
	if err != nil {
		return errors.Wrap(myerr.ErrBadRequest, "$can't find price step form$")
	}

	dbLot, err := l.client.UpdateLotPrice(context.Background(),
		&lotspb.BuyLot{UserID: userID, LotID: lotID, Price: priceStep, IsWS: true})
	if err != nil {
		errStatus, _ := status.FromError(err)
		return errors.Wrap(myerr.ConvGRPCStatusToMyError(errStatus), "cant update price")
	}
	respLot, err := responce.ConvertGRPCToRespLot(dbLot)
	if err != nil {
		return errors.Wrap(err, "cant convert grpc to resp lot")
	}
	res, err := json.Marshal(respLot)
	if err != nil {
		return errors.Wrap(err, "can't marshal message")
	}
	l.Hub.Broadcast <- BroadcastWithID{LotID: lotID, Data: res}
	return nil
}

func (l *LotWSHandler) WSUpdateLot(w http.ResponseWriter, r *http.Request) error {
	conn, err := l.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errors.Wrap(err, "can't upgrade connection")
	}
	var id int64
	lotID := r.URL.Query().Get("id")
	if lotID == "" {
		id = AllLotsID
	} else {
		id, err = strconv.ParseInt(lotID, 10, 64)
		if err != nil {
			log.Printf("wrong id: %s", lotID)
			return nil
		}
	}
	client := &Client{hub: l.Hub, conn: conn, lotID: id, send: make(chan BroadcastWithID, 256)}
	client.hub.register <- client
	go client.writePump()
	return nil
}
