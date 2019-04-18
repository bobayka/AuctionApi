package userWeb

import (
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/auth-api/handlers/HTMLHandlers"
	utility "gitlab.com/bobayka/courseproject/cmd/utilities"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/services"
	"html/template"
	"net/http"
)

const baseHTMLDirectory = "cmd/auth-api/handlers/HTMLHandlers/user-handlers/html/"

type WebUserHandler struct {
	webUser services.UserService
	templs  HTMLHandlers.Templates
}

func NewWebUserHandler(storage *postgres.UsersStorage) *WebUserHandler {
	templ := HTMLHandlers.Templates{
		"index": template.Must(template.ParseFiles(baseHTMLDirectory+"index.html", baseHTMLDirectory+"base.html")),
	}
	return &WebUserHandler{webUser: services.UserService{StmtsStorage: storage}, templs: templ}
}
func (wb *WebUserHandler) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/{id:[0-9]*}/lots", utility.MakeHandler(wb.GetUserLotsHandler))

	return r
}

func (wb *WebUserHandler) GetUserLotsHandler(w http.ResponseWriter, r *http.Request) error {
	UserID, err := utility.GetUserIDURL(r)
	if err != nil {
		return errors.Wrap(err, "cant get id url param") //после отладки можно убрать
	}
	dbLots, err := wb.webUser.GetUserLotsByID(UserID, "own")
	if err != nil {
		return errors.Wrap(err, "lot cant be get")
	}
	wb.templs.RenderTemplate(w, "index", "utilities", dbLots)

	return nil
}
