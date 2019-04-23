package storage

import (
	"database/sql"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/internal/domains"
)

const (
	sessionInsertFields     = `session_id, user_id`
	sessionInsertQuery      = `INSERT INTO sessions (` + sessionInsertFields + `) VALUES ($1, $2)`
	findSessionByTokenQuery = `SELECT * FROM sessions WHERE session_id = $1`
)

type SessionStorage struct {
	statementStorage
	insertSessionStmt      *sql.Stmt
	findSessionByTokenStmt *sql.Stmt
}

func NewSessionsStorage(db *sql.DB) (*SessionStorage, error) {
	storage := &SessionStorage{statementStorage: newStatementsStorage(db)}
	statements := []stmt{
		{Query: sessionInsertQuery, Dst: &storage.insertSessionStmt},
		{Query: findSessionByTokenQuery, Dst: &storage.findSessionByTokenStmt},
	}
	if err := storage.initStatements(statements); err != nil {
		return nil, errors.Wrap(err, "can't create statements")
	}
	return storage, nil
}

func scanSession(scanner sqlScanner, s *domains.Session) error {
	return scanner.Scan(&s.SessionID, &s.UserID, &s.CreatedAt, &s.ValidUntil)
}

func (s *SessionStorage) AddSession(db *domains.User) (string, error) {
	token := RandomString(20)
	_, err := s.insertSessionStmt.Exec(token, db.ID)
	if err != nil {
		return "", errors.Wrap(err, "Can't insert user") // можно добавить проверку  на уникальность токена
	}
	return token, nil
}

func (s *SessionStorage) FindSessionByToken(token string) (*domains.Session, error) {
	var ses domains.Session
	row := s.findSessionByTokenStmt.QueryRow(token)
	if err := scanSession(row, &ses); err != nil {
		return nil, errors.Wrap(err, "can't check session by token")
	}
	return &ses, nil
}
