package postgres

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jinzhu/gorm"
)

// BatchInsert - Since it is not officially supported in the GORM API
func BatchInsert(db *gorm.DB, table string, data []map[string]interface{}) {
	batchSize := 500
	var batches [][]map[string]interface{}
	for batchSize < len(data) {
		data, batches = data[batchSize:], append(batches, data[0:batchSize:batchSize])
	}
	batches = append(batches, data)

	columns := make([]string, len(data))
	i := 0
	for k := range data[0] {
		columns[i] = k
		i++
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
			for i := 0; i < len(batches[i]); i++ {
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
