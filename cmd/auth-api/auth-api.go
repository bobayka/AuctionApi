package authapi

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gitlab.com/bobayka/courseproject/cmd/auth-api/handlers/HTMLHandlers/auth-handlers"
	"gitlab.com/bobayka/courseproject/cmd/auth-api/handlers/HTMLHandlers/lot-handlers"
	"gitlab.com/bobayka/courseproject/cmd/auth-api/handlers/HTMLHandlers/user-handlers"
	"gitlab.com/bobayka/courseproject/cmd/auth-api/handlers/JSONHandlers/auth-handlers"
	"gitlab.com/bobayka/courseproject/cmd/auth-api/handlers/JSONHandlers/lot-handlers"
	"gitlab.com/bobayka/courseproject/cmd/auth-api/handlers/JSONHandlers/user-handlers"
	"gitlab.com/bobayka/courseproject/cmd/auth-api/handlers/JSONHandlers/websocket-handlers"
	"gitlab.com/bobayka/courseproject/cmd/utilities"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/postgres/PGBackground"
)

type authApi struct {
	storage *postgres.UsersStorage
	auth    *authhandlers.AuthHandler
	lot     *lothandlers.LotServiceHandler
	user    *userhandlers.UserHandler
	webAuth *authWeb.WebAuthHandler
	webUser *userWeb.WebUserHandler
	webLot  *lotWeb.WebLotHandler
	wsLot   *websocket_handlers.LotWSHandler
}

func NewAuthApi(storage *postgres.UsersStorage) *authApi {
	auth := authhandlers.NewAuthHandler(storage)
	lot := lothandlers.NewLotServiceHandler(storage)
	user := userhandlers.NewUserHandlers(storage)
	webUser := userWeb.NewWebUserHandler(storage)
	webLot := lotWeb.NewWebLotHandler(storage)
	webAuth := authWeb.NewWebAuthHandler(storage)
	wsLot := websocket_handlers.NewLotWSHandler(storage)
	PGBackground.StartDBBackgroundProcesses(storage, wsLot.Hub)
	return &authApi{storage: storage, auth: auth, lot: lot, user: user, webUser: webUser, webLot: webLot, webAuth: webAuth, wsLot: wsLot}
}

func (a *authApi) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	authRouter := a.auth.Routes()
	userRouter := a.user.Routes()
	lotRouter := a.lot.Routes()
	webUserRouter := a.webUser.Routes()
	webAuthRouter := a.webAuth.Routes()
	webLotRouter := a.webLot.Routes()
	wsLotsRouter := a.wsLot.Routes()

	router.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Mount("/v1/auction", authRouter)
		r.Group(func(r chi.Router) {
			r.Use(utility.CheckTokenMiddleware(a.storage))
			r.Mount("/v1/auction/users", userRouter)
			r.Mount("/v1/auction/lots", lotRouter)
		})
	})
	router.Mount("/", webAuthRouter)
	router.Group(func(r chi.Router) {
		r.Use(utility.CheckCookieMiddleware(a.storage))
		r.Mount("/w/auction/user", webUserRouter)
		r.Mount("/w/auction/lots", webLotRouter)
		r.Mount("/auction", wsLotsRouter)
	})
	return router
}
