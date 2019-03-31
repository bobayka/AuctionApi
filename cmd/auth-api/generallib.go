package main

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/services"
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
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func CheckTokenMiddleware(store *postgres.UsersStorage) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) error {
			if r.Header.Get("token_type") != "bearer" {
				resp, err := myerr.ErrMarshal("invalid token type")
				if err != nil {
					return errors.Wrap(err, "marshal error")
				}
				myerr.Error(w, string(resp), http.StatusBadRequest)
				return nil
			}
			s, err := services.CheckValidToken(r.Header.Get("access_token"), store)
			if errors.Cause(err) != myerr.Success {
				switch errors.Cause(err) {
				case myerr.Unauthorized:
					resp, errM := myerr.ErrMarshal("unauthorized")
					if errM != nil {
						return errors.Wrap(errM, "marshal error")
					}
					myerr.Error(w, string(resp), http.StatusUnauthorized)
					return nil
				default:
					return errors.Wrap(err, "lot cant be create")
				}
			}
			ctx := context.WithValue(r.Context(), "user_id", s.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return nil
		}
		return http.HandlerFunc(makeHandler(fn))
	}
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

//func readDataCheckTokenType(r *http.Request, u request.TokenTypeGetter) ([]byte, error) {
//	if err := readReqData(r, &u); err != nil {
//		resp, errM := myerr.ErrMarshal(myerr.GetClientErr(err.Error()))
//		if errM != nil {
//			return nil, errors.Wrap(err, "marshal error")
//		}
//		return resp, nil
//	}
//	//if u.GetTokenType() != "bearer" {
//	//	resp, err := myerr.ErrMarshal("invalid token type")
//	//	if err != nil {
//	//		return nil, errors.Wrap(err, "marshal error")
//	//	}
//	//	return resp, nil
//	//}
//	return nil, nil
//}
