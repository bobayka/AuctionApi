package postgres

import (
	"database/sql"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/session"
	"gitlab.com/bobayka/courseproject/internal/user"
	"gitlab.com/bobayka/courseproject/pkg/customTime"
	"gitlab.com/bobayka/courseproject/pkg/myerr"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

const (
	findUserByIDQuery              = `SELECT * FROM users WHERE id = $1`
	userInsertFields               = `first_name, last_name, email, password, birthday`
	insertUserQuery                = `INSERT INTO users(` + userInsertFields + `) VALUES ($1, $2, $3, $4, $5)`
	findUserByEmailQuery           = `SELECT * FROM users WHERE email = $1`
	sessionInsertFields            = `session_id, user_id`
	sessionInsertQuery             = `INSERT INTO sessions (` + sessionInsertFields + `) VALUES ($1, $2)`
	findSessionByTokenQuery        = ` SELECT * FROM sessions WHERE session_id = $1`
	userUpdateFields               = `first_name = $1, last_name = $2, birthday = $3 `
	updateUserQuery                = `UPDATE users SET ` + userUpdateFields + `WHERE id = $4`
	userUpdateFieldsWithoutBirtday = `first_name = $1, last_name = $2 `
	updateUserQueryWithoutBirtday  = `UPDATE users SET ` + userUpdateFieldsWithoutBirtday + `WHERE id = $3`
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

func scanUser(scanner sqlScanner, u *user.User) error {
	return scanner.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.Birthday,
		&u.CreatedAt, &u.UpdatedAt,
	)
}
func scanSession(scanner sqlScanner, s *session.Session) error {
	return scanner.Scan(&s.SessionId, &s.UserID, &s.CreatedAt, &s.ValidUntil)
}

type UsersStorage struct {
	statementStorage

	findUserByIDStmt         *sql.Stmt
	insertUserStmt           *sql.Stmt
	findUserByEmailStmt      *sql.Stmt
	insertSessionStmt        *sql.Stmt
	findSessionByTokenStmt   *sql.Stmt
	updateUser               *sql.Stmt
	updateUserWithoutBirtday *sql.Stmt
}

func NewUsersStorage(db *sql.DB) (*UsersStorage, error) {
	storage := &UsersStorage{statementStorage: newStatementsStorage(db)}

	statements := []stmt{
		{Query: insertUserQuery, Dst: &storage.insertUserStmt},
		{Query: findUserByEmailQuery, Dst: &storage.findUserByEmailStmt},
		{Query: sessionInsertQuery, Dst: &storage.insertSessionStmt},
		{Query: findSessionByTokenQuery, Dst: &storage.findSessionByTokenStmt},
		{Query: findUserByIDQuery, Dst: &storage.findUserByIDStmt},
		{Query: updateUserQuery, Dst: &storage.updateUser},
		{Query: updateUserQueryWithoutBirtday, Dst: &storage.updateUserWithoutBirtday},
	}

	if err := storage.initStatements(statements); err != nil {
		return nil, errors.Wrap(err, "can't create statements")
	}

	return storage, nil
}

func (s *UsersStorage) FindUserByEmail(email string) (*user.User, error) {
	var u user.User
	row := s.findUserByEmailStmt.QueryRow(email)
	if err := scanUser(row, &u); err != nil {
		return nil, errors.Wrap(err, "can't check user by email")
	}
	return &u, nil
}

func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func TimeToNullString(t *customTime.CustomTime) sql.NullString {
	if t.Time.UnixNano() == customTime.NilTime {
		return sql.NullString{}
	}
	return sql.NullString{String: t.Time.Format(customTime.CTLayout), Valid: true}
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

func (s *UsersStorage) AddSession(db *user.User) (string, error) {
	token := RandomString(20)
	_, err := s.insertSessionStmt.Exec(token, db.ID)
	if err != nil {
		return "", errors.Wrap(err, "Can't insert user")
	}
	return token, nil
}
func (s *UsersStorage) FindSessionByToken(token string) (*session.Session, error) {
	var ses session.Session
	row := s.findSessionByTokenStmt.QueryRow(token)
	if err := scanSession(row, &ses); err != nil {
		return nil, errors.Wrap(err, "can't check session by token")
	}
	return &ses, nil
}

func (s *UsersStorage) FindUserByID(id int64) (*user.User, error) {
	var u user.User
	row := s.findUserByIDStmt.QueryRow(id)
	if err := scanUser(row, &u); err != nil {
		return nil, errors.Wrapf(err, "can't scan user by id %d", id)
	}

	return &u, nil
}

func (s *UsersStorage) UpdateUserBD(id int64, u *request.UpdateUser) error {
	if u.Birthday.Time.UnixNano() == customTime.NilTime {
		_, err := s.updateUserWithoutBirtday.Exec(u.FirstName, u.LastName, id)
		if err != nil {
			return errors.Wrap(err, "Can't update user without birthday")
		}
		return nil
	}
	_, err := s.updateUser.Exec(u.FirstName, u.LastName, u.Birthday.Time, id)
	if err != nil {
		return errors.Wrap(err, "Can't update user bd")
	}
	return nil
}
