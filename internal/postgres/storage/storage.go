package storage

import (
	"database/sql"
	"github.com/pkg/errors"
)

type Storage struct {
	Lots     *LotsStorage
	Sessions *SessionStorage
	Users    *UsersStorage
}

func NewStorage(db *sql.DB) (Storage, error) {
	usersStorage, err := NewUsersStorage(db)
	if err != nil {
		return Storage{}, errors.Wrap(err, "error in creation users storage")
	}
	lotsStorage, err := NewLotsStorage(db)
	if err != nil {
		return Storage{}, errors.Wrap(err, "error in creation lots storage")
	}
	sessionsStorage, err := NewSessionsStorage(db)
	if err != nil {
		return Storage{}, errors.Wrap(err, "error in creation sessions storage")
	}
	return Storage{Lots: lotsStorage, Sessions: sessionsStorage, Users: usersStorage}, nil
}

func (s *Storage) Close() {
	s.Lots.Close()
	s.Sessions.Close()
	s.Users.Close()
}
