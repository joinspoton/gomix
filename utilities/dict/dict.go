package dict

import "strings"

// LowercaseKeys - Make the field keys lowercase
func LowercaseKeys(data map[string]interface{}) map[string]interface{} {
	converted := make(map[string]interface{}, len(data))
	for k, v := range data {
		converted[strings.ToLower(k)] = v
	}
	return converted
}
