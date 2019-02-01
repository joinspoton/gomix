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

// BatchStructsInsert - Insert an array of structs into a table
func BatchStructsInsert(db *gorm.DB, table string, objArr []interface{}) {
	// If there is no data, nothing to do.
	if len(objArr) == 0 {
		return
	}

	batchSize := 500
	var batches [][]interface{}
	for batchSize < len(objArr) {
		objArr, batches = objArr[batchSize:], append(batches, objArr[0:batchSize:batchSize])
	}
	batches = append(batches, objArr)

	for _, batch := range batches {
		mainObj := batch[0]
		mainScope := db.NewScope(mainObj)
		mainFields := mainScope.Fields()
		quoted := make([]string, 0, len(mainFields))
		for i := range mainFields {
			// If primary key has blank value (0 for int, "" for string, nil for interface ...), skip it.
			// If field is ignore field, skip it.
			if (mainFields[i].IsPrimaryKey && mainFields[i].IsBlank) || (mainFields[i].IsIgnored) {
				continue
			}
			quoted = append(quoted, mainScope.Quote(mainFields[i].DBName))
		}

		placeholdersArr := make([]string, 0, len(batch))

		for _, obj := range batch {
			scope := db.NewScope(obj)
			fields := scope.Fields()
			placeholders := make([]string, 0, len(fields))
			for i := range fields {
				if (fields[i].IsPrimaryKey && fields[i].IsBlank) || (fields[i].IsIgnored) {
					continue
				}
				placeholders = append(placeholders, mainScope.AddToVars(fields[i].Field.Interface()))
			}
			placeholdersStr := "(" + strings.Join(placeholders, ", ") + ")"
			placeholdersArr = append(placeholdersArr, placeholdersStr)
			// add real variables for the replacement of placeholders' '?' letter later.
			mainScope.SQLVars = append(mainScope.SQLVars, scope.SQLVars...)
		}

		mainScope.Raw(fmt.Sprintf("INSERT INTO %s(%s) VALUES %s",
			table,
			strings.Join(quoted, ", "),
			strings.Join(placeholdersArr, ", "),
		))

		if _, err := mainScope.SQLDB().Exec(mainScope.SQL, mainScope.SQLVars...); err != nil {
			fmt.Printf("%+v\n", mainScope.SQLVars)
			panic(err)
		}
	}
}
