package websocket_handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/auth-api/handlers/HTMLHandlers"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/cmd/utilities"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/responce"
	"gitlab.com/bobayka/courseproject/internal/services"
	"html/template"
	"net/http"
	"strconv"
)

const baseHTMLDirectory = "cmd/auth-api/handlers/JSONHandlers/websocket-handlers/html/"

type LotWSHandler struct {
	lotServ  services.LotService
	wsServ   services.WSService
	templs   HTMLHandlers.Templates
	upgrader websocket.Upgrader
}

func NewLotWSHandler(storage *postgres.UsersStorage) *LotWSHandler {
	templ := HTMLHandlers.Templates{
		"getActiveLots": template.Must(template.ParseFiles(baseHTMLDirectory+"buyLotScript.html", baseHTMLDirectory+"getActiveLots.html", baseHTMLDirectory+"base.html")),
		//"getAllLots":template.Must(template.ParseFiles(
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return &LotWSHandler{lotServ: services.LotService{StmtsStorage: storage},
		wsServ:   services.WSService{StmtsStorage: storage},
		templs:   templ,
		upgrader: upgrader}
}

func (l *LotWSHandler) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/{id:[0-9]*}", utility.MakeHandler(l.GetActiveLot))
	hub := newHub()
	go hub.run()
	r.Put("/lots/{id:[0-9]*}/updprice",
		utility.MakeHandler(func(w http.ResponseWriter, r *http.Request) error {
			err := l.UpdatePriceHandler(hub, w, r)
			return err
		}))
	r.HandleFunc("/ws/{id:[0-9]*}/buy",
		utility.MakeHandler(func(w http.ResponseWriter, r *http.Request) error {
			err := l.WSUpdateLot(hub, w, r)
			return err
		}))
	return r
}

func (l *LotWSHandler) GetActiveLot(w http.ResponseWriter, r *http.Request) error {
	lotID, err := utility.GetIDURLParam(r)
	if err != nil {
		return errors.Wrap(err, "Wrong Lot ID")
	}
	dbLot, err := l.lotServ.GetLotByID(lotID)
	if err != nil {
		return errors.Wrap(err, "lot cant be get")
	}
	if dbLot.Status != "active" {
		return myerr.ErrNotFound
	}
	userID, err := utility.GetTokenUserID(r)
	if err != nil {
		return errors.Wrap(err, "cant get token user id")
	}
	templDBLot := struct {
		Lot    *responce.RespLot
		UserID int64
	}{
		dbLot, userID,
	}
	l.templs.RenderTemplate(w, "getActiveLots", "base", templDBLot)
	return nil
}

func (l *LotWSHandler) UpdatePriceHandler(hub *Hub, w http.ResponseWriter, r *http.Request) error {
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
	dbLot, err := l.wsServ.UpdatePrice(userID, lotID, priceStep)
	if err != nil {
		return errors.Wrap(err, "lot cant be update")
	}
	res, err := json.Marshal(dbLot)
	if err != nil {
		return errors.Wrap(err, "can't marshal message")
	}
	hub.broadcast <- broadcastWithID{lotID: lotID, data: res}
	return nil
}

func (l *LotWSHandler) WSUpdateLot(hub *Hub, w http.ResponseWriter, r *http.Request) error {
	conn, err := l.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errors.Wrap(err, "can't upgrade connection")
	}
	lotID, err := utility.GetIDURLParam(r)
	if err != nil {
		return errors.Wrap(err, "Wrong Lot ID")
	}
	client := &Client{hub: hub, conn: conn, lotID: lotID, send: make(chan broadcastWithID, 256)}
	client.hub.register <- client
	//defer conn.Close()
	go client.writePump()
	return nil

}

//func (l *LotWSHandler) GetAllLotsHandler(w http.ResponseWriter, r *http.Request) error{
//
//}
