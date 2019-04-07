package webapi

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/utilities"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/services"
	"html/template"
	"net/http"
)

type templates map[string]*template.Template

func (t templates) renderTemplate(w http.ResponseWriter, name string, template string, viewModel interface{}) {
	tmpl, ok := t[name]
	if !ok {
		http.Error(w, "can't find template", http.StatusInternalServerError)
	}
	err := tmpl.ExecuteTemplate(w, template, viewModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NewWebHandler(storage *postgres.UsersStorage) *WebHandler {
	templ := templates{
		"index": template.Must(template.ParseFiles("html/index.html", "html/base.html"))}
	return &WebHandler{webServ: services.Auth{StmtsStorage: storage}, templs: templ}
}

type WebHandler struct {
	webServ services.Auth
	templs  templates
}

func (wb *WebHandler) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(utility.CheckTokenMiddleware(wb.webServ.StmtsStorage))
	r.Get("/users/{id:[0-9]*}/lots", utility.MakeHandler(wb.GetUserLots))
	return r
}
func (wb *WebHandler) GetUserLots(w http.ResponseWriter, r *http.Request) error {
	UserID, err := utility.GetUserIDURL(r)
	if err != nil {
		return errors.Wrap(err, "cant get id url param") //после отладки можно убрать
	}
	dbLots, err := wb.webServ.GetUserLotsByID(UserID, "own")
	if err != nil {
		return errors.Wrap(err, "lot cant be get")
	}
	wb.templs.renderTemplate(w, "index", "utilities", dbLots)

	return nil
}
