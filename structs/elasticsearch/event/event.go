package event

import "time"

// Event - Structure and standardize raw events from different vendors
type Event struct {
	Vendor         string
	Name           string
	Description    string
	StartTimestamp time.Time
	EndTimestamp   time.Time
	Price          float64
	Currency       string
	Longitude      float64
	Latitude       float64
	City           string
	Country        string
}
