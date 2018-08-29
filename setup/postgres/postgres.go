package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Connect - Create a new PostGreSQL connection
func Connect(host string, port string, db string, username string, password string) (*gorm.DB, error) {
	connection := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", host, port, db, username, password)
	client, err := gorm.Open("postgres", connection)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to postgres with: %+v", connection)
	}
	return client, nil
}
