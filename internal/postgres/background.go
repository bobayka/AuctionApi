package postgres

import (
	"log"
	"time"
)

type Background struct {
	stmt *UsersStorage
}

func StartDBBackgroundProcesses(stmt *UsersStorage) {
	back := Background{stmt: stmt}
	back.FinishEndedLots(time.Second)
}
func (b *Background) FinishEndedLots(d time.Duration) {
	go func() {
		for range time.Tick(d) {
			_, err := b.stmt.updateFinishedLotsStmt.Exec()
			if err != nil {
				log.Printf("Background: Can't update lots bd: %s", err)
			}
		}
	}()

}
