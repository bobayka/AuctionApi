package lotWeb

import (
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/auth-api/handlers/HTMLHandlers"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/cmd/utilities"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	request "gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/services"
	"html/template"
	"net/http"
)

const baseHTMLDirectory = "cmd/auth-api/handlers/HTMLHandlers/"

type WebLotHandler struct {
	LotsStatus map[string]bool
	webLot     services.LotService
	templs     HTMLHandlers.Templates
}

func NewWebLotHandler(storage *postgres.UsersStorage) *WebLotHandler {
	var lotsStatus = map[string]bool{
		"created":  true,
		"active":   true,
		"finished": true,
		"":         true,
	}
	templ := HTMLHandlers.Templates{
		"getAllLots": template.Must(template.ParseFiles(baseHTMLDirectory+"lot-handlers/html/getAllLots.html", baseHTMLDirectory+"base.html")),
		"getOneLot":  template.Must(template.ParseFiles(baseHTMLDirectory+"lot-handlers/html/getOneLot.html", baseHTMLDirectory+"base.html")),
	}
	return &WebLotHandler{LotsStatus: lotsStatus, webLot: services.LotService{StmtsStorage: storage}, templs: templ}
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

	if !wb.LotsStatus[lotStat] {
		return errors.Wrap(myerr.ErrBadRequest, "$Wrong lot status$")
	}

	dbLots, err := wb.webLot.GetAllLots(lotStat)
	if err != nil {
		return errors.Wrap(err, "cant get all lots")
	}

	wb.templs.RenderTemplate(w, "getAllLots", "base", dbLots)
	return nil
}

func (wb *WebLotHandler) GetHandler(w http.ResponseWriter, r *http.Request) error {
	lotID, err := utility.GetIDURLParam(r)
	if err != nil {
		return errors.Wrap(err, "Wrong Lot ID")
	}
	dbLot, err := wb.webLot.GetLotByID(lotID)
	if err != nil {
		return errors.Wrap(err, "lot cant be get")
	}
	wb.templs.RenderTemplate(w, "getOneLot", "base", dbLot)
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
	dbLot, err := wb.webLot.UpdatePrice(userID, lotID, price.Price)
	if err != nil {
		return errors.Wrap(err, "cant update price")
	}
	wb.templs.RenderTemplate(w, "getOneLot", "base", dbLot)
	return nil
}
