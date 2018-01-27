package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DB interface {
	user
	deposit
}

func GetDB(kind, path string) (DB, error) {
	database, err := sql.Open(kind, path)
	if err != nil {
		log.Printf("error opening database connection\n%v", err)
		return nil, err
	}
	return &sqlDb{database}, nil
}

type NotFound struct {
	resource string
}

func (err *NotFound) Error() string {
	return fmt.Sprintf("%v not found", err.resource)
}

func IsNotFound(err error) bool {
	if _, ok := err.(*NotFound); ok {
		return true
	}
	return false
}

type BadQuery struct {
	reason string
}

func (err *BadQuery) Error() string {
	return err.reason
}

func IsBadQuery(err error) bool {
	if _, ok := err.(*BadQuery); ok {
		return true
	}
	return false
}

type sqlDb struct {
	db *sql.DB
}
