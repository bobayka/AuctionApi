package gatewayApi

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gitlab.com/bobayka/courseproject/cmd/Protobuf"
	"gitlab.com/bobayka/courseproject/cmd/gateway-api/handlers/HTMLHandlers/auth-handlers"
	"gitlab.com/bobayka/courseproject/cmd/gateway-api/handlers/HTMLHandlers/lot-handlers"
	"gitlab.com/bobayka/courseproject/cmd/gateway-api/handlers/JSONHandlers/auth-handlers"
	"gitlab.com/bobayka/courseproject/cmd/gateway-api/handlers/JSONHandlers/lot-handlers"
	"gitlab.com/bobayka/courseproject/cmd/gateway-api/handlers/JSONHandlers/user-handlers"
	"gitlab.com/bobayka/courseproject/cmd/gateway-api/handlers/JSONHandlers/websocket-handlers"
	"gitlab.com/bobayka/courseproject/cmd/utilities"
	"gitlab.com/bobayka/courseproject/internal/postgres/storage"
	"time"
)

type gatewayApi struct {
	lotServiceClient lotspb.LotsServiceClient
	storage          *storage.SessionStorage
	auth             *authhandlers.AuthHandler
	lot              *lothandlers.LotServiceHandler
	user             *userhandlers.UserHandler
	webAuth          *authWeb.WebAuthHandler
	webLot           *lotWeb.WebLotHandler
	wsLot            *websocket_handlers.LotWSHandler
}

func NewGatewayApi(storage storage.Storage, client lotspb.LotsServiceClient) *gatewayApi {
	auth := authhandlers.NewAuthHandler(storage)
	lot := lothandlers.NewLotHandler(client)
	user := userhandlers.NewUserHandlers(storage.Users)
	webLot := lotWeb.NewWebLotHandler(client)
	webAuth := authWeb.NewWebAuthHandler(storage)
	wsLot := websocket_handlers.NewLotWSHandler(client)
	wsLot.BackGroundFinishEndedLots(time.Second)
	return &gatewayApi{storage: storage.Sessions, auth: auth, lot: lot, user: user, webLot: webLot, webAuth: webAuth, wsLot: wsLot}

}

func (a *gatewayApi) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	authRouter := a.auth.Routes()
	userRouter := a.user.Routes()
	lotRouter := a.lot.Routes()
	webAuthRouter := a.webAuth.Routes()
	webLotRouter := a.webLot.Routes()
	wsLotsRouter := a.wsLot.Routes()

	router.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Mount("/v1/auction", authRouter)
		r.Group(func(r chi.Router) {
			r.Use(utility.CheckTokenMiddleware(a.storage))
			r.Get("/v1/auction/users/{id:[0-9]*}/lots", utility.MakeHandler(a.lot.GetUserLotsHandler))
			r.Mount("/v1/auction/users", userRouter)
			r.Mount("/v1/auction/lots", lotRouter)
		})
	})
	router.Mount("/", webAuthRouter)
	router.Group(func(r chi.Router) {
		r.Use(utility.CheckCookieMiddleware(a.storage))
		r.Get("/w/auction/user/{id:[0-9]*}/lots", utility.MakeHandler(a.webLot.GetUserLotsHandler))
		r.Mount("/w/auction/lots", webLotRouter)
		r.Mount("/auction", wsLotsRouter)
	})
	return router
}
