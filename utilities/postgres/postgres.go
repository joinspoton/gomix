package postgres

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

// BatchInsert - Since it is not officially supported in the GORM API
func BatchInsert(db *gorm.DB, table string, columns []string, data []interface{}) error {
	cmd := fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES ",
		table,
		strings.Join(columns, ","),
	)

	var placeholders []string
	for i := 0; i < len(columns); i++ {
		placeholders = append(placeholders, "?")
	}
	row := fmt.Sprintf("(%s)", strings.Join(placeholders, ","))

	var inserts []string
	for i := 0; i < len(data)/len(columns); i++ {
		inserts = append(inserts, row)
	}

	cmd += strings.Join(inserts, ",")

	if err := db.Exec(cmd, data...).Error; err != nil {
		return err
	}
	return nil
}
