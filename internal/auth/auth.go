package auth

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/session"
	"gitlab.com/bobayka/courseproject/internal/user"
	"gitlab.com/bobayka/courseproject/pkg/myerr"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"regexp"
	"time"
)

var userDatabase = make(map[string]*user.User)
var sessionDatabase = make(map[string]*session.Session)
var Bearer = regexp.MustCompile("(b|B)earer:([ a-z]*)")
var idUser int

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(int('a') + rand.Intn('z'-'a'))
	}
	if sessionDatabase[string(bytes)] != nil {
		return RandomString(len)
	}
	return string(bytes)
}

func addUser(u *request.RegUser) error {
	idUser++
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrapf(myerr.UnprocessableEntity, "can't generate password hash: %s", err)
		//return myerr.NewErr(
		//	errors.Wrap(err, "can't generate password hash"),
		//	"wrong password",
		//	422)
	}
	userDatabase[u.Email] = &user.User{
		ID:        idUser,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Birthday:  u.Birthday,
		Email:     u.Email,
		Password:  string(passwordHash),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	fmt.Printf("AddUser:\n%+v", *userDatabase[u.Email])
	return nil
}

func RegisterUser(u *request.RegUser) error {
	if userDatabase[u.Email] != nil {
		return errors.Wrap(myerr.BadRequest, "user is already registered")
	}
	err := addUser(u)
	if err != nil {
		return errors.Wrap(err, "user can't be add")
	}
	return myerr.Created
}

func checkValidToken(u request.TokenGetter) (*session.Session, error) {
	var s *session.Session
	if Bearer.FindStringSubmatch(u.GetToken()) == nil {
		return nil, errors.Wrap(myerr.UnprocessableEntity, "token doesnt match pattern `^(b|B)earer:([ a-z]*)$")
	}
	if s = sessionDatabase[u.GetToken()[len("bearer:"):]]; s != nil && s.CheckTokenTime() {
		return s, errors.Wrap(myerr.Accepted, "Accepted")
	}
	return nil, errors.Wrap(myerr.Unauthorized, "Unauthorized")
}

func addSession(db *user.User) string {
	token := RandomString(20)
	sessionDatabase[token] = &session.Session{
		SessionId:  token,
		UserID:     db.ID,
		CreatedAt:  time.Now(),
		ValidUntil: time.Now().AddDate(0, 3, 0),
	}
	fmt.Printf("AddSession:\n%+v", *sessionDatabase[token])
	return token
}

func AuthorizeUser(u *request.AuthUser) error {
	_, err := checkValidToken(u)
	if errors.Cause(err) != myerr.Unauthorized {
		return err
	}
	var db *user.User
	if db = userDatabase[u.Email]; db == nil {
		return errors.Wrap(myerr.BadRequest, "wrong email")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(db.Password), []byte(u.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword || err == bcrypt.ErrHashTooShort {
			return errors.Wrap(myerr.BadRequest, "wrong password")
		} else {
			return errors.Wrapf(myerr.BadRequest, "password cant be compared: %s", err)
		}
	}

	token := addSession(db)
	return errors.Wrap(myerr.Accepted, "Bearer: "+token)

}

func getUserByID(id int) *user.User {
	for _, v := range userDatabase {
		if id == v.ID {
			return v
		}
	}
	return nil
}

func update(db *user.User, u *request.UpdateUser) {
	db.FirstName = u.FirstName
	db.LastName = u.LastName
	if u.Birthday != nil && !u.Birthday.IsZero() {
		db.Birthday = u.Birthday
	}
	db.UpdatedAt = time.Now()
}

func UpdateUser(u *request.UpdateUser) error {
	s, err := checkValidToken(u)
	if errors.Cause(err) != myerr.Accepted {
		return err
	}
	db := getUserByID(s.UserID)
	update(db, u)
	fmt.Printf("UpdateUser:\n%+v", userDatabase[db.Email])
	return errors.Wrap(myerr.Success, "Update successful")
}
