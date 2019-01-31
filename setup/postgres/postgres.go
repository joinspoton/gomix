package postgres

import (
	"fmt"
	"gomix/utilities/paramstore"
	"gomix/utilities/system"
	"strings"
	"sync"

	"github.com/jinzhu/gorm"
)

// Connect - Create a new PostGreSQL connection without specifying credentials
func Connect(db string) (*gorm.DB, error) {
	path := fmt.Sprintf("/%s/postgres/", system.GetEnv("stage", "staging"))

	host, _ := paramstore.GetConfig(path + "host")
	port, _ := paramstore.GetConfig(path + "port")
	username, _ := paramstore.GetConfig(path + "username")
	password, _ := paramstore.GetConfig(path + "password")

	return ManuallyConnect(
		host,
		port,
		db,
		username,
		password,
	)
}

// ConnectToProd - Create a new PostGreSQL connection to Production Database
func ConnectToProd(db string) (*gorm.DB, error) {
	path := "/production/postgres/"

	host, _ := paramstore.GetConfig(path + "host")
	port, _ := paramstore.GetConfig(path + "port")
	username, _ := paramstore.GetConfig(path + "username")
	password, _ := paramstore.GetConfig(path + "password")

	return ManuallyConnect(
		host,
		port,
		db,
		username,
		password,
	)
}

// ConnectToStaging - Create a new PostGreSQL connection to Staging Database
func ConnectToStaging(db string) (*gorm.DB, error) {
	path := "/staging/postgres/"

	host, _ := paramstore.GetConfig(path + "host")
	port, _ := paramstore.GetConfig(path + "port")
	username, _ := paramstore.GetConfig(path + "username")
	password, _ := paramstore.GetConfig(path + "password")

	return ManuallyConnect(
		host,
		port,
		db,
		username,
		password,
	)
}

// ManuallyConnect - Create a new PostGreSQL connection with specified credentials
func ManuallyConnect(host string, port string, db string, username string, password string) (*gorm.DB, error) {
	connection := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", host, port, db, username, password)
	client, err := gorm.Open("postgres", connection)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to postgres with: %+v", connection)
	}
	return client, nil
}

// BatchInsert - Since it is not officially supported in the GORM API
func BatchInsert(db *gorm.DB, table string, data []map[string]interface{}, onConflict string) {
	if len(data) < 1 {
		return
	}

	batchSize := 500
	var batches [][]map[string]interface{}
	for batchSize < len(data) {
		data, batches = data[batchSize:], append(batches, data[0:batchSize:batchSize])
	}
	batches = append(batches, data)

	var columns []string
	for k := range data[0] {
		columns = append(columns, k)
	}

	var placeholders []string
	for i := 0; i < len(columns); i++ {
		placeholders = append(placeholders, "?")
	}
	row := fmt.Sprintf("(%s)", strings.Join(placeholders, ","))

	cmd := fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES ",
		table,
		strings.Join(columns, ","),
	)

	var wg sync.WaitGroup
	wg.Add(len(batches))

	for i := 0; i < len(batches); i++ {
		go func(i int) {
			defer wg.Done()

			var inserts []string
			for j := 0; j < len(batches[i]); j++ {
				inserts = append(inserts, row)
			}

			query := cmd + strings.Join(inserts, ",")

			var values []interface{}
			for _, entry := range batches[i] {
				for _, column := range columns {
					values = append(values, entry[column])
				}
			}

			if onConflict != "" {
				query += "\n" + onConflict
			}

			db.Exec(query, values...)
		}(i)
	}

	wg.Wait()
}
