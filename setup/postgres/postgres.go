package postgres

import (
	"fmt"
	"strings"
	"sync"

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

// BatchInsert - Since it is not officially supported in the GORM API
func BatchInsert(db *gorm.DB, table string, data []map[string]interface{}) {
	if len(data) <= 0 {
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

			db.Exec(query, values...)
		}(i)
	}

	wg.Wait()
}
