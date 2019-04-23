package storage

import (
	"database/sql"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/internal/domains"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/responce"
	"gitlab.com/bobayka/courseproject/pkg/customTime"
	"log"
	"math/rand"
	"time"
)

//nolint: gosec
const (
	findShortUserByIDQuery   = `SELECT id, first_name, last_name FROM users WHERE id = $1`
	lotInsertFields          = `title, description, min_price, price_step, end_at, status, creator_id `
	insertLotQuery           = `INSERT INTO lots(` + lotInsertFields + `) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	findLotFields            = `id, title, description, buy_price, min_price, price_step, status, end_at, created_at, updated_at, creator_id, buyer_id, deleted_at`
	findLotsQuery            = `SELECT ` + findLotFields + ` FROM lots `
	findAllLotsOrderbyQuery  = findLotsQuery + `order by id`
	findLotByStatusQuery     = findLotsQuery + `where status = $1 ORDER BY ID `
	findLotByIDQuery         = findLotsQuery + `where id =  $1 ORDER BY id`
	findLotsByCreatorIDQuery = findLotsQuery + `where creator_id =  $1 ORDER BY id`
	findLotsByBuyerIDQuery   = findLotsQuery + `where buyer_id =  $1 ORDER BY id`
	findAllUserLots          = findLotsQuery + `where buyer_id = $1 OR creator_id =  $1 ORDER BY id`
	lotUpdateFields          = `title = $1, description = $2, min_price = $3, price_step = $4, end_at = $5, status = $6 ,updated_at = NOW()`
	updateLotQuery           = `UPDATE lots SET ` + lotUpdateFields + ` WHERE id = $7 AND status = 'created'`
	deleteLotQuery           = `UPDATE lots SET deleted_at = NOW() WHERE ID = $1`
	updateLotPriceQuery      = `UPDATE lots SET buy_price = $1, buyer_id = $2, updated_at = NOW() WHERE id = $3`
	updateFinishedLotsQuery  = `UPDATE lots SET status = 'finished', updated_at = NOW() WHERE NOW()>end_at and NOT status='finished' RETURNING id`
)

func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(int('a') + rand.Intn('z'-'a'))
	}
	return string(bytes)
}

func rowsLotsToSlice(rows *sql.Rows, lots []*domains.Lot) ([]*domains.Lot, error) {
	for rows.Next() {
		var lot domains.Lot
		if err := scanLot(rows, &lot); err != nil {
			return nil, errors.Wrapf(err, "can't scan lot")
		}
		lots = append(lots, &lot)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows return error")
	}
	return lots, nil
}

func scanShortUser(scanner sqlScanner, u *responce.ShortUser) error {
	return scanner.Scan(&u.ID, &u.FirstName, &u.LastName)
}

func scanLot(scanner sqlScanner, s *domains.Lot) error {
	return scanner.Scan(&s.ID, &s.Title, &s.Description, &s.BuyPrice, &s.MinPrice,
		&s.PriceStep, &s.Status, &s.EndAt, &s.CreatedAt, &s.UpdatedAt,
		&s.CreatorID, &s.BuyerID, &s.DeletedAt)
}
func TimeToNullString(t *customtime.CustomTime) sql.NullString {
	if t == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: time.Time(*t).Format(customtime.CTLayout), Valid: true}
}

type LotsStorage struct {
	statementStorage
	findShortUserByIDStmt   *sql.Stmt
	insertLotStmt           *sql.Stmt
	findLotByIDStmt         *sql.Stmt
	updateLotStmt           *sql.Stmt
	deleteLotStmt           *sql.Stmt
	findLotsByCreatorIDStmt *sql.Stmt
	findLotsByBuyerIDStmt   *sql.Stmt
	findAllUserLotsStmt     *sql.Stmt
	findLotByStatusStmt     *sql.Stmt
	findAllLotsStmt         *sql.Stmt
	updateLotPriceStmt      *sql.Stmt
	updateFinishedLotsStmt  *sql.Stmt
}

func NewLotsStorage(db *sql.DB) (*LotsStorage, error) {
	storage := &LotsStorage{statementStorage: newStatementsStorage(db)}

	statements := []stmt{

		{Query: findShortUserByIDQuery, Dst: &storage.findShortUserByIDStmt},
		{Query: insertLotQuery, Dst: &storage.insertLotStmt},
		{Query: findLotByIDQuery, Dst: &storage.findLotByIDStmt},
		{Query: updateLotQuery, Dst: &storage.updateLotStmt},
		{Query: deleteLotQuery, Dst: &storage.deleteLotStmt},
		{Query: findAllUserLots, Dst: &storage.findAllUserLotsStmt},
		{Query: findAllLotsOrderbyQuery, Dst: &storage.findAllLotsStmt},
		{Query: findLotsByBuyerIDQuery, Dst: &storage.findLotsByBuyerIDStmt},
		{Query: findLotsByCreatorIDQuery, Dst: &storage.findLotsByCreatorIDStmt},
		{Query: findLotByStatusQuery, Dst: &storage.findLotByStatusStmt},
		{Query: updateLotPriceQuery, Dst: &storage.updateLotPriceStmt},
		{Query: updateFinishedLotsQuery, Dst: &storage.updateFinishedLotsStmt},
	}

	if err := storage.initStatements(statements); err != nil {
		return nil, errors.Wrap(err, "can't create statements")
	}
	return storage, nil
}

func (s *LotsStorage) FindShortUserByID(id int64) (*responce.ShortUser, error) {
	var u responce.ShortUser
	row := s.findShortUserByIDStmt.QueryRow(id)
	if err := scanShortUser(row, &u); err != nil {
		return nil, errors.Wrapf(err, "can't scan ShortUser by id %d", id)
	}
	return &u, nil
}

func (s *LotsStorage) InsertLot(userID int64, l *request.LotCreateUpdate) (int64, error) {
	var LotID int64
	err := s.insertLotStmt.QueryRow(l.Title, l.Description, l.MinPrice,
		l.PriceStep, l.EndAt, l.Status, userID).Scan(&LotID)
	if err != nil {
		return 0, errors.Wrap(err, "Can't insert lot in bd")
	}
	return LotID, nil
}

func (s *LotsStorage) FindLotByID(id int64) (*domains.Lot, error) {
	var l domains.Lot
	row := s.findLotByIDStmt.QueryRow(id)
	if err := scanLot(row, &l); err != nil {
		return nil, errors.Wrapf(err, "can't scan lot by id %d", id)
	}
	return &l, nil
}

func (s *LotsStorage) UpdateLotBD(id int64, l *request.LotCreateUpdate) error {
	res, err := s.updateLotStmt.Exec(l.Title, l.Description, l.MinPrice,
		l.PriceStep, l.EndAt, l.Status, id)
	if err != nil {
		return errors.Wrap(err, "Can't update lot bd")
	}
	count, err2 := res.RowsAffected()
	if err2 != nil {
		return errors.Wrap(err, "Can't get affected rows")
	}
	if count == 0 {
		return myerr.ErrNotFound
	}
	return nil
}

func (s *LotsStorage) DeleteLotBD(id int64) error {
	res, err := s.deleteLotStmt.Exec(id)
	if err != nil {
		return errors.Wrap(err, "Can't delete lot bd")
	}
	count, err2 := res.RowsAffected()
	if err2 != nil {
		return errors.Wrap(err, "Can't get affected rows ")
	}
	if count == 0 {
		return myerr.ErrNotFound
	}
	return nil
}

func (s *LotsStorage) FindUserLotsBD(userID int64, lotsType string) ([]*domains.Lot, error) {
	var rows *sql.Rows
	var err error //	fmt.Printf("%+v", dbLots)
	switch lotsType {
	case "own":
		rows, err = s.findLotsByCreatorIDStmt.Query(userID)
	case "buyed":
		rows, err = s.findLotsByBuyerIDStmt.Query(userID)
	case "":
		rows, err = s.findAllUserLotsStmt.Query(userID)
	default:
		return nil, errors.New("query param doesnt match") // временнное
	}

	if err != nil {
		return nil, errors.Wrap(err, "Can't select lots bd")
	}
	defer rows.Close()
	var lots []*domains.Lot
	lots, err = rowsLotsToSlice(rows, lots)
	if err != nil {
		return nil, errors.Wrap(err, "error in rows lot to slice")
	}
	return lots, nil

}

func (s *LotsStorage) BackgroundUpdateLotsBD() ([]*int64, error) {
	rows, err := s.updateFinishedLotsStmt.Query()
	if err != nil {
		log.Printf("Can't update lots bd: %s", err)
	}
	defer rows.Close()
	var lotsID []*int64
	for rows.Next() {
		var lotID int64
		if err := rows.Scan(&lotID); err != nil {
			return nil, errors.Wrapf(err, "can't scan lot")
		}
		lotsID = append(lotsID, &lotID)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows return error")
	}
	if err != nil {
		return nil, errors.Wrap(err, "error in rows lot to slice")
	}
	return lotsID, nil

}
func (s *LotsStorage) FindAllLotsBD(status string) ([]*domains.Lot, error) {
	var rows *sql.Rows
	var err error //	fmt.Printf("%+v", dbLots)
	if status != "" {
		rows, err = s.findLotByStatusStmt.Query(status)
	} else {
		rows, err = s.findAllLotsStmt.Query()
	}
	if err != nil {
		return nil, errors.Wrap(err, "Can't select lots bd")
	}
	defer rows.Close()
	var lots []*domains.Lot
	lots, err = rowsLotsToSlice(rows, lots)
	if err != nil {
		return nil, errors.Wrap(err, "error in rows lot to slice")
	}
	return lots, nil

}

func (s *LotsStorage) UpdateLotPriceBD(userID int64, lotID int64, price float64) error {
	_, err := s.updateLotPriceStmt.Exec(price, userID, lotID)
	if err != nil {
		return errors.Wrap(err, "Can't update lot price bd")
	}
	return nil
}
