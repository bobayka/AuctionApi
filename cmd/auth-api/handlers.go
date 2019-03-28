package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	myAuth "gitlab.com/bobayka/courseproject/internal/auth"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/pkg/myerr"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var validEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var Bearer = regexp.MustCompile("(b|B)earer:([ a-z]*)")

func makeHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func NewAuthHandler(storage *postgres.UsersStorage) *AuthHandler {
	return &AuthHandler{myAuth.AuthService{StmtsStorage: storage}}
}

type AuthHandler struct {
	authServ myAuth.AuthService
}

func (h *AuthHandler) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))

	r.Route("/", func(r chi.Router) {
		r.Post("/signup", makeHandler(h.RegistrationHandler))
		r.Post("/signin", makeHandler(h.AuthorizationHandler))
		r.Put("/users/0", makeHandler(h.UpdateHandler))
	})
	return r
}

func readReqData(r *http.Request, userData interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrapf(err, "$read error$")
	}
	err = json.Unmarshal(body, userData)
	if err != nil {
		return errors.Wrapf(err, "$unmarshal error$")
	}
	return nil
}

func checkEmail(email string) error {
	if validEmail.FindStringSubmatch(email) == nil {
		return errors.Wrap(myerr.BadRequest, "$email doesnt match pattern$")
	}
	return nil
}

func readReqAndCheckEmail(r *http.Request, userData request.EmailGetter) error {

	if err := readReqData(r, userData); err != nil {
		return errors.Wrap(err, "read req data")

	}
	if err := checkEmail(userData.GetEmail()); err != nil {
		return errors.Wrap(err, "error in check email")
	}
	return nil
}

func (h *AuthHandler) RegistrationHandler(w http.ResponseWriter, r *http.Request) error {

	var user request.RegUser
	if err := readReqAndCheckEmail(r, &user); err != nil {
		resp, err := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
		if err != nil {
			return errors.Wrap(err, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusBadRequest)
		return nil
	}
	err := h.authServ.RegisterUser(&user)
	switch errors.Cause(err) {
	case myerr.Conflict:
		resp, err := myerr.ErrMarshal("email already exists")
		if err != nil {
			return errors.Wrap(err, "marshal error")
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
		resp, err := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
		if err != nil {
			return errors.Wrap(err, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusBadRequest)
		return nil
	}
	token, err := h.authServ.AuthorizeUser(&user)
	switch errors.Cause(err) {
	case myerr.Unauthorized:
		resp, err := myerr.ErrMarshal("invalid email or password")
		if err != nil {
			return errors.Wrap(err, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusConflict)
	case myerr.BadRequest:
		resp, err := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
		if err != nil {
			return errors.Wrap(err, "marshal error")
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
		resp, err := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
		if err != nil {
			return errors.Wrap(err, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusBadRequest)
		return nil
	}
	fmt.Println(user.Birthday)
	if user.TokenType != "bearer" {
		resp, err := myerr.ErrMarshal("invalid token type")
		if err != nil {
			return errors.Wrap(err, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusBadRequest)
		return nil
	}
	dbUser, err := h.authServ.UpdateUser(&user)
	switch errors.Cause(err) {
	case myerr.Unauthorized:
		resp, err := myerr.ErrMarshal("unauthorized")
		if err != nil {
			return errors.Wrap(err, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusUnauthorized)
	case myerr.Success:
		resp, err := json.Marshal(dbUser)
		if err != nil {
			return errors.Wrap(err, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusOK)
	default:
		return errors.Wrap(err, "user cant be update")
	}
	return nil
}
