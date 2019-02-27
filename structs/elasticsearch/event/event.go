package event

import (
	"fmt"
	"time"
)

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

// AddDefaultValues - Give the event default values
func AddDefaultValues(event *Event) {
	event.Vendor = ""
	event.Name = ""
	event.Description = ""
	event.Currency = "USD"
	event.StartTime = time.Time{}
	event.EndTime = time.Time{}
	event.PriceMin = float64(0)
	event.PriceMax = float64(0)
	event.City = ""
	event.Country = ""
	event.Longitude = float64(0)
	event.Latitude = float64(0)
	event.ImageURL = ""
	event.Tags = ""
	fmt.Println(event)
}
