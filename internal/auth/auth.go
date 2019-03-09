package auth

import (
	"crypto/md5"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/internal/requests"
	"gitlab.com/bobayka/courseproject/internal/session"
	"gitlab.com/bobayka/courseproject/internal/user"
	"math/rand"
	"regexp"
	"time"
)

var userDatabase = make(map[string]*user.User)
var sessionDatabase = make(map[string]*session.Session)
var bearer = regexp.MustCompile("^(b|B)earer:([ a-z]*)$")
var idUser int

func RandomString(len int) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(97 + rand.Intn(25)) //a=65 and z = 97+25
	}
	if sessionDatabase[string(bytes)] != nil {
		return RandomString(len)
	}
	return string(bytes)
}

func hashString(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func addUser(u *request.RegUser) {
	idUser++
	userDatabase[u.Email] = &user.User{
		ID:        idUser,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Birthday:  u.Birthday,
		Email:     u.Email,
		Password:  hashString(u.Password),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	fmt.Printf("AddUser:\n%+v", *userDatabase[u.Email])
}

func RegisterUser(u *request.RegUser) error {
	if userDatabase[u.Email] != nil {
		return errors.New("User is already registered")
	}
	addUser(u)
	return errors.New("Registration complete")
}

func checkValidToken(u request.User) (*session.Session, error) {
	var s *session.Session
	if bearer.FindStringSubmatch(u.GetToken()) == nil {
		return nil, errors.New("token doesnt match pattern 	`Bearer:*******` ")
	}
	if s = sessionDatabase[u.GetToken()[len("bearer:"):]]; s != nil && s.CheckTokenTime() {
		return s, errors.New("Accept")
	}
	return nil, errors.New("Wrong token")
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
	if err.Error() != "Wrong token" {
		return err
	}
	var db *user.User
	if db = userDatabase[u.Email]; db == nil {
		return errors.New("Wrong email")
	}
	if hashString(u.Password) != db.Password {
		return errors.New("Wrong password")
	}
	token := addSession(db)
	return errors.New("Bearer: " + token)

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
	db.Birthday = u.Birthday
	db.UpdatedAt = time.Now()
}

func UpdateUser(u *request.UpdateUser) error {
	s, err := checkValidToken(u)
	if err.Error() != "Accept" {
		return err
	}
	db := getUserByID(s.UserID)
	update(db, u)
	fmt.Printf("UpdateUser:\n%+v", userDatabase[db.Email])
	return errors.New("Update successful")
}
