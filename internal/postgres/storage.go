package postgres

import (
	"database/sql"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/internal/domains"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/responce"
	"gitlab.com/bobayka/courseproject/pkg/customTime"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"time"
)

//nolint: gosec
const (
	findUserQuery           = `SELECT * FROM users WHERE `
	findUserByIDQuery       = findUserQuery + `id = $1`
	findShortUserByIDQuery  = `SELECT id, first_name, last_name FROM users WHERE id = $1`
	userInsertFields        = `first_name, last_name, email, password, birthday`
	insertUserQuery         = `INSERT INTO users(` + userInsertFields + `) VALUES ($1, $2, $3, $4, $5)`
	findUserByEmailQuery    = findUserQuery + `email = $1`
	sessionInsertFields     = `session_id, user_id`
	sessionInsertQuery      = `INSERT INTO sessions (` + sessionInsertFields + `) VALUES ($1, $2)`
	findSessionByTokenQuery = `SELECT * FROM sessions WHERE session_id = $1`
	userUpdateFields        = `first_name = $1, last_name = $2, birthday = $3, updated_at = NOW() `
	updateUserQuery         = `UPDATE users SET ` + userUpdateFields + `WHERE id = $4`
	lotInsertFields         = `title, description, min_price, price_step, end_at, creator_id`
	insertLotQuery          = `INSERT INTO lots(` + lotInsertFields + `) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	//findLotFields           = `l.id, l.title, l.description, l.buy_price, l.min_price, l.price_step, l.status, l.end_at, l.created_at, l.updated_at, u.id, u.first_name, u.last_name, d.id, d.first_name, d.last_name`
	//findLotsQuery            = `SELECT ` + findLotFields + ` FROM lots as l INNER JOIN users as u ON l.creator_id = u.id  LEFT JOIN  users as d ON l.buyer_id = d.id ` //inner join
	findLotFields            = `id, title, description, buy_price, min_price, price_step, status, end_at, created_at, updated_at, creator_id, buyer_id`
	findLotsQuery            = `SELECT ` + findLotFields + ` FROM lots `
	findLotByStatusQuery     = findLotsQuery + `where status = $1`
	findLotByIDQuery         = findLotsQuery + `where id =  $1`
	findLotsByCreatorIDQuery = findLotsQuery + `where creator_id =  $1`
	findLotsByBuyerIDQuery   = findLotsQuery + `where buyer_id =  $1`
	findAllUserLots          = findLotsQuery + `where buyer_id = $1 OR creator_id =  $1`
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

func scanShortUser(scanner sqlScanner, u *responce.ShortUSer) error {
	return scanner.Scan(&u.ID, &u.FirstName, &u.LastName)
}

func scanUser(scanner sqlScanner, u *domains.User) error {
	return scanner.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.Birthday,
		&u.CreatedAt, &u.UpdatedAt)
}
func scanSession(scanner sqlScanner, s *domains.Session) error {
	return scanner.Scan(&s.SessionID, &s.UserID, &s.CreatedAt, &s.ValidUntil)
}

func scanLot(scanner sqlScanner, s *domains.Lot) error {
	return scanner.Scan(&s.ID, &s.Title, &s.Description, &s.BuyPrice, &s.MinPrice,
		&s.PriceStep, &s.Status, &s.EndAt, &s.CreatedAt, &s.UpdatedAt,
		&s.CreatorID, &s.BuyerID)
}
func TimeToNullString(t *customtime.CustomTime) sql.NullString {
	if t == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: time.Time(*t).Format(customtime.CTLayout), Valid: true}
}

type UsersStorage struct {
	statementStorage

	findUserByIDStmt        *sql.Stmt
	findShortUserByIDStmt   *sql.Stmt
	insertUserStmt          *sql.Stmt
	findUserByEmailStmt     *sql.Stmt
	insertSessionStmt       *sql.Stmt
	findSessionByTokenStmt  *sql.Stmt
	updateUserStmt          *sql.Stmt
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

func NewUsersStorage(db *sql.DB) (*UsersStorage, error) {
	storage := &UsersStorage{statementStorage: newStatementsStorage(db)}

	statements := []stmt{
		{Query: insertUserQuery, Dst: &storage.insertUserStmt},
		{Query: findUserByEmailQuery, Dst: &storage.findUserByEmailStmt},
		{Query: findShortUserByIDQuery, Dst: &storage.findShortUserByIDStmt},
		{Query: sessionInsertQuery, Dst: &storage.insertSessionStmt},
		{Query: findSessionByTokenQuery, Dst: &storage.findSessionByTokenStmt},
		{Query: findUserByIDQuery, Dst: &storage.findUserByIDStmt},
		{Query: updateUserQuery, Dst: &storage.updateUserStmt},
		{Query: insertLotQuery, Dst: &storage.insertLotStmt},
		{Query: findLotByIDQuery, Dst: &storage.findLotByIDStmt},
		{Query: updateLotQuery, Dst: &storage.updateLotStmt},
		{Query: deleteLotQuery, Dst: &storage.deleteLotStmt},
		{Query: findAllUserLots, Dst: &storage.findAllUserLotsStmt},
		{Query: findLotsQuery, Dst: &storage.findAllLotsStmt},
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

func (s *UsersStorage) FindUserByEmail(email string) (*domains.User, error) {
	var u domains.User
	row := s.findUserByEmailStmt.QueryRow(email)
	if err := scanUser(row, &u); err != nil {
		return nil, errors.Wrap(err, "can't check user by email")
	}
	return &u, nil
}

func (s *UsersStorage) AddUser(u *request.RegUser) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrapf(myerr.ErrBadRequest, "$can't generate password hash$: %s", err)
	}
	_, err = s.insertUserStmt.Exec(u.FirstName, u.LastName, u.Email,
		string(passwordHash), TimeToNullString(u.Birthday))
	if err != nil {
		return errors.Wrap(err, "Can't insert user")
	}
	return nil
}

func (s *UsersStorage) AddSession(db *domains.User) (string, error) {
	token := RandomString(20)
	_, err := s.insertSessionStmt.Exec(token, db.ID)
	if err != nil {
		return "", errors.Wrap(err, "Can't insert user") // можно добавить проверку  на уникальность токена
	}
	return token, nil
}

func (s *UsersStorage) FindSessionByToken(token string) (*domains.Session, error) {
	var ses domains.Session
	row := s.findSessionByTokenStmt.QueryRow(token)
	if err := scanSession(row, &ses); err != nil {
		return nil, errors.Wrap(err, "can't check session by token")
	}
	return &ses, nil
}

func (s *UsersStorage) FindShortUserByID(id int64) (*responce.ShortUSer, error) {
	var u responce.ShortUSer
	row := s.findShortUserByIDStmt.QueryRow(id)
	if err := scanShortUser(row, &u); err != nil {
		return nil, errors.Wrapf(err, "can't scan ShortUser by id %d", id)
	}
	return &u, nil
}

func (s *UsersStorage) FindUserByID(id int64) (*domains.User, error) {
	var u domains.User
	row := s.findUserByIDStmt.QueryRow(id)
	if err := scanUser(row, &u); err != nil {
		return nil, errors.Wrapf(err, "can't scan user by id %d", id)
	}
	return &u, nil
}

func (s *UsersStorage) UpdateUserBD(id int64, u *request.UpdateUser) error {
	_, err := s.updateUserStmt.Exec(u.FirstName, u.LastName, TimeToNullString(u.Birthday), id)
	if err != nil {
		return errors.Wrap(err, "Can't update user bd")
	}
	return nil
}

func (s *UsersStorage) InsertLot(userID int64, l *request.LotToCreateUpdate) (int64, error) {
	var LotID int64
	err := s.insertLotStmt.QueryRow(l.Title, l.Description, l.MinPrice,
		l.PriceStep, l.EndAt, userID).Scan(&LotID)
	if err != nil {
		return 0, errors.Wrap(err, "Can't insert lot in bd")
	}
	return LotID, nil
}

func (s *UsersStorage) FindLotByID(id int64) (*domains.Lot, error) {
	var l domains.Lot
	row := s.findLotByIDStmt.QueryRow(id)
	if err := scanLot(row, &l); err != nil {
		return nil, errors.Wrapf(err, "can't scan lot by id %d", id)
	}
	return &l, nil
}

func (s *UsersStorage) UpdateLotBD(id int64, l *request.LotToCreateUpdate) error {
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

func (s *UsersStorage) DeleteLotBD(id int64) error {
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

func (s *UsersStorage) FindUserLotsBD(userID int64, lotsType string) ([]*domains.Lot, error) {
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

func (s *UsersStorage) ConvertLotToRespLot(dbLot *domains.Lot) (*responce.RespLot, error) {

	creator, err := s.FindShortUserByID(dbLot.CreatorID)
	if err != nil {
		return nil, err
	}
	var buyer *responce.ShortUSer
	if dbLot.BuyerID != nil {
		buyer, err = s.FindShortUserByID(*dbLot.BuyerID)
	}
	return &responce.RespLot{LotGeneral: dbLot.LotGeneral, Creator: *creator, Buyer: buyer}, nil
}

func (s *UsersStorage) BackgroundUpdateLotsBD() ([]*responce.RespLot, error) {
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
	var respLots []*responce.RespLot
	for _, v := range lotsID {
		dbLot, err := s.FindLotByID(*v)
		respLot, err := s.ConvertLotToRespLot(dbLot)
		if err != nil {
			return nil, err
		}
		respLots = append(respLots, respLot)
	}
	return respLots, nil
}
func (s *UsersStorage) FindAllLotsBD(status string) ([]*domains.Lot, error) {
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

func (s *UsersStorage) UpdateLotPriceBD(userID int64, lotID int64, price float64) error {
	_, err := s.updateLotPriceStmt.Exec(price, userID, lotID)
	if err != nil {
		return errors.Wrap(err, "Can't update lot price bd")
	}
	return nil
}
