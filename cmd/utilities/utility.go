package utility

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	request "gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/responce"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type key int

const UserIDKey key = 0

func MakeHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := myerr.ErrHandler(w, fn(w, r))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func GetIDURLParam(r *http.Request) (int64, error) {
	id := chi.URLParam(r, "id")
	key, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, errors.Wrap(myerr.ErrBadRequest, "$Wrong ID$")
	}
	return key, nil
}

func ReadReqData(r *http.Request, data interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrapf(myerr.ErrBadRequest, "$read error$: %s", err)
	}
	err = json.Unmarshal(body, data)
	if err != nil {
		return errors.Wrapf(myerr.ErrBadRequest, "$unmarshal error$: %s", err)
	}
	return nil
}

func GetTokenUserID(r *http.Request) (int64, error) {
	ctx := r.Context()
	OwnID, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return 0, errors.New("r.context doesn't contain user_id")
	}
	return OwnID, nil
}

func ReadReqAndCheckEmail(r *http.Request, userData request.EmailGetter) error {

	if err := ReadReqData(r, userData); err != nil {
		return errors.Wrap(err, "read req data")
	}
	if err := CheckEmail(userData.GetEmail()); err != nil {
		return errors.Wrap(err, "error in check email")
	}
	return nil
}

func GetUserIDURL(r *http.Request) (int64, error) {
	userID, err := GetIDURLParam(r)
	if err != nil {
		return 0, errors.Wrap(myerr.ErrBadRequest, "$Wrong User ID$")
	}
	if userID == 0 {
		userID, err = GetTokenUserID(r)
		if err != nil {
			return 0, err
		}
	}
	return userID, nil
}

func MarshalAndRespondJSON(w http.ResponseWriter, data interface{}) error {
	resp, errM := json.Marshal(data)
	if errM != nil {
		return errors.Wrap(errM, "marshal error")
	}
	responce.RespondJSON(w, string(resp), http.StatusOK)
	return nil
}
