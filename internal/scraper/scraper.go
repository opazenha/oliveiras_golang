package scraper

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"

	"github.com/zenha/oliveiras/internal/models"
)

// Service handles scraping operations
type Service struct {
	pythonPath string
	scriptPath string
}

// NewService creates a new scraper service
func NewService(pythonPath, scriptPath string) *Service {
	return &Service{
		pythonPath: pythonPath,
		scriptPath: scriptPath,
	}
}

// ScrapeListings scrapes both Airbnb and Booking.com listings
func (s *Service) ScrapeListings(startDate, endDate string) (*models.ListingAnalysis, *models.ListingAnalysis, error) {
	cmd := exec.Command(s.pythonPath, s.scriptPath, startDate, endDate)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("cmd.CombinedOutput() failed:", err)
		log.Println("Output:", string(output))
		return nil, nil, err
	}
	
	return s.parseScrapingResults(string(output))
}

// parseScrapingResults parses the Python script output
func (s *Service) parseScrapingResults(output string) (*models.ListingAnalysis, *models.ListingAnalysis, error) {
	finalOutputSplit := strings.Split(output, "Airbnb Listings Data:")
	bookingStart := strings.Index(finalOutputSplit[1], "Booking Listings Data:")
	airbnbData := strings.TrimSpace(finalOutputSplit[1][:bookingStart])
	bookingData := strings.TrimSpace(finalOutputSplit[1][bookingStart:])
	bookingData = strings.ReplaceAll(bookingData, "Booking Listings Data:", "")

	airbnbData = strings.ReplaceAll(airbnbData, "'", "\"")
	bookingData = strings.ReplaceAll(bookingData, "'", "\"")

	var airbnbAnalysis models.ListingAnalysis
	var bookingAnalysis models.ListingAnalysis

	if err := json.Unmarshal([]byte(airbnbData), &airbnbAnalysis); err != nil {
		log.Println("Airbnb json.Unmarshal() failed:", err)
		return nil, nil, err
	}

	if err := json.Unmarshal([]byte(bookingData), &bookingAnalysis); err != nil {
		log.Println("Booking json.Unmarshal() failed:", err)
		return nil, nil, err
	}

	return &airbnbAnalysis, &bookingAnalysis, nil
}
