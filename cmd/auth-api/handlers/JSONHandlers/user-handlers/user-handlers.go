package userhandlers

import (
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/cmd/utilities"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/services"
	"net/http"
)

// nolint: gochecknoglobals
var lotsTypes = map[string]bool{
	"own":   true,
	"buyed": true,
	"":      true,
}

func NewUserHandlers(storage *postgres.UsersStorage) *UserHandler {
	return &UserHandler{services.UserService{StmtsStorage: storage}}
}

type UserHandler struct {
	authServ services.UserService
}

func (u *UserHandler) Routes() *chi.Mux {

	r := chi.NewRouter()

	r.Put("/0", utility.MakeHandler(u.UpdateHandler))
	r.Get("/{id:[0-9]*}", utility.MakeHandler(u.GetHandler))
	r.Get("/{id:[0-9]*}/lots", utility.MakeHandler(u.GetUserLots))

	return r
}

func (u *UserHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	var user request.UpdateUser

	if err := utility.ReadReqData(r, &user); err != nil {
		return errors.Wrap(err, "cant be read req")
	}
	userID, err := utility.GetTokenUserID(r)
	if err != nil {
		return errors.Wrap(err, "cant get token user id")
	}
	dbUser, err := u.authServ.UpdateUser(&user, userID)
	if err != nil {
		return errors.Wrap(err, "user cant be update")
	}

	err = utility.MarshalAndRespondJSON(w, dbUser)
	if err != nil {
		return errors.Wrap(err, "marshal and respondJSON")
	}
	return nil
}

func (u *UserHandler) GetHandler(w http.ResponseWriter, r *http.Request) error {
	userID, err := utility.GetUserIDURL(r)
	if err != nil {
		return errors.Wrap(err, "cant get token user id")
	}

	dbUser, err := u.authServ.GetUserByID(userID)
	if err != nil {
		return errors.Wrap(err, "lot cant be get")
	}

	err = utility.MarshalAndRespondJSON(w, dbUser)
	if err != nil {
		return errors.Wrap(err, "marshal and respondJSON")
	}
	return nil
}

func (u *UserHandler) GetUserLots(w http.ResponseWriter, r *http.Request) error {
	UserID, err := utility.GetUserIDURL(r)
	if err != nil {
		return errors.Wrap(err, "cant get id url param") //после отладки можно убрать
	}
	lotType := r.URL.Query().Get("type")
	if !lotsTypes[lotType] {
		return errors.Wrap(myerr.ErrBadRequest, "$Wrong lot type$")
	}
	dbLots, err := u.authServ.GetUserLotsByID(UserID, lotType)
	if err != nil {
		return errors.Wrap(err, "lot cant be get")
	}
	err = utility.MarshalAndRespondJSON(w, dbLots)
	if err != nil {
		return errors.Wrap(err, "marshal and respondJSON")
	}
	return nil
}
