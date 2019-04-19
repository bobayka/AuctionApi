package PGBackground

import (
	"encoding/json"
	"fmt"
	"gitlab.com/bobayka/courseproject/cmd/auth-api/handlers/JSONHandlers/websocket-handlers"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"log"
	"time"
)

type Background struct {
	hub  *websocket_handlers.Hub
	stmt *postgres.UsersStorage
}

func StartDBBackgroundProcesses(stmt *postgres.UsersStorage, hub *websocket_handlers.Hub) {
	back := Background{hub: hub, stmt: stmt}
	back.FinishEndedLots(time.Second)
}
func (b *Background) FinishEndedLots(d time.Duration) {
	go func() {
		for range time.Tick(d) {
			lots, err := b.stmt.BackgroundUpdateLotsBD()
			if err != nil {
				log.Printf("Background: %s", err)
			}
			if lots == nil {
				continue
			}
			for _, lot := range lots {
				res, err := json.Marshal(lot)
				if err != nil {
					log.Printf("Background: can't marshal message: %s", err)
				}
				fmt.Println(string(res))
				b.hub.Broadcast <- websocket_handlers.BroadcastWithID{LotID: lot.ID, Data: res}
				b.hub.Broadcast <- websocket_handlers.BroadcastWithID{LotID: -1, Data: res}
			}
		}
	}()

}
