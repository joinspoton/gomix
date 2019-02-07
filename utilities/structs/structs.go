package structs

import (
	"reflect"
)

// GetFieldNames - Get field names of a struct as a slice
func GetFieldNames(struc interface{}) []string {
	e := reflect.ValueOf(&struc).Elem()
	fields := []string{}
	for i := 0; i < e.NumField(); i++ {
		varName := e.Type().Field(i).Name
		fields = append(fields, varName)
		// varType := e.Type().Field(i).Type
		// varValue := e.Field(i).Interface()
		// fmt.Printf("%v %v %v\n", varName, varType, varValue)
	}
	return fields
}
