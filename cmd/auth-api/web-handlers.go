package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/services"
	"gitlab.com/bobayka/courseproject/pkg/myerr"
	"html/template"
	"net/http"
)

var templates map[string]*template.Template

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templates["index"] = template.Must(template.ParseFiles("html/index.html", "html/base.html"))
}

func renderTemplate(w http.ResponseWriter, name string, template string, viewModel interface{}) {
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, "can't find template", http.StatusInternalServerError)
	}
	err := tmpl.ExecuteTemplate(w, template, viewModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NewWebHandler(storage *postgres.UsersStorage) *WebHandler {
	return &WebHandler{services.Auth{StmtsStorage: storage}}
}

type WebHandler struct {
	webServ services.Auth
}

func (wb *WebHandler) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(CheckTokenMiddleware(wb.webServ.StmtsStorage))
	r.Get("/users/{id:[0-9]*}/lots", makeHandler(wb.GetUserLots))
	return r
}
func (wb *WebHandler) GetUserLots(w http.ResponseWriter, r *http.Request) error {
	UserID, err := GetUserIDQueryParam(r)
	if err != nil {
		if err == myerr.BadRequest {
			jsonRespond(w, "Wrong User ID", http.StatusBadRequest)
			return nil
		}
		return errors.Wrap(err, "cant get id url param") //после отладки можно убрать
	}
	dbLots, err := wb.webServ.GetUserLotsByID(UserID, "own")
	switch errors.Cause(err) {
	case myerr.NotFound:
		jsonRespond(w, "Content by the passed ID could not be found", http.StatusNotFound)
	case myerr.Success:
		renderTemplate(w, "index", "base", dbLots)
	default:
		return errors.Wrap(err, "lot cant be get")
	}
	return nil
}
