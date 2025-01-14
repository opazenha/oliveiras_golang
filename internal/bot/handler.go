package bot

import (
	"fmt"
	"strings"

	"github.com/zenha/oliveiras/internal/models"
	"github.com/zenha/oliveiras/internal/scraper"
	"github.com/zenha/oliveiras/internal/telegram"
)

// Handler manages bot message handling
type Handler struct {
	telegramClient *telegram.Client
	scraperService *scraper.Service
}

// NewHandler creates a new bot handler
func NewHandler(telegramClient *telegram.Client, scraperService *scraper.Service) *Handler {
	return &Handler{
		telegramClient: telegramClient,
		scraperService: scraperService,
	}
}

// HandleMessage processes incoming bot messages
func (h *Handler) HandleMessage(chatID int, message string) error {
	parts := strings.Split(message, " ")
	if len(parts) == 0 {
		return nil
	}

	switch parts[0] {
	case "/scrape":
		if len(parts) != 3 {
			return h.telegramClient.SendMessage(chatID, "Usage: /scrape start_date end_date")
		}

		startDate := parts[1]
		endDate := parts[2]

		airbnbAnalysis, bookingAnalysis, err := h.scraperService.ScrapeListings(startDate, endDate)
		if err != nil {
			return h.telegramClient.SendMessage(chatID, "Error: "+err.Error())
		}

		response := formatAnalysisResponse(airbnbAnalysis, bookingAnalysis)
		return h.telegramClient.SendMessage(chatID, response)

	default:
		return h.telegramClient.SendMessage(chatID, "Unknown command: "+parts[0])
	}
}

// formatAnalysisResponse formats the analysis results into a readable message
func formatAnalysisResponse(airbnb, booking *models.ListingAnalysis) string {
	return fmt.Sprintf("Airbnb Listings Data:\nAverage Price: %.2f\nHighest Price: %.2f\nLowest Price: %.2f\nTotal Listings: %d\n\nBooking Listings Data:\nAverage Price: %.2f\nHighest Price: %.2f\nLowest Price: %.2f\nTotal Listings: %d",
		airbnb.AveragePrice, airbnb.HighestPrice, airbnb.LowestPrice, airbnb.TotalListings,
		booking.AveragePrice, booking.HighestPrice, booking.LowestPrice, booking.TotalListings)
}
