package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/services"
	"gitlab.com/bobayka/courseproject/pkg/myerr"
	"net/http"
)

type LotServiceHandler struct {
	lotServ services.LotServ
}

func NewLotServiceHandler(storage *postgres.UsersStorage) *LotServiceHandler {
	return &LotServiceHandler{services.LotServ{StmtsStorage: storage}}
}

func (l *LotServiceHandler) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(CheckTokenMiddleware(l.lotServ.StmtsStorage))

	r.Post("/", makeHandler(l.CreateHandler))
	r.Put("/{id:[0-9]*}", makeHandler(l.UpdateHandler))
	r.Get("/{id:[0-9]*}", makeHandler(l.GetHandler))
	r.Delete("/{id:[0-9]*}", makeHandler(l.DeleteHandler))

	return r
}

func (l *LotServiceHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	var lot request.LotToCreateUpdate
	if err := readReqData(r, &lot); err != nil {
		jsonRespond(w, myerr.GetClientErr(err.Error()), http.StatusBadRequest)
		return nil
	}
	ctx := r.Context()
	id, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return errors.New("r.context doesn't contain user_id")
	}
	dbLot, err := l.lotServ.CreateLot(&lot, id)
	switch errors.Cause(err) {
	case myerr.BadRequest:
		jsonRespond(w, myerr.GetClientErr(err.Error()), http.StatusBadRequest)
	case myerr.Success:
		resp, _ := json.Marshal(dbLot)
		myerr.Error(w, string(resp), http.StatusOK)
	default:
		return errors.Wrap(err, "lot cant be create")
	}
	return nil
}

func (l *LotServiceHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	var lot request.LotToCreateUpdate
	if err := readReqData(r, &lot); err != nil {
		jsonRespond(w, myerr.GetClientErr(err.Error()), http.StatusBadRequest)
		return nil
	}
	lotID, err := getIDURLParam(r)
	if err != nil {
		jsonRespond(w, "Wrong Lot ID", http.StatusBadRequest)
		return nil
	}
	dbLot, err := l.lotServ.UpdateLot(&lot, lotID)
	switch errors.Cause(err) {
	case myerr.NotFound:
		jsonRespond(w, "Content by the passed ID could not be found", http.StatusNotFound)
	case myerr.BadRequest:
		jsonRespond(w, myerr.GetClientErr(err.Error()), http.StatusBadRequest)
	case myerr.Success:
		resp, errM := json.Marshal(dbLot)
		if errM != nil {
			return errors.Wrap(errM, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusOK)
	default:
		return errors.Wrap(err, "lot cant be update")
	}
	return nil
}

func (l *LotServiceHandler) GetHandler(w http.ResponseWriter, r *http.Request) error {
	lotID, err := getIDURLParam(r)
	if err != nil {
		jsonRespond(w, "Wrong Lot ID", http.StatusBadRequest)
		return nil
	}
	dbLot, err := l.lotServ.GetLotByID(lotID)
	switch errors.Cause(err) {
	case myerr.NotFound:
		jsonRespond(w, "Content by the passed ID could not be found", http.StatusNotFound)
	case myerr.Success:
		resp, errM := json.Marshal(dbLot)
		if errM != nil {
			return errors.Wrap(errM, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusOK)
	default:
		return errors.Wrap(err, "lot cant be get")
	}
	return nil
}

func (l *LotServiceHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	lotID, err := getIDURLParam(r)
	if err != nil {
		jsonRespond(w, "Wrong Lot ID", http.StatusBadRequest)
		return nil
	}
	err = l.lotServ.DeleteLotByID(lotID)
	switch errors.Cause(err) {
	case myerr.NotFound:
		jsonRespond(w, "Content by the passed ID could not be found", http.StatusNotFound)
	case myerr.Success:
		w.WriteHeader(http.StatusOK)
	default:
		return errors.Wrap(err, "lot cant be get")
	}
	return nil
}
