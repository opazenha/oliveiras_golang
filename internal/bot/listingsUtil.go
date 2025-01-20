package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/zenha/oliveiras/internal/models"
)

// formatAnalysisResponse formats the analysis results into a readable message
func formatAnalysisResponse(airbnb, booking *models.ListingAnalysis) string {
	return fmt.Sprintf("Airbnb Listings Data:\nAverage Price: %.2f\nHighest Price: %.2f\nLowest Price: %.2f\nTotal Listings: %d\n\nBooking Listings Data:\nAverage Price: %.2f\nHighest Price: %.2f\nLowest Price: %.2f\nTotal Listings: %d",
		airbnb.AveragePrice, airbnb.HighestPrice, airbnb.LowestPrice, airbnb.TotalListings,
		booking.AveragePrice, booking.HighestPrice, booking.LowestPrice, booking.TotalListings)
}

func checkAirbnbDataUpToDate(airbnbListings []models.AirbnbData) (bool, error) {
	for _, listing := range airbnbListings {
		if !strings.Contains(listing.InsertedAt, "Z") {
			listing.InsertedAt += "Z"
		}
		insertedAt, err := time.Parse(time.RFC3339, listing.InsertedAt)
		if err != nil {
			log.Println("Error parsing Airbnb inserted_at:", err)
			return false, err
		}

		if insertedAt.After(time.Now().AddDate(0, 0, -7)) {
			return true, nil
		}
	}
	return false, nil
}

func checkBookingDataUpToDate(bookingListings []models.BookingData) (bool, error) {
	for _, listing := range bookingListings {
		if !strings.Contains(listing.InsertedAt, "Z") {
			listing.InsertedAt += "Z"
		}
		insertedAt, err := time.Parse(time.RFC3339, listing.InsertedAt)
		if err != nil {
			log.Println("Error parsing Booking inserted_at:", err)
			return false, err
		}

		if insertedAt.After(time.Now().AddDate(0, 0, -7)) {
			return true, nil
		}
	}
	return false, nil
}

func filterBookingDataUpToDate(bookingListings []models.BookingData) ([]models.BookingData, error) {
	upToDateList := []models.BookingData{}
	for _, listing := range bookingListings {
		if !strings.Contains(listing.InsertedAt, "Z") {
			listing.InsertedAt += "Z"
		}
		insertedAt, err := time.Parse(time.RFC3339, listing.InsertedAt)
		if err != nil {
			log.Println("Error parsing Booking inserted_at:", err)
			return []models.BookingData{}, err
		}

		if insertedAt.After(time.Now().AddDate(0, 0, -7)) {
			upToDateList = append(upToDateList, bookingListings...)
		}
	}
	return upToDateList, nil
}

func filterAirbnbDataUpToDate(airbnbListings []models.AirbnbData) ([]models.AirbnbData, error) {
	upToDateList := []models.AirbnbData{}
	for _, listing := range airbnbListings {
		if !strings.Contains(listing.InsertedAt, "Z") {
			listing.InsertedAt += "Z"
		}
		insertedAt, err := time.Parse(time.RFC3339, listing.InsertedAt)
		if err != nil {
			log.Println("Error parsing Airbnb inserted_at:", err)
			return []models.AirbnbData{}, err
		}

		if insertedAt.After(time.Now().AddDate(0, 0, -7)) {
			upToDateList = append(upToDateList, airbnbListings...)
		}
	}
	return upToDateList, nil
}

func separateAirbnbByDate(listings []models.AirbnbData) (map[string][]models.AirbnbData, error) {
	result := make(map[string][]models.AirbnbData)
	for _, listing := range listings {
		result[listing.StartDate] = append(result[listing.StartDate], listing)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no listings found")
	}
	return result, nil
}

func separateBookingByDate(listings []models.BookingData) (map[string][]models.BookingData, error) {
	result := make(map[string][]models.BookingData)
	for _, listing := range listings {
		result[listing.StartDate] = append(result[listing.StartDate], listing)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no listings found")
	}
	return result, nil
}

func getAirbnbOutOfDateList(lists map[string][]models.AirbnbData) string {
	outdatedDates := []string{}
	for date, listings := range lists {
		isUpToDate, err := checkAirbnbDataUpToDate(listings)
		if err != nil {
			log.Println("Error checking Airbnb data:", err)
			return ""
		}
		if !isUpToDate {
			outdatedDates = append(outdatedDates, date)
		}
	}
	if len(outdatedDates) > 0 {
		return strings.Join(outdatedDates, ",")
	}
	return ""
}

func getBookingOutOfDateList(lists map[string][]models.BookingData) string {
	outdatedDates := []string{}
	for date, listings := range lists {
		isUpToDate, err := checkBookingDataUpToDate(listings)
		if err != nil {
			log.Println("Error checking Booking data:", err)
			return ""
		}
		if !isUpToDate {
			outdatedDates = append(outdatedDates, date)
		}
	}
	if len(outdatedDates) > 0 {
		return strings.Join(outdatedDates, ",")
	}
	return ""
}
