package storage

import (
	"database/sql"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/internal/domains"
	request "gitlab.com/bobayka/courseproject/internal/requests"
	"golang.org/x/crypto/bcrypt"
)

const (
	findUserQuery        = `SELECT * FROM users WHERE `
	findUserByIDQuery    = findUserQuery + `id = $1`
	userInsertFields     = `first_name, last_name, email, password, birthday`
	insertUserQuery      = `INSERT INTO users(` + userInsertFields + `) VALUES ($1, $2, $3, $4, $5)`
	findUserByEmailQuery = findUserQuery + `email = $1`
	userUpdateFields     = `first_name = $1, last_name = $2, birthday = $3, updated_at = NOW() `
	updateUserQuery      = `UPDATE users SET ` + userUpdateFields + `WHERE id = $4`
)

type UsersStorage struct {
	statementStorage

	findUserByIDStmt    *sql.Stmt
	insertUserStmt      *sql.Stmt
	findUserByEmailStmt *sql.Stmt
	updateUserStmt      *sql.Stmt
}

func NewUsersStorage(db *sql.DB) (*UsersStorage, error) {
	storage := &UsersStorage{statementStorage: newStatementsStorage(db)}
	statements := []stmt{
		{Query: insertUserQuery, Dst: &storage.insertUserStmt},
		{Query: findUserByEmailQuery, Dst: &storage.findUserByEmailStmt},
		{Query: findUserByIDQuery, Dst: &storage.findUserByIDStmt},
		{Query: updateUserQuery, Dst: &storage.updateUserStmt},
	}
	if err := storage.initStatements(statements); err != nil {
		return nil, errors.Wrap(err, "can't create statements")
	}
	return storage, nil
}
func scanUser(scanner sqlScanner, u *domains.User) error {
	return scanner.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.Birthday,
		&u.CreatedAt, &u.UpdatedAt)
}

func (u *UsersStorage) FindUserByEmail(email string) (*domains.User, error) {
	var user domains.User
	row := u.findUserByEmailStmt.QueryRow(email)
	if err := scanUser(row, &user); err != nil {
		return nil, errors.Wrap(err, "can't check user by email")
	}
	return &user, nil
}

func (u *UsersStorage) AddUser(user *request.RegUser) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrapf(myerr.ErrBadRequest, "$can't generate password hash$: %u", err)
	}
	_, err = u.insertUserStmt.Exec(user.FirstName, user.LastName, user.Email,
		string(passwordHash), TimeToNullString(user.Birthday))
	if err != nil {
		return errors.Wrap(err, "Can't insert user")
	}
	return nil
}

func (u *UsersStorage) FindUserByID(id int64) (*domains.User, error) {
	var user domains.User
	row := u.findUserByIDStmt.QueryRow(id)
	if err := scanUser(row, &user); err != nil {
		return nil, errors.Wrapf(err, "can't scan user by id %d", id)
	}
	return &user, nil
}

func (u *UsersStorage) UpdateUserBD(id int64, user *request.UpdateUser) error {
	_, err := u.updateUserStmt.Exec(user.FirstName, user.LastName, TimeToNullString(user.Birthday), id)
	if err != nil {
		return errors.Wrap(err, "Can't update user bd")
	}
	return nil
}
