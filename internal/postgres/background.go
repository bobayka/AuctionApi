package postgres

import (
	"gitlab.com/bobayka/courseproject/internal/domains"
	"log"
	"time"
)

type Background struct {
	stmt *UsersStorage
}

func (b *Background) FinishEndedLots(d time.Duration) {
	for _ = range time.Tick(d) {
		var lots []domains.Lot
		rows, err := b.stmt.findAllLotsStmt.Query()
		if err != nil {
			log.Printf("Background: Can't select lots bd: %s", err)
		}
		lots, err = rowsLotsToSlice(rows, lots)
		if err != nil {
			log.Printf("Background: error in rows lot to slice: %s", err)
		}

	}

}
