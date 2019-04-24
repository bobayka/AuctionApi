package gatewayapi

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gitlab.com/bobayka/courseproject/cmd/Protobuf"
	"gitlab.com/bobayka/courseproject/cmd/gateway-api/handlers/htmlhandlers/auth-handlers"
	"gitlab.com/bobayka/courseproject/cmd/gateway-api/handlers/htmlhandlers/lot-handlers"
	"gitlab.com/bobayka/courseproject/cmd/gateway-api/handlers/jsonhandlers/auth-handlers"
	"gitlab.com/bobayka/courseproject/cmd/gateway-api/handlers/jsonhandlers/lot-handlers"
	"gitlab.com/bobayka/courseproject/cmd/gateway-api/handlers/jsonhandlers/user-handlers"
	"gitlab.com/bobayka/courseproject/cmd/gateway-api/handlers/jsonhandlers/websocket-handlers"
	"gitlab.com/bobayka/courseproject/cmd/utilities"
	"gitlab.com/bobayka/courseproject/internal/postgres/storage"
	"time"
)

type GatewayAPI struct {
	//lotServiceClient lotspb.LotsServiceClient
	storage *storage.SessionStorage
	auth    *authhandlers.AuthHandler
	lot     *lothandlers.LotServiceHandler
	user    *userhandlers.UserHandler
	webAuth *authweb.WebAuthHandler
	webLot  *lotweb.WebLotHandler
	wsLot   *websockethandlers.LotWSHandler
}

func NewGatewayAPI(storage storage.Storage, client lotspb.LotsServiceClient) *GatewayAPI {
	auth := authhandlers.NewAuthHandler(storage)
	lot := lothandlers.NewLotHandler(client)
	user := userhandlers.NewUserHandlers(storage.Users)
	webLot := lotweb.NewWebLotHandler(client)
	webAuth := authweb.NewWebAuthHandler(storage)
	wsLot := websockethandlers.NewLotWSHandler(client)
	wsLot.BackGroundFinishEndedLots(time.Second)
	return &GatewayAPI{storage: storage.Sessions, auth: auth, lot: lot, user: user, webLot: webLot, webAuth: webAuth, wsLot: wsLot}

}

func (a *GatewayAPI) Routes() *chi.Mux {
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
