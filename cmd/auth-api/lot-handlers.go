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

func (l *LotServiceHandler) Routes(r *chi.Mux) *chi.Mux {
	r.Group(func(r chi.Router) {

		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.AllowContentType("application/json"))
		r.Use(CheckTokenMiddleware(l.lotServ.StmtsStorage))
		r.Route("/lots", func(r chi.Router) {
			r.Post("/", makeHandler(l.CreateHandler))
			r.Put("/{id}", makeHandler(l.Updatehandler))
		})

	})
	return r
}

func (l *LotServiceHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	var lot request.LotToCreateUpdate
	if err := readReqData(r, &lot); err != nil {
		resp, errM := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
		if errM != nil {
			return errors.Wrap(err, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusBadRequest)
		return nil
	}
	ctx := r.Context()
	id, ok := ctx.Value("user_id").(int64)
	if !ok {
		return errors.New("r.context doesn't contain user_id")
	}
	dbLot, err := l.lotServ.CreateLot(&lot, id)
	switch errors.Cause(err) {
	case myerr.Unauthorized:
		resp, errM := myerr.ErrMarshal("unauthorized")
		if errM != nil {
			return errors.Wrap(errM, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusUnauthorized)
	case myerr.BadRequest:
		resp, errM := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
		if errM != nil {
			return errors.Wrap(errM, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusBadRequest)
	case myerr.Success:
		resp, errM := json.Marshal(dbLot)
		if errM != nil {
			return errors.Wrap(errM, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusOK)
	default:
		return errors.Wrap(err, "lot cant be create")
	}
	return nil
}

func (l *LotServiceHandler) Updatehandler(w http.ResponseWriter, r *http.Request) error {
	var lot request.LotToCreateUpdate
	if err := readReqData(r, &lot); err != nil {
		resp, errM := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
		if errM != nil {
			return errors.Wrap(err, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusBadRequest)
		return nil
	}
	ctx := r.Context()
	_, ok := ctx.Value("user_id").(int64)
	if !ok {
		return errors.New("r.context doesn't contain user_id")
	}
	return nil
}
