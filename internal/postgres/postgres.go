package postgres

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

const (
	UniqueViolation     = "23505"
	ForeignKeyViolation = "23503"
	CheckViolation      = "23514"
)

func PGInit(host string, port int, user string, password string, dbname string) (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to db")
	}
	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "can't ping db")
	}
	return db, nil
}
