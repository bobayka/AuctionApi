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
	"gitlab.com/bobayka/courseproject/internal/services"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

const baseHTMLDirectory = "cmd/auth-api/handlers/JSONHandlers/websocket-handlers/html/"

type LotWSHandler struct {
	Hub      *Hub
	lotServ  services.LotService
	wsServ   services.WSService
	templs   HTMLHandlers.Templates
	upgrader websocket.Upgrader
}

func NewLotWSHandler(storage *postgres.UsersStorage) *LotWSHandler {
	templ := HTMLHandlers.Templates{
		"getActiveLots": template.Must(template.ParseFiles(baseHTMLDirectory+"buyLotScript.html", baseHTMLDirectory+"getLot.html", baseHTMLDirectory+"base.html")),
		"getAllLots":    template.Must(template.ParseFiles(baseHTMLDirectory+"getAllLots.html", baseHTMLDirectory+"scriptGetAllLots.html", baseHTMLDirectory+"base.html")),
	}

	hub := newHub()
	go hub.run()

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return &LotWSHandler{
		Hub:      hub,
		lotServ:  services.LotService{StmtsStorage: storage},
		wsServ:   services.WSService{StmtsStorage: storage},
		templs:   templ,
		upgrader: upgrader}
}

func (l *LotWSHandler) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/{id:[0-9]*}", utility.MakeHandler(l.GetActiveLot))
	r.Get("/", utility.MakeHandler(l.GetAllLots))

	r.Put("/lots/{id:[0-9]*}/updprice", utility.MakeHandler(l.UpdatePriceHandler))
	r.HandleFunc("/websocket", utility.MakeHandler(l.WSUpdateLot))
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
	//userID, err := utility.GetTokenUserID(r)
	//if err != nil {
	//	return errors.Wrap(err, "cant get token user id")
	//}
	//templDBLot := struct {
	//	Lot    *responce.RespLot
	//	UserID int64
	//}{
	//	dbLot, userID,
	//}
	l.templs.RenderTemplate(w, "getActiveLots", "base", dbLot)
	return nil
}

func (l *LotWSHandler) GetAllLots(w http.ResponseWriter, r *http.Request) error {
	dbLots, err := l.lotServ.GetAllLots("")
	if err != nil {
		return errors.Wrap(err, "cant get all lots")
	}
	l.templs.RenderTemplate(w, "getAllLots", "base", dbLots)
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
	dbLot, err := l.wsServ.UpdatePrice(userID, lotID, priceStep)
	if err != nil {
		return errors.Wrap(err, "lot cant be update")
	}
	res, err := json.Marshal(dbLot)
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
		id = -1
	} else {
		id, err = strconv.ParseInt(lotID, 10, 64)
		if err != nil {
			log.Println("Wrong ID: %s", lotID)
			return nil
		}
	}
	client := &Client{hub: l.Hub, conn: conn, lotID: id, send: make(chan BroadcastWithID, 256)}
	client.hub.register <- client
	go client.writePump()
	return nil

}

//func (l *LotWSHandler) GetAllLotsHandler(w http.ResponseWriter, r *http.Request) error{
//
//}
