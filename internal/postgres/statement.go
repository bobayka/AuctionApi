package postgres

import (
	"database/sql"
	"github.com/pkg/errors"
)

type statementStorage struct {
	db         *sql.DB
	statements []*sql.Stmt
}

func newStatementsStorage(db *sql.DB) statementStorage {
	s := statementStorage{
		db: db,
	}
	return s
}

type sqlScanner interface {
	Scan(dest ...interface{}) error
}

type stmt struct {
	Query string
	Dst   **sql.Stmt
}

func (s *statementStorage) initStatements(statements []stmt) error {
	for i := range statements {
		statement, err := s.db.Prepare(statements[i].Query)
		if err != nil {
			return errors.Wrapf(err, "can't prepare query %q", statements[i].Query)
		}

		*statements[i].Dst = statement
		s.statements = append(s.statements, statement)
	}

	return nil
}

// Close implements io.Closer interface. It is used for close statements (graceful shutdown)
func (s *statementStorage) Close() error {
	for _, stmt := range s.statements {
		if err := stmt.Close(); err != nil {
			return errors.Wrap(err, "can't close statement")
		}
	}

	return nil
}
