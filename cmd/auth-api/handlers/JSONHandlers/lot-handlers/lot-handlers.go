package lothandlers

import (
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/cmd/utilities"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/services"
	"net/http"
)

// nolint: gochecknoglobals

type LotServiceHandler struct {
	LotsStatus map[string]bool
	LotServ    services.LotService
}

func NewLotServiceHandler(storage *postgres.UsersStorage) *LotServiceHandler {
	var lotsStatus = map[string]bool{
		"created":  true,
		"active":   true,
		"finished": true,
		"":         true,
	}
	return &LotServiceHandler{LotsStatus: lotsStatus, LotServ: services.LotService{StmtsStorage: storage}}
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
	var lot request.LotToCreateUpdate
	if err := utility.ReadReqData(r, &lot); err != nil {
		return errors.Wrap(err, "cant be read req")
	}
	userID, err := utility.GetTokenUserID(r)
	if err != nil {
		return errors.Wrap(err, "cant get token user id")
	}
	dbLot, err := l.LotServ.CreateLot(&lot, userID)
	if err != nil {
		return errors.Wrap(err, "lot cant be create")
	}
	err = utility.MarshalAndRespondJSON(w, dbLot)
	if err != nil {
		return errors.Wrap(err, "marshal and respondJSON")
	}
	return nil
}

func (l *LotServiceHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error { // сделать вложенную функцию
	var lot request.LotToCreateUpdate
	if err := utility.ReadReqData(r, &lot); err != nil {
		return errors.Wrap(err, "cant be read req")
	}
	lotID, err := utility.GetIDURLParam(r)
	if err != nil {
		return errors.Wrap(err, "Wrong Lot ID")
	}
	dbLot, err := l.LotServ.UpdateLot(&lot, lotID)
	if err != nil {
		return errors.Wrap(err, "lot cant be update")
	}

	err = utility.MarshalAndRespondJSON(w, dbLot)
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
	dbLot, err := l.LotServ.UpdatePrice(userID, lotID, price.Price)
	if err != nil {
		return errors.Wrap(err, "cant update price")
	}
	err = utility.MarshalAndRespondJSON(w, dbLot)
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
	dbLot, err := l.LotServ.GetLotByID(lotID)
	if err != nil {
		return errors.Wrap(err, "lot cant be get")
	}

	err = utility.MarshalAndRespondJSON(w, dbLot)
	if err != nil {
		return errors.Wrap(err, "marshal and respondJSON")
	}
	return nil
}

func (l *LotServiceHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) error {
	lotStat := r.URL.Query().Get("status")
	if !l.LotsStatus[lotStat] {
		return errors.Wrap(myerr.ErrBadRequest, "$Wrong lot status$")
	}
	dbLots, err := l.LotServ.GetAllLots(lotStat)
	if err != nil {
		return errors.Wrap(err, "cant get all lots")
	}
	err = utility.MarshalAndRespondJSON(w, dbLots)
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
	err = l.LotServ.DeleteLotByID(lotID)
	if err != nil {
		return errors.Wrap(err, "lot cant be get")
	}
	w.WriteHeader(http.StatusOK)

	return nil
}
