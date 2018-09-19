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
	batches := len(data)/batchSize + 1

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
	wg.Add(batches)

	for i := 0; i < batches; i++ {
		go func(i int) {
			defer wg.Done()

			var inserts []string
			for i := 0; i < len(data)/len(columns); i++ {
				inserts = append(inserts, row)
			}

			stop := i + 500*len(headers)
			subset := data[i:stop]

			if stop > len(data) {
				subset = data[i:len(data)]
			}
			query := cmd + strings.Join(inserts, ",")
			db.Exec(query, data...)
		}(i)
	}

	wg.Wait()
}
