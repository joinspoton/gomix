package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Connect - Create a new PostGreSQL connection
func Connect(host string, port int, db string, username string, password string) (*gorm.DB, error) {
	connection := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s", host, port, db, username, password)
	client, err := gorm.Open("postgres", connection)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to postgres with: %+v", connection)
	}
	return client, nil
}
