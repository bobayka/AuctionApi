package postgres

import (
	"database/sql"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/internal/domains"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/pkg/customTime"
	"gitlab.com/bobayka/courseproject/pkg/myerr"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

const (
	findUserByIDQuery       = `SELECT * FROM users WHERE id = $1`
	userInsertFields        = `first_name, last_name, email, password, birthday`
	insertUserQuery         = `INSERT INTO users(` + userInsertFields + `) VALUES ($1, $2, $3, $4, $5)`
	findUserByEmailQuery    = `SELECT * FROM users WHERE email = $1`
	sessionInsertFields     = `session_id, user_id`
	sessionInsertQuery      = `INSERT INTO sessions (` + sessionInsertFields + `) VALUES ($1, $2)`
	findSessionByTokenQuery = `SELECT * FROM sessions WHERE session_id = $1`
	userUpdateFields        = `first_name = $1, last_name = $2, birthday = $3`
	updateUserQuery         = `UPDATE users SET ` + userUpdateFields + `WHERE id = $4`
	lotInsertFields         = `title, description, min_price, price_step, end_at, creator_id, buyer_id`
	insertLotQuery          = `INSERT INTO lots(` + lotInsertFields + `) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	findLotFields           = `l.id, l.title, l.description, l.min_price, l.price_step, l.status, l.end_at, l.created_at, l.updated_at, u.id, u.first_name, u.last_name, d.id, d.first_name, d.last_name`
	findLotByIDQuery        = `SELECT ` + findLotFields + ` FROM lots as l INNER JOIN users as u ON l.creator_id = u.id  INNER JOIN  users as d ON l.buyer_id = d.id where l.id =  $1` //inner join
	lotUpdateFields         = `title = $1, description = $2, min_price = $3, price_step = $4, end_at = $5`
	updateLotQuery          = `UPDATE lots SET ` + lotUpdateFields + ` WHERE id = $6 AND status = 'created'`
	//selectAllLotsQuery      = `SELECT `+ findLotFields+ ` FROM lots as l INNER JOIN users as u ON l.creator_id = u.id  INNER JOIN  users as d ON l.buyer_id = d.id`
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(int('a') + rand.Intn('z'-'a'))
	}
	return string(bytes)
}

func scanUser(scanner sqlScanner, u *domains.User) error {
	return scanner.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.Birthday,
		&u.CreatedAt, &u.UpdatedAt,
	)
}
func scanSession(scanner sqlScanner, s *domains.Session) error {
	return scanner.Scan(&s.SessionId, &s.UserID, &s.CreatedAt, &s.ValidUntil)
}

func scanLot(scanner sqlScanner, s *domains.Lot) error {
	return scanner.Scan(&s.ID, &s.Title, &s.Description, &s.MinPrice,
		&s.PriceStep, &s.Status, &s.EndAt, &s.CreatedAt, &s.UpdatedAt,
		&s.CreatorID.ID, &s.CreatorID.FirstName, &s.CreatorID.LastName,
		&s.BuyerID.ID, &s.BuyerID.FirstName, &s.BuyerID.LastName)
}

func TimeToNullString(t *customTime.CustomTime) sql.NullString {
	if t == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: time.Time(*t).Format(customTime.CTLayout), Valid: true}
}

type UsersStorage struct {
	statementStorage

	findUserByIDStmt       *sql.Stmt
	insertUserStmt         *sql.Stmt
	findUserByEmailStmt    *sql.Stmt
	insertSessionStmt      *sql.Stmt
	findSessionByTokenStmt *sql.Stmt
	updateUserStmt         *sql.Stmt
	insertLotStmt          *sql.Stmt
	findLotByIDStmt        *sql.Stmt
	updateLotStmt          *sql.Stmt
}

func NewUsersStorage(db *sql.DB) (*UsersStorage, error) {
	storage := &UsersStorage{statementStorage: newStatementsStorage(db)}

	statements := []stmt{
		{Query: insertUserQuery, Dst: &storage.insertUserStmt},
		{Query: findUserByEmailQuery, Dst: &storage.findUserByEmailStmt},
		{Query: sessionInsertQuery, Dst: &storage.insertSessionStmt},
		{Query: findSessionByTokenQuery, Dst: &storage.findSessionByTokenStmt},
		{Query: findUserByIDQuery, Dst: &storage.findUserByIDStmt},
		{Query: updateUserQuery, Dst: &storage.updateUserStmt},
		{Query: insertLotQuery, Dst: &storage.insertLotStmt},
		{Query: findLotByIDQuery, Dst: &storage.findLotByIDStmt},
		{Query: updateLotQuery, Dst: &storage.updateLotStmt},
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
		return errors.Wrapf(myerr.BadRequest, "can't generate password hash: %s", err)
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

func (s *UsersStorage) InsertLot(UserID int64, l *request.LotToCreateUpdate) (int64, error) {
	var LotID int64
	err := s.insertLotStmt.QueryRow(l.Title, l.Description, l.MinPrice,
		l.PriceStep, l.EndAt, UserID, UserID).Scan(&LotID)
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
		l.PriceStep, l.EndAt, id)
	if err != nil {
		return errors.Wrap(err, "Can't update lot bd")
	}
	count, err2 := res.RowsAffected()
	if err2 != nil {
		return errors.Wrap(err, "Can't get rows affected")
	}
	if count == 0 {
		return myerr.NotFound
	}
	return nil
}

//func (s *UsersStorage) SelectAllLotsDB() ([]domains.Lot, error) {
//	rows, err := s.selectAllLotsStmt.Query()
//	if err != nil {
//		return nil, errors.Wrap(err, "Can't select lots bd")
//	}
//	defer rows.Close()
//	var lots []domains.Lot
//	for rows.Next() {
//		var lot domains.Lot
//		if err := scanLot(rows, &lot); err != nil {
//			return nil, errors.Wrap(err, "can't scan lots")
//		}
//		lots = append(lots, lot)
//	}
//	if err = rows.Err(); err != nil {
//		return nil, errors.Wrap(err, "rows return error")
//	}
//	return lots, nil
//
//}
