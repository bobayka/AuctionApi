package authhandlers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/cmd/utilities"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/responce"
	"gitlab.com/bobayka/courseproject/internal/services"
	"net/http"
)

// nolint: gochecknoglobals
var lotsTypes = map[string]bool{
	"own":   true,
	"buyed": true,
	"":      true,
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
	r.Post("/signup", utility.MakeHandler(h.RegistrationHandler))
	r.Post("/signin", utility.MakeHandler(h.AuthorizationHandler))
	r.Route("/users", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(utility.CheckTokenMiddleware(h.authServ.StmtsStorage))
			r.Put("/0", utility.MakeHandler(h.UpdateHandler))
			r.Get("/{id:[0-9]*}", utility.MakeHandler(h.GetHandler))
			r.Get("/{id:[0-9]*}/lots", utility.MakeHandler(h.GetUserLots))
		})

	})
	return r
}

func (h *AuthHandler) RegistrationHandler(w http.ResponseWriter, r *http.Request) error {

	var user request.RegUser
	if err := utility.ReadReqAndCheckEmail(r, &user); err != nil {
		return errors.Wrap(err, "cant be read req or checked email")
	}
	if len(user.Password) < 6 {
		return errors.Wrap(myerr.ErrBadRequest, "$Password less than 6 characters$")
	}
	err := h.authServ.RegisterUser(&user)
	if err != nil {
		return errors.Wrap(err, "user cant be registered")
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h *AuthHandler) AuthorizationHandler(w http.ResponseWriter, r *http.Request) error {
	var user request.AuthUser
	if err := utility.ReadReqAndCheckEmail(r, &user); err != nil {
		return errors.Wrap(err, "cant be read req or checked email")
	}
	token, err := h.authServ.AuthorizeUser(&user)
	if err != nil {
		return errors.Wrap(err, "cant authorize user")
	}
	resp := `{"token_type": "bearer","access_token": "` + token + `"}`
	responce.RespondJSON(w, resp, http.StatusOK)
	return nil
}

func (h *AuthHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	var user request.UpdateUser

	if err := utility.ReadReqData(r, &user); err != nil {
		return errors.Wrap(err, "cant be read req")
	}
	userID, err := utility.GetTokenUserID(r)
	if err != nil {
		return errors.Wrap(err, "cant get token user id")
	}
	dbUser, err := h.authServ.UpdateUser(&user, userID)
	if err != nil {
		return errors.Wrap(err, "user cant be update")
	}

	err = utility.MarshalAndRespondJSON(w, dbUser)
	if err != nil {
		return errors.Wrap(err, "marshal and respondJSON")
	}
	return nil
}

func (h *AuthHandler) GetHandler(w http.ResponseWriter, r *http.Request) error {
	userID, err := utility.GetUserIDURL(r)
	if err != nil {
		return errors.Wrap(err, "cant get token user id")
	}

	dbUser, err := h.authServ.GetUserByID(userID)
	if err != nil {
		return errors.Wrap(err, "lot cant be get")
	}

	err = utility.MarshalAndRespondJSON(w, dbUser)
	if err != nil {
		return errors.Wrap(err, "marshal and respondJSON")
	}
	return nil
}

func (h *AuthHandler) GetUserLots(w http.ResponseWriter, r *http.Request) error {
	UserID, err := utility.GetUserIDURL(r)
	if err != nil {
		return errors.Wrap(err, "cant get id url param") //после отладки можно убрать
	}
	lotType := r.URL.Query().Get("type")
	if !lotsTypes[lotType] {
		return errors.Wrap(myerr.ErrBadRequest, "$Wrong lot type$")
	}
	dbLots, err := h.authServ.GetUserLotsByID(UserID, lotType)
	if err != nil {
		return errors.Wrap(err, "lot cant be get")
	}
	err = utility.MarshalAndRespondJSON(w, dbLots)
	if err != nil {
		return errors.Wrap(err, "marshal and respondJSON")
	}
	return nil
}
