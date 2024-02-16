package db

import (
	"database/sql"
)

type iDatabase interface {
}

type database struct {
	db *sql.DB
}

var db iDatabase

func DB() iDatabase { return db }

func init() {

}
