package models

// AirbnbData represents an Airbnb listing with its metadata
type AirbnbData struct {
	ID struct {
		Oid string `json:"$oid"`
	} `json:"_id"`
	Timestamp  string  `json:"timestamp"`
	URL        string  `json:"url"`
	StartDate  string  `json:"start_date"`
	EndDate    string  `json:"end_date"`
	Listing    Listing `json:"listing"`
	InsertedAt string  `json:"inserted_at"`
}

// Listing represents the core listing data
type Listing struct {
	Name             string  `json:"name"`
	Price            float64 `json:"price"`
	Rating           float64 `json:"rating"`
	BedConfiguration string  `json:"bed_configuration"`
}

// BookingData represents a Booking.com listing
type BookingData struct {
	ID struct {
		Oid string `json:"$oid"`
	} `json:"_id"`
	Timestamp        string  `json:"timestamp"`
	URL              string  `json:"url"`
	StartDate        string  `json:"start_date"`
	EndDate          string  `json:"end_date"`
	Name             string  `json:"name"`
	Price            float32 `json:"price"`
	Rating           string  `json:"rating"`
	BedConfiguration string  `json:"bed_configuration"`
	InsertedAt       string  `json:"inserted_at"`
}

// ListingAnalysis represents analyzed data for listings
type ListingAnalysis struct {
	AveragePrice  float64 `json:"average_price"`
	HighestPrice  float64 `json:"highest_price"`
	LowestPrice   float64 `json:"lowest_price"`
	TotalListings int     `json:"total_listings"`
}
