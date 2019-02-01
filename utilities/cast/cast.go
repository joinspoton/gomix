package cast

// GetBoolOrNull - Type asserts and returns the value as a boolean
func GetBoolOrNull(value interface{}) bool {
	b, ok := value.(bool)
	if !ok {
		return false
	}
	return b
}

// GetStringOrNull - Type asserts and returns the value as a string
func GetStringOrNull(value interface{}) string {
	s, ok := value.(string)
	if !ok {
		return ""
	}
	return s
}

// GetFloat64OrNull - Type asserts and returns the value as a float64
func GetFloat64OrNull(value interface{}) float64 {

}
