package event

import "time"

// Event - Structure and standardize raw events from different vendors
type Event struct {
	Vendor      string
	Name        string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	PriceMin    float64
	PriceMax    float64
	Currency    string
	Address     string
	Longitude   float64
	Latitude    float64
	City        string
	Country     string
	Tags        string
	ImageURL    string
}
