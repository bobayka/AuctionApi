package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	myAuth "gitlab.com/bobayka/courseproject/internal/auth"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/pkg/myerr"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var validEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func makeHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		switch errors.Cause(err) {
		case myerr.BadRequest:
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
		case myerr.UnprocessableEntity:
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		case myerr.Success:
			http.Error(w, "", http.StatusOK)
		case myerr.Accepted:
			http.Error(w, myAuth.Bearer.FindString(err.Error()), http.StatusAccepted)
		case myerr.Unauthorized:
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
		case myerr.Created:
			http.Error(w, "", http.StatusCreated)
		default:
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)

		}
	}
}

type AuthHandler struct {
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
		return errors.Wrapf(myerr.UnprocessableEntity, "read error: %s", err)
	}

	err = json.Unmarshal(body, userData)
	if err != nil {
		return errors.Wrapf(myerr.UnprocessableEntity, "unmarshal error: %s", err)
	}
	return nil
}

func checkEmail(email string) error {
	if validEmail.FindStringSubmatch(email) == nil {
		return errors.Wrap(myerr.UnprocessableEntity, "email doesnt match pattern")
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
		return errors.Wrap(err, "error in readreqandcheckemail")
	}
	fmt.Printf("%+v", user)
	err := myAuth.RegisterUser(&user)
	return errors.Wrap(err, "error in registeruser method")
}

func (h *AuthHandler) AuthorizationHandler(w http.ResponseWriter, r *http.Request) error {
	var user request.AuthUser
	if err := readReqData(r, &user); err != nil {
		return errors.Wrap(err, "error in readreqdata")
	}
	err := myAuth.AuthorizeUser(&user)
	return errors.Wrap(err, "error in authorizeuser method")

}

func (h *AuthHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	var user request.UpdateUser
	if err := readReqData(r, &user); err != nil {
		return errors.Wrap(err, "error in readreqdata")
	}
	err := myAuth.UpdateUser(&user)
	return errors.Wrap(err, "error in updateuser method")
}
