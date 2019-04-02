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

var lotsTypes map[string]bool

func init() {
	if lotsTypes == nil {
		lotsTypes = make(map[string]bool)
	}
	lotsTypes["own"] = true
	lotsTypes["buyed"] = true
	lotsTypes[""] = true
}

func NewAuthHandler(storage *postgres.UsersStorage) *AuthHandler {
	return &AuthHandler{services.Auth{StmtsStorage: storage}}
}

type AuthHandler struct {
	authServ services.Auth
}

func (h *AuthHandler) Routes() *chi.Mux {
	//r.Use(middleware.RealIP)
	//r.Use(middleware.RequestID)
	r := chi.NewRouter()

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
			r.Get("/{id:[0-9]*}/lots", makeHandler(h.GetUserLots))
		})

	})
	return r
}

func (h *AuthHandler) RegistrationHandler(w http.ResponseWriter, r *http.Request) error {

	var user request.RegUser
	if err := readReqAndCheckEmail(r, &user); err != nil {
		jsonRespond(w, myerr.GetClientErr(err.Error()), http.StatusBadRequest)
		return nil
	}
	if len(user.Password) < 6 {
		jsonRespond(w, "Password less than 6 characters", http.StatusBadRequest)
		return nil
	}
	err := h.authServ.RegisterUser(&user)
	switch errors.Cause(err) {
	case myerr.Conflict:
		jsonRespond(w, "email already exists", http.StatusConflict)
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
		jsonRespond(w, myerr.GetClientErr(err.Error()), http.StatusBadRequest)
		return nil
	}
	token, err := h.authServ.AuthorizeUser(&user)
	switch errors.Cause(err) {
	case myerr.Unauthorized:
		jsonRespond(w, "invalid email or password", http.StatusConflict)
	case myerr.BadRequest:
		jsonRespond(w, myerr.GetClientErr(err.Error()), http.StatusBadRequest)
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
		jsonRespond(w, myerr.GetClientErr(err.Error()), http.StatusBadRequest)
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
		jsonRespond(w, "unauthorized", http.StatusUnauthorized)
	case myerr.Success:
		resp, _ := json.Marshal(dbUser)
		myerr.Error(w, string(resp), http.StatusOK)
	default:
		return errors.Wrap(err, "user cant be update")
	}
	return nil
}

func (h *AuthHandler) GetHandler(w http.ResponseWriter, r *http.Request) error {
	UserID, err := getIDURLParam(r)
	if err != nil {
		jsonRespond(w, "Wrong User ID", http.StatusBadRequest)
		return nil
	}
	if UserID == 0 {
		ctx := r.Context()
		OwnID, ok := ctx.Value(UserIDKey).(int64)
		if !ok {
			return errors.New("r.context doesn't contain user_id")
		}
		UserID = OwnID
	}
	dbUser, err := h.authServ.GetUserByID(UserID)
	switch errors.Cause(err) {
	case myerr.NotFound:
		jsonRespond(w, "Content by the passed ID could not be found", http.StatusNotFound)
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

func GetUserIDQueryParam(r *http.Request) (int64, error) {
	UserID, err := getIDURLParam(r)
	if err != nil {
		return 0, myerr.BadRequest
	}
	if UserID == 0 {
		ctx := r.Context()
		OwnID, ok := ctx.Value(UserIDKey).(int64)
		if !ok {
			return 0, errors.New("r.context doesn't contain user_id")
		}
		UserID = OwnID
	}
	return UserID, nil
}

func (h *AuthHandler) GetUserLots(w http.ResponseWriter, r *http.Request) error {
	UserID, err := GetUserIDQueryParam(r)
	if err != nil {
		if err == myerr.BadRequest {
			jsonRespond(w, "Wrong User ID", http.StatusBadRequest)
			return nil
		}
		return errors.Wrap(err, "cant get id url param") //после отладки можно убрать
	}
	lotType := r.URL.Query().Get("type")
	if !lotsTypes[lotType] {
		jsonRespond(w, "Wrong lot type", http.StatusBadRequest)
		return nil
	}
	dbLots, err := h.authServ.GetUserLotsByID(UserID, lotType)
	switch errors.Cause(err) {
	case myerr.NotFound:
		jsonRespond(w, "Content by the passed ID could not be found", http.StatusNotFound)
	case myerr.Success:
		resp, errM := json.Marshal(dbLots)
		if errM != nil {
			return errors.Wrap(errM, "marshal error")
		}
		myerr.Error(w, string(resp), http.StatusOK)
	default:
		return errors.Wrap(err, "lot cant be get")
	}
	return nil
}
