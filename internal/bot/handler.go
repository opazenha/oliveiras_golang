package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/zenha/oliveiras/internal/database"
	"github.com/zenha/oliveiras/internal/gemini"
	"github.com/zenha/oliveiras/internal/scraper"
	"github.com/zenha/oliveiras/internal/telegram"
	"github.com/zenha/oliveiras/pkg/config"
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
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize MongoDB client
	mongoClient, err := database.NewClient(cfg.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongoClient.Disconnect()

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

	case "/getprices":
		if len(parts) != 3 {
			return h.telegramClient.SendMessage(chatID, "Usage: /getprices start_date end_date")
		}

		startDate := parts[1]
		endDate := parts[2]

		airbnbListings, err := mongoClient.GetAirbnbByDate(startDate, endDate)
		if err != nil {
			return h.telegramClient.SendMessage(chatID, "Error: "+err.Error())
		}
		bookingListings, err := mongoClient.GetBookingByDate(startDate, endDate)
		if err != nil {
			return h.telegramClient.SendMessage(chatID, "Error: "+err.Error())
		}

		airbnbDateList, err := separateAirbnbByDate(airbnbListings)
		if err != nil {
			return h.telegramClient.SendMessage(chatID, "Error: "+err.Error())
		}
		bookingDateList, err := separateBookingByDate(bookingListings)
		if err != nil {
			return h.telegramClient.SendMessage(chatID, "Error: "+err.Error())
		}

		airbnbOutOfDateList := getAirbnbOutOfDateList(airbnbDateList)
		bookingOutOfDateList := getBookingOutOfDateList(bookingDateList)
		if airbnbOutOfDateList != "" || bookingOutOfDateList != "" {
			return h.telegramClient.SendMessage(chatID, "Data is not up to date. Please run /scrape command. Airbnbs: "+airbnbOutOfDateList+". Bookings: "+bookingOutOfDateList+".")
		}

		geminiClient, err := gemini.NewClient(cfg.GeminiKey)
		if err != nil {
			log.Fatal("Failed to create Gemini client:", err)
		}
		bookingPrices, err := gemini.GenerateContent(geminiClient, gemini.PrepareBookingPrompt(bookingListings))
		if err != nil {
			return h.telegramClient.SendMessage(chatID, "Error: "+err.Error())
		}
		airbnbPrices, err := gemini.GenerateContent(geminiClient, gemini.PrepareAirbnbPrompt(airbnbListings))
		if err != nil {
			return h.telegramClient.SendMessage(chatID, "Error: "+err.Error())
		}

		telegramMessage := fmt.Sprintf("Booking Prices:\n%v\n\nAirbnb Prices:\n%v", bookingPrices, airbnbPrices)
		return h.telegramClient.SendMessage(chatID, telegramMessage)

	default:
		return h.telegramClient.SendMessage(chatID, "Unknown command: "+parts[0]+".\nUse /scrape command to scrape and analyze listings.\nUse /getprices command to get the prices suggestions.")
	}
}
