package postgres

import (
	"github.com/go-pg/pg"
)

// Connect - Create a new PostGreSQL connection
func Connect(address string, db string, username string, password string) *pg.DB {
	return pg.Connect(&pg.Options{
		Addr:     address,
		Database: db,
		User:     username,
		Password: password,
	})
}
