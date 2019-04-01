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
	"strconv"
)

type LotServiceHandler struct {
	lotServ services.LotServ
}

func NewLotServiceHandler(storage *postgres.UsersStorage) *LotServiceHandler {
	return &LotServiceHandler{services.LotServ{StmtsStorage: storage}}
}

func (l *LotServiceHandler) Routes(r *chi.Mux) *chi.Mux {
	r.Group(func(r chi.Router) {

		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.AllowContentType("application/json"))
		r.Use(CheckTokenMiddleware(l.lotServ.StmtsStorage))
		r.Route("/lots", func(r chi.Router) {
			r.Post("/", makeHandler(l.CreateHandler))
			r.Put("/{id:[0-9]*}", makeHandler(l.UpdateHandler))
			r.Get("/{id:[0-9]*}", makeHandler(l.GetHandler))
		})

	})
	return r
}

func (l *LotServiceHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	var lot request.LotToCreateUpdate
	if err := readReqData(r, &lot); err != nil {
		resp, _ := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
		myerr.Error(w, string(resp), http.StatusBadRequest)
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
		resp, _ := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
		myerr.Error(w, string(resp), http.StatusBadRequest)
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
		resp, _ := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
		myerr.Error(w, string(resp), http.StatusBadRequest)
		return nil
	}
	id := chi.URLParam(r, "id")
	lotID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		resp, _ := myerr.ErrMarshal("Wrong Lot ID")
		myerr.Error(w, string(resp), http.StatusBadRequest)
		return nil
	}
	dbLot, err := l.lotServ.UpdateLot(&lot, lotID)
	switch errors.Cause(err) {
	case myerr.NotFound:
		resp, _ := myerr.ErrMarshal("Content by the passed ID could not be found")
		myerr.Error(w, string(resp), http.StatusNotFound)
	case myerr.BadRequest:
		resp, _ := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
		myerr.Error(w, string(resp), http.StatusBadRequest)
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
	id := chi.URLParam(r, "id")
	lotID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		resp, _ := myerr.ErrMarshal("Wrong Lot ID")
		myerr.Error(w, string(resp), http.StatusBadRequest)
		return nil
	}
	dbLot, err := l.lotServ.GetLotByID(lotID)
	switch errors.Cause(err) {
	case myerr.NotFound:
		resp, _ := myerr.ErrMarshal("Content by the passed ID could not be found")
		myerr.Error(w, string(resp), http.StatusNotFound)
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
