package main

import (
	"encoding/json"
	"errors"
	"fmt"
	errors2 "github.com/pkg/errors"
	myAuth "gitlab.com/bobayka/courseproject/internal/auth"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var validEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type AuthHandler struct {
}

func readReqData(w http.ResponseWriter, r *http.Request, userData request.User) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return errors2.Wrap(err, "read error")

	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return errors2.Wrap(err, "unmarshal error")
	}
	return nil
}

func checkEmail(w http.ResponseWriter, email string) error {
	if validEmail.FindStringSubmatch(email) == nil {
		_, err := w.Write([]byte("Email doesnt match pattern"))
		if err != nil {
			log.Println("error in write method checkEmail")
		}
		return errors.New("email doesnt match pattern")
	}
	return nil
}

func readReqAndCheckEmail(w http.ResponseWriter, r *http.Request, userData request.User) error {

	if err := readReqData(w, r, userData); err != nil {
		return errors2.Wrap(err, "read req data")
	}
	if err := checkEmail(w, userData.GetEmail()); err != nil {
		return errors2.Wrap(err, "check email")
	}
	return nil
}

func (auth AuthHandler) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var user request.RegUser
	err := readReqAndCheckEmail(w, r, &user)
	if err != nil {
		log.Printf("registration handler: %s", err)
		return
	}
	status := myAuth.RegisterUser(&user)
	_, err = w.Write([]byte(errors2.Cause(status).Error()))
	if err != nil {
		log.Println("error in write method registrationhandler")
	}
}

func (auth AuthHandler) AuthorizationHandler(w http.ResponseWriter, r *http.Request) {
	var user request.AuthUser
	err := readReqAndCheckEmail(w, r, &user)
	if err != nil {
		log.Printf("authorization handler: %s", err)
		return
	}
	status := myAuth.AuthorizeUser(&user)
	_, err = w.Write([]byte(errors2.Cause(status).Error()))
	if err != nil {
		log.Println("error in write method authorizationhandler")
	}

}

func (auth AuthHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var user request.UpdateUser
	err := readReqData(w, r, &user)
	if err != nil {
		log.Printf("update handler: %s", err)
		return
	}
	status := myAuth.UpdateUser(&user)
	_, err = w.Write([]byte(errors2.Cause(status).Error()))
	if err != nil {
		log.Println("error in write method UpdateHandler")
	}
}
