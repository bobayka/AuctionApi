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
		r.Group(func(r chi.Router) {
			r.Use(CheckTokenMiddleware(h.authServ.StmtsStorage))
			r.Put("/users/0", makeHandler(h.UpdateHandler))

		})
	})
	return r
}

func (h *AuthHandler) RegistrationHandler(w http.ResponseWriter, r *http.Request) error {

	var user request.RegUser
	if err := readReqAndCheckEmail(r, &user); err != nil {
		resp, errM := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
		if errM != nil {
			return errors.Wrap(errM, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusBadRequest)
		return nil
	}
	err := h.authServ.RegisterUser(&user)
	switch errors.Cause(err) {
	case myerr.Conflict:
		resp, errM := myerr.ErrMarshal("email already exists")
		if errM != nil {
			return errors.Wrap(errM, "marshal error")
		}
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
		resp, errM := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
		if errM != nil {
			return errors.Wrap(errM, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusBadRequest)
		return nil
	}
	token, err := h.authServ.AuthorizeUser(&user)
	switch errors.Cause(err) {
	case myerr.Unauthorized:
		resp, errM := myerr.ErrMarshal("invalid email or password")
		if errM != nil {
			return errors.Wrap(errM, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusConflict)
	case myerr.BadRequest:
		resp, errM := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
		if errM != nil {
			return errors.Wrap(errM, "marshal error")
		}
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
	dbUser, err := h.authServ.UpdateUser(&user, id)
	switch errors.Cause(err) {
	case myerr.Unauthorized:
		resp, errM := myerr.ErrMarshal("unauthorized")
		if errM != nil {
			return errors.Wrap(errM, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusUnauthorized)
	case myerr.Success:
		resp, errM := json.Marshal(dbUser)
		if errM != nil {
			return errors.Wrap(errM, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusOK)
	default:
		return errors.Wrap(err, "user cant be update")
	}
	return nil
}
