package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
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

func makeHandler(fn func(http.ResponseWriter, *http.Request) *myerr.AppError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err.Message != "" {
			http.Error(w, err.Message, err.Code)
		}
		if err.Err != nil {
			log.Println(err.Err)
		}
	}
}

type AuthHandler struct {
}

func (h *AuthHandler) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/signup", makeHandler(h.RegistrationHandler))
		r.Post("/signin", makeHandler(h.AuthorizationHandler))
		r.Put("/users/0", makeHandler(h.UpdateHandler))
	})
	return r
}

func readReqData(r *http.Request, userData interface{}) *myerr.AppError {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "read error")
		return myerr.NewErr(err, "Unprocessable Entity", 422)
	}

	err = json.Unmarshal(body, userData)
	if err != nil {
		err = errors.Wrap(err, "unmarshal error")
		return myerr.NewErr(err, "Unprocessable Entity", 422)
	}
	return nil
}

func checkEmail(email string) *myerr.AppError {
	if validEmail.FindStringSubmatch(email) == nil {
		return myerr.NewErr(nil, "email doesnt match pattern", 400)
	}
	return nil
}

func readReqAndCheckEmail(r *http.Request, userData request.EmailGetter) *myerr.AppError {

	if err := readReqData(r, userData); err != nil {
		return err.MyWrap("read req data")

	}
	if err := checkEmail(userData.GetEmail()); err != nil {
		return err.MyWrap("error in check email")
	}
	return nil
}

func (h *AuthHandler) RegistrationHandler(w http.ResponseWriter, r *http.Request) *myerr.AppError {
	var user request.RegUser

	if err := readReqAndCheckEmail(r, &user); err != nil {
		return err.MyWrap("error in readreqandcheckemail")
	}
	fmt.Printf("%+v", user)
	err := myAuth.RegisterUser(&user)
	return err.MyWrap("error in registeruser method")
}

func (h *AuthHandler) AuthorizationHandler(w http.ResponseWriter, r *http.Request) *myerr.AppError {
	var user request.AuthUser
	if err := readReqAndCheckEmail(r, &user); err != nil {
		return err.MyWrap("error in readreqandcheckemail")
	}
	err := myAuth.AuthorizeUser(&user)
	return err.MyWrap("error in authorizeuser method")

}

func (h *AuthHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) *myerr.AppError {
	var user request.UpdateUser
	if err := readReqData(r, &user); err != nil {
		return err.MyWrap("error in readreqdata")
	}
	err := myAuth.UpdateUser(&user)
	return err.MyWrap("error in updateuser method")
}
