package authhandlers

import (
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/cmd/utilities"
	"gitlab.com/bobayka/courseproject/internal/postgres/storage"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/responce"
	"gitlab.com/bobayka/courseproject/internal/services"
	"net/http"
)

func NewAuthHandler(storage storage.Storage) *AuthHandler {
	return &AuthHandler{services.AuthService{StmtsStorage: storage}}
}

type AuthHandler struct {
	authServ services.AuthService
}

func (h *AuthHandler) Routes() *chi.Mux {

	r := chi.NewRouter()
	r.Post("/signup", utility.MakeHandler(h.RegistrationHandler))
	r.Post("/signin", utility.MakeHandler(h.AuthorizationHandler))
	return r
}

func (h *AuthHandler) RegistrationHandler(w http.ResponseWriter, r *http.Request) error {
	var user request.RegUser
	//if err := utility.ReadReqAndCheckEmail(r, &user); err != nil {
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
