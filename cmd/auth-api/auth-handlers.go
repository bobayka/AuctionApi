package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/internal/domains"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/services"
	"gitlab.com/bobayka/courseproject/pkg/myerr"
	"net/http"
	"strconv"
)

func NewAuthHandler(storage *postgres.UsersStorage) *AuthHandler {
	return &AuthHandler{services.Auth{StmtsStorage: storage}}
}

type AuthHandler struct {
	authServ services.Auth
}

func (h *AuthHandler) Routes(r *chi.Mux) *chi.Mux {
	//r.Use(middleware.RealIP)
	//r.Use(middleware.RequestID)
	r.Group(func(r chi.Router) {

		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/signup", makeHandler(h.RegistrationHandler))
		r.Post("/signin", makeHandler(h.AuthorizationHandler))
		r.Route("/users", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(CheckTokenMiddleware(h.authServ.StmtsStorage))
				r.Put("/0", makeHandler(h.UpdateHandler))
				r.Get("/{id:[0-9]*}", makeHandler(h.GetHandler))
			})
		})

	})
	return r
}

func (h *AuthHandler) RegistrationHandler(w http.ResponseWriter, r *http.Request) error {

	var user request.RegUser
	if err := readReqAndCheckEmail(r, &user); err != nil {
		resp, _ := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
		myerr.Error(w, string(resp), http.StatusBadRequest)
		return nil
	}
	err := h.authServ.RegisterUser(&user)
	switch errors.Cause(err) {
	case myerr.Conflict:
		resp, _ := myerr.ErrMarshal("email already exists")
		myerr.Error(w, string(resp), http.StatusConflict)
	case myerr.Created:
		w.WriteHeader(http.StatusCreated)
	default:
		return errors.Wrap(err, "user cant be registered")
	}
	return nil
}

func (h *AuthHandler) AuthorizationHandler(w http.ResponseWriter, r *http.Request) error {
	var user request.AuthUser
	if err := readReqAndCheckEmail(r, &user); err != nil {
		resp, _ := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
		myerr.Error(w, string(resp), http.StatusBadRequest)
		return nil
	}
	token, err := h.authServ.AuthorizeUser(&user)
	switch errors.Cause(err) {
	case myerr.Unauthorized:
		resp, _ := myerr.ErrMarshal("invalid email or password")
		myerr.Error(w, string(resp), http.StatusConflict)
	case myerr.BadRequest:
		resp, _ := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
		myerr.Error(w, string(resp), http.StatusBadRequest)
	case myerr.Success:
		resp := `{"token_type": "bearer","access_token": "` + token + `"}`
		myerr.Error(w, resp, http.StatusOK)
	default:
		return errors.Wrap(err, "cant authorize user")
	}
	return nil
}

func (h *AuthHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	var user request.UpdateUser

	if err := readReqData(r, &user); err != nil {
		resp, _ := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
		myerr.Error(w, string(resp), http.StatusBadRequest)
		return nil
	}
	ctx := r.Context()
	id, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return errors.New("r.context doesn't contain user_id")
	}
	dbUser, err := h.authServ.UpdateUser(&user, id)
	switch errors.Cause(err) {
	case myerr.Unauthorized:
		resp, _ := myerr.ErrMarshal("unauthorized")
		myerr.Error(w, string(resp), http.StatusUnauthorized)
	case myerr.Success:
		resp, _ := json.Marshal(dbUser)
		myerr.Error(w, string(resp), http.StatusOK)
	default:
		return errors.Wrap(err, "user cant be update")
	}
	return nil
}

func (h *AuthHandler) GetHandler(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	UserID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		resp, _ := myerr.ErrMarshal("Wrong User ID")
		myerr.Error(w, string(resp), http.StatusBadRequest)
		return nil
	}
	ctx := r.Context()
	OwnId, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return errors.New("r.context doesn't contain user_id")
	}
	var dbUser *domains.User
	if UserID == 0 {
		dbUser, err = h.authServ.GetUserByID(OwnId)
	} else {
		dbUser, err = h.authServ.GetUserByID(UserID)
	}
	switch errors.Cause(err) {
	case myerr.NotFound:
		resp, _ := myerr.ErrMarshal("Content by the passed ID could not be found")
		myerr.Error(w, string(resp), http.StatusNotFound)
	case myerr.Success:
		resp, errM := json.Marshal(dbUser)
		if errM != nil {
			return errors.Wrap(errM, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusOK)
	default:
		return errors.Wrap(err, "lot cant be get")
	}
	return nil
}
