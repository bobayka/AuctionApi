package authweb

import (
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/gateway-api/handlers/htmlhandlers"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/cmd/utilities"
	"gitlab.com/bobayka/courseproject/internal/postgres/storage"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/services"
	"gitlab.com/bobayka/courseproject/pkg/customTime"
	"html/template"
	"net/http"
	"time"
)

const sessionCookieName = "x-authorization"
const baseHTMLDirectory = "cmd/gateway-api/handlers/htmlhandlers/auth-handlers/html/"

func NewWebAuthHandler(storage storage.Storage) *WebAuthHandler {
	templ := htmlhandlers.Templates{
		"auth":    template.Must(template.ParseFiles(baseHTMLDirectory+"auth.html", baseHTMLDirectory+"base.html")),
		"registr": template.Must(template.ParseFiles(baseHTMLDirectory+"registr.html", baseHTMLDirectory+"base.html")),
	}
	return &WebAuthHandler{webAuth: services.AuthService{StmtsStorage: storage}, templs: templ}
}

type WebAuthHandler struct {
	webAuth services.AuthService
	templs  htmlhandlers.Templates
}

func (wa *WebAuthHandler) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", wa.InitialHandler)
	r.Get("/signin", wa.GetAuthPageHandler)
	r.Post("/signin", utility.MakeHandler(wa.AuthorizationHandler))
	r.Get("/signup", wa.GetRegisterPageHandler)
	r.Post("/signup", utility.MakeHandler(wa.RegistrationHandler))

	return r
}
func (wa *WebAuthHandler) InitialHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:5000/signin", http.StatusFound)
}
func (wa *WebAuthHandler) GetAuthPageHandler(w http.ResponseWriter, r *http.Request) {
	wa.templs.RenderTemplate(w, "auth", "utilities", nil)
}

func (wa *WebAuthHandler) GetRegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	wa.templs.RenderTemplate(w, "registr", "utilities", nil)
}

func (wa *WebAuthHandler) RegistrationHandler(w http.ResponseWriter, r *http.Request) error {
	firstName := r.PostFormValue("firstName")
	lastName := r.PostFormValue("lastName")
	birthday := r.PostFormValue("birthday")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	if err := utility.CheckEmail(email); err != nil {
		return errors.Wrap(err, "error in check email")
	}
	if len(password) < 6 {
		return errors.Wrap(myerr.ErrBadRequest, "$Password less than 6 characters$")
	}
	var birthdayPtr *customtime.CustomTime
	if birthday == "" {
		birthdayPtr = nil
	} else {
		t, err := time.Parse(customtime.CTLayout, birthday)
		if err != nil {
			return errors.Wrap(myerr.ErrBadRequest, "$birthday date can be '2006-01-02' type$")
		}
		birthdayDate := customtime.CustomTime(t)
		birthdayPtr = &birthdayDate
	}
	webUser := request.RegUser{
		BasicInfo: request.BasicInfo{
			FirstName: firstName,
			LastName:  lastName,
			Birthday:  birthdayPtr,
		},
		GeneralInfo: request.GeneralInfo{
			Email:    email,
			Password: password,
		},
	}
	err := wa.webAuth.RegisterUser(&webUser)
	if err != nil {
		return errors.Wrap(err, "user cant be registered")
	}
	http.Redirect(w, r, "http://localhost:5000/signin", http.StatusFound)
	return nil
}

func (wa *WebAuthHandler) AuthorizationHandler(w http.ResponseWriter, r *http.Request) error {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	if err := utility.CheckEmail(email); err != nil {
		return errors.Wrap(err, "error in check email")
	}
	if len(password) < 6 {
		return errors.Wrap(myerr.ErrBadRequest, "$Password less than 6 characters$")
	}
	user := &request.AuthUser{GeneralInfo: request.GeneralInfo{Email: email, Password: password}}
	token, err := wa.webAuth.AuthorizeUser(user)
	if err != nil {
		if errors.Cause(err) == myerr.ErrUnauthorized {
			http.Redirect(w, r, "http://localhost:5000/signup", http.StatusFound)
			return nil
		}
		return errors.Wrap(err, "cant authorize user")
	}
	http.SetCookie(w, &http.Cookie{Name: sessionCookieName, Value: token})
	//http.Redirect(w, r, )
	return nil
}
