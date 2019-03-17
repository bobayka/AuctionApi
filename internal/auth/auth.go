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
var bearer = regexp.MustCompile("^(b|B)earer:([ a-z]*)$")
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

func addUser(u *request.RegUser) *myerr.AppError {
	idUser++
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return myerr.NewErr(
			errors.Wrap(err, "can't generate password hash"),
			"wrong password",
			422)
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

func RegisterUser(u *request.RegUser) *myerr.AppError {
	if userDatabase[u.Email] != nil {
		return myerr.NewErr(nil, "User is already registered", 422)
	}
	err := addUser(u)
	if err != nil {
		return err.MyWrap("user can't be add")
	}
	return myerr.NewErr(nil, "Registration complete", 200)
}

func checkValidToken(u request.TokenGetter) (*session.Session, *myerr.AppError) {
	var s *session.Session
	if bearer.FindStringSubmatch(u.GetToken()) == nil {
		return nil, myerr.NewErr(nil, "token doesnt match pattern `Bearer:*******` ", 449)
	}
	if s = sessionDatabase[u.GetToken()[len("bearer:"):]]; s != nil && s.CheckTokenTime() {
		return s, myerr.NewErr(nil, "Accept", 200)
	}
	return nil, myerr.NewErr(nil, "Wrong token", 401)
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

func AuthorizeUser(u *request.AuthUser) *myerr.AppError {
	_, err := checkValidToken(u)
	if err.Message != "Wrong token" {
		return err
	}
	var db *user.User
	if db = userDatabase[u.Email]; db == nil {
		return myerr.NewErr(nil, "Wrong email", 400)
	}
	if bcrypt.CompareHashAndPassword([]byte(db.Password), []byte(u.Password)) != nil {
		return myerr.NewErr(nil, "Wrong password", 400)
	}
	//if err := bcrypt.CompareHashAndPassword([]byte(db.Password), []byte(u.Password)); err != nil {
	//	err = errors.Wrap(err, "cant convert password to hash")
	//	return myerr.NewErr(err, "Wrong password", 400)
	//}
	token := addSession(db)
	return myerr.NewErr(nil, "Bearer: "+token, 200)

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

func UpdateUser(u *request.UpdateUser) *myerr.AppError {
	s, err := checkValidToken(u)
	if err.Message != "Accept" {
		return err
	}
	db := getUserByID(s.UserID)
	update(db, u)
	fmt.Printf("UpdateUser:\n%+v", userDatabase[db.Email])
	return myerr.NewErr(nil, "Update successful", 200)
}
