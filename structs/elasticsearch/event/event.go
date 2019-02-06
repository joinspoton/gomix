package event

import "time"

// Event - Structure and standardize raw events from different vendors
type Event struct {
	Vendor      string `gorm:"unique_index:all_columns"`
	Name        string `gorm:"unique_index:all_columns"`
	Description string `gorm:"unique_index:all_columns"`
	StartTime   time.Time `gorm:"unique_index:all_columns"`
	EndTime     time.Time `gorm:"unique_index:all_columns"`
	PriceMin    float64 `gorm:"unique_index:all_columns"`
	PriceMax    float64 `gorm:"unique_index:all_columns"`
	Currency    string `gorm:"unique_index:all_columns"`
	Address     string `gorm:"unique_index:all_columns"`
	Longitude   float64 `gorm:"unique_index:all_columns"`
	Latitude    float64 `gorm:"unique_index:all_columns"`
	City        string `gorm:"unique_index:all_columns"`
	Country     string `gorm:"unique_index:all_columns"`
}
